// Version 1.0.2
// Debian 9 Officially supported

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
      cronname    string  // Name of the cronjob
      account     string  // Account the job belongs to
      email       string  // Address to send alerts to
      ipaddress   string  // Source IP address
      frequency   string  // How often a job should check in
      tolerance   string  // Additional time before an alert is thrown TODO Depricate
      lastruntime string  // Unix timestamp
      alerted     bool    // set to true if an alert has already been thrown
}

var config Config
var version string = "1.0.2"
var verbose bool = false



func main() {
      // Handle optional command line args
      if len(os.Args) > 1 {
            var args []string = os.Args

            // Return the current version --version
            if strings.Contains(args[1], "--version") {
                  println(version)
                  return

            // Enable verbose logging. All syslog will be printed standard out
            } else if strings.Contains(args[1], "--verbose") {
                  verbose = true
            }
      }


      // Read in the config file
      yamlFile, err := ioutil.ReadFile("/etc/gocron/.config.yml")
      if err != nil {
            checkError(err)
            return
      }


      // Set the global config var
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            checkError(err)
            return
      }


      // Start the status checking timer on a new thread
      go timer()


      // Start the web server on port 8080
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}



// Called from main. Parses a GET / POST into a cron sctruct
// and then passes them to the updateDatabase() function
func cronStatus(w http.ResponseWriter, req *http.Request) {
      // Method agnostic vars
      var currentTime int = int(time.Now().Unix())
      var socket = strings.Split(req.RemoteAddr, ":")
      var cronJob Cron
      var method string = ""


      // Handle GET / POST methods
      switch req.Method {
      case "GET":
            method = "GET"
            cronJob.cronname = req.URL.Query().Get("cronname")
            cronJob.account = req.URL.Query().Get("account")
            cronJob.email = req.URL.Query().Get("email")
            cronJob.frequency = req.URL.Query().Get("frequency")
            cronJob.tolerance = req.URL.Query().Get("tolerance")
            cronJob.lastruntime = strconv.Itoa(currentTime)
            cronJob.ipaddress = socket[0]

      case "POST":
            // POST not supported currently
            method = "POST"
            cronJob.cronname = ""
            cronJob.account = ""
            cronJob.email = ""
            cronJob.frequency = ""
            cronJob.tolerance = ""
            cronJob.lastruntime = strconv.Itoa(currentTime)
            cronJob.ipaddress = socket[0]

            return  // TODO Remove the return once POST is handled

      default:
            // Log an error
            cronLog("Incoming request from " + cronJob.ipaddress + " is not a GET or POST.")
            return
      }


      if checkLength(cronJob) == true {
            go updateDatabase(cronJob)

      } else {
            cronLog(method + " from " + cronJob.ipaddress + " not valid. Dropping.")
      }
}


func updateDatabase(c Cron) {
      // Build the database query
      var query string
      query = "INSERT INTO gocron (cronname, account, email, ipaddress, frequency, tolerance, lastruntime, alerted) VALUES ('" +
             c.cronname + "','" +
             c.account + "','" +
             c.email + "','" +
             c.ipaddress + "','" +
             c.frequency + "','" +
             c.tolerance + "','" +
             c.lastruntime + "','" +
             "false" + "') " +
             "ON CONFLICT (cronname, account) DO UPDATE " +
             "SET email = " + "'" + c.email + "'," +
             "ipaddress = " + "'" + c.ipaddress + "'," +
             "frequency = " + "'" + c.frequency + "'," +
             "tolerance = " + "'" + c.tolerance + "'," +
             "lastruntime = " + "'" + c.lastruntime + "'" +
             ";"


      db, err := sql.Open("postgres", databaseString())
      if err != nil {
            checkError(err)
      }
      defer db.Close()


      _, err = db.Exec(query)
      if err != nil {
            checkError(err)

      } else {
            cronLog("Heartbeat from " + c.cronname + ": " + c.account + " \n")
      }
}
