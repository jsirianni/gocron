package main
import (
      "os"
      "time"
      "strings"
      "strconv"
      "net/http"
      "io/ioutil"
      "gopkg.in/yaml.v2"
      "database/sql"; _ "github.com/lib/pq";
)


type Config struct {
      Dbfqdn       string
      Dbport       string
      Dbuser       string
      Dbpass       string
      Dbdatabase   string
      Smtpserver   string
      Smtpport     string
      Smtpaddress  string
      Smtppassword string
}

type Cron struct {
      cronname    string   // Name of the cronjob
      account     string   // Account the job belongs to
      email       string   // Address to send alerts to
      ipaddress   string   // Source IP address
      frequency   string   // How often a job should check in
      lastruntime string   // Unix timestamp
      alerted     bool     // set to true if an alert has already been thrown
}


// Global const and vars
const version string     = "1.0.7"
const confPath string    = "/etc/gocron/config.yml"
const socket string      = ":8080"
const okResp string      = "Update received\n"
const errorResp string   = "Internal Server Error\n"
const contentType string = "plain/text"

var config Config               // Stores configuration values in a Cron struct
var verbose bool  = false       // Flag enabling / disabling verbosity
var args []string = os.Args     // Command line arguments


func main() {
      // Parse arguments
      handleArgs(args)


      // Read in the config file
      yamlFile, err := ioutil.ReadFile(confPath)
      if err != nil {
            checkError(err)
            os.Exit(1)
      }


      // Set the global config var
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            checkError(err)
            os.Exit(1)
      }


      // Start the status checking timer on a new thread
      go timer()


      // Start the web server on port 8080
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(socket, nil)
}



// Validate the request and then pass to updateDatabase()
func cronStatus(w http.ResponseWriter, req *http.Request) {

      var currentTime int = int(time.Now().Unix())
      var socket = strings.Split(req.RemoteAddr, ":")
      var cronJob Cron
      var method string = ""


      switch req.Method {
      case "GET":
            method = "GET"
            cronJob.cronname = req.URL.Query().Get("cronname")
            cronJob.account = req.URL.Query().Get("account")
            cronJob.email = req.URL.Query().Get("email")
            cronJob.frequency = req.URL.Query().Get("frequency")
            cronJob.lastruntime = strconv.Itoa(currentTime)
            cronJob.ipaddress = socket[0]

      default:
            // Log an error and do not respond
            cronLog("Incoming request from " + cronJob.ipaddress +
                   " is not a GET or POST.")
            return
      }


      if validateArgs(cronJob) == true {
            if updateDatabase(cronJob) == true {
                  returnCreated(w)

            } else {
                  returnServerError(w)
            }

      } else {
            returnNotFound(w)
            cronLog(method + " from " + cronJob.ipaddress + " not valid. Dropping.")
      }
}



// Return a 201 Created
func returnCreated(w http.ResponseWriter) {
      w.Header().Set("Content-Type", contentType)
      w.WriteHeader(http.StatusCreated)
      w.Write([]byte(okResp))
}

// Return a 500 Server Error
func returnServerError(w http.ResponseWriter) {
      w.Header().Set("Content-Type", contentType)
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(errorResp))
}

// Return 404 Not Found
func returnNotFound(w http.ResponseWriter) {
      w.WriteHeader(http.StatusNotFound)
}



func updateDatabase(c Cron) bool {
      // Build the database query
      var query string
      query = "INSERT INTO gocron (cronname, account, email, ipaddress, frequency, lastruntime, alerted) VALUES ('" +
             c.cronname + "','" +
             c.account + "','" +
             c.email + "','" +
             c.ipaddress + "','" +
             c.frequency + "','" +
             c.lastruntime + "','" +
             "false" + "') " +
             "ON CONFLICT (cronname, account) DO UPDATE " +
             "SET email = " + "'" + c.email + "'," +
             "ipaddress = " + "'" + c.ipaddress + "'," +
             "frequency = " + "'" + c.frequency + "'," +
             "lastruntime = " + "'" + c.lastruntime + "'" +
             ";"


      db, err := sql.Open("postgres", databaseString())
      if err != nil {
            checkError(err)
            return false
      }
      defer db.Close()


      _, err = db.Exec(query)
      if err != nil {
            checkError(err)
            return false

      } else {
            cronLog("Heartbeat from " + c.cronname + ": " + c.account + " \n")
            return true
      }
}



func handleArgs(args []string) {
      if len(os.Args) > 1 {

            // Return the current version
            if strings.Contains(args[1], "--version") {
                  println(version)
                  os.Exit(0)

            // When enabled, all logging will also print to screen
            } else if strings.Contains(args[1], "--verbose") {
                  verbose = true
                  cronLog("gocron started with --verbose.")
                  return

            } else {
                  return
            }
      }
}
