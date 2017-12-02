package main
import (
      "os"
      "time"
      "strings"
      "strconv"
      "net/http"
      "database/sql"; _ "github.com/lib/pq";
      "gocronlib"
)


// Global const and vars
const version string     = "2.0.0"
const socket string      = ":8080"
const errorResp string   = "Internal Server Error\n"
const contentType string = "plain/text"

var verbose bool  = false       // Flag enabling / disabling verbosity
var args []string = os.Args     // Command line arguments


func main() {
      // Parse arguments
      handleArgs(args)

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
      var cronJob gocronlib.Cron
      var method string = ""

      switch req.Method {
      case "GET":
            method = "GET"
            cronJob.Cronname = req.URL.Query().Get("cronname")
            cronJob.Account = req.URL.Query().Get("account")
            cronJob.Email = req.URL.Query().Get("email")
            cronJob.Frequency = req.URL.Query().Get("frequency")
            cronJob.Lastruntime = strconv.Itoa(currentTime)
            cronJob.Ipaddress = socket[0]
      default:
            // Log an error and do not respond
            gocronlib.CronLog("Incoming request from " + cronJob.Ipaddress +
                   " is not a GET or POST.", verbose)
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
            gocronlib.CronLog(method + " from " + cronJob.Ipaddress + " not valid. Dropping.", verbose)
      }
}


// Return a 201 Created
func returnCreated(w http.ResponseWriter) {
      w.Header().Set("Content-Type", contentType)
      w.WriteHeader(http.StatusCreated)
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


func updateDatabase(c gocronlib.Cron) bool {
      // Build the database query
      var query string
      query = "INSERT INTO gocron " +
                   "(cronname, account, email, ipaddress, frequency, lastruntime, alerted) " +
              "VALUES ('" +
                   c.Cronname + "','" +
                   c.Account + "','" +
                   c.Email + "','" +
                   c.Ipaddress + "','" +
                   c.Frequency + "','" +
                   c.Lastruntime + "','" +
                   "false" + "') " +
              "ON CONFLICT (cronname, account) DO UPDATE " +
                   "SET email = " + "'" + c.Email + "'," +
                   "ipaddress = " + "'" + c.Ipaddress + "'," +
                   "frequency = " + "'" + c.Frequency + "'," +
                   "lastruntime = " + "'" + c.Lastruntime + "'" +
              ";"

      db, err := sql.Open("postgres", gocronlib.DatabaseString(gocronlib.GetConfig(verbose)))
      if err != nil {
            gocronlib.CheckError(err, verbose)
            return false
      }
      defer db.Close()

      _, err = db.Exec(query)
      if err != nil {
            gocronlib.CheckError(err, verbose)
            return false

      } else {
            gocronlib.CronLog("Heartbeat from " + c.Cronname + ": " + c.Account + " \n", verbose)
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
                  gocronlib.CronLog("gocron started with --verbose.", verbose)
                  return

            } else {
                  return
            }
      }
}


// Function validates SQL variables
func validateArgs(c gocronlib.Cron) bool {
      // Flag determines the return value
      var valid bool = false

      // Perform validation of parameters
      if checkLength(c) == true {
            valid = true
      }

      // Log result if verbose is enabled
      if verbose == true {
            if valid == true {
                  gocronlib.CronLog("Parameters from " + c.Ipaddress + " passed validation", verbose)
                  return true

            } else {
                  gocronlib.CronLog("Parameters from " + c.Ipaddress + " failed validation!", verbose)
                  return false
            }
      }

      // Return true or false
      return valid
}


// Validate that parameters are present
func checkLength(c gocronlib.Cron) bool {
      if len(c.Account) == 0 {
            return false

      } else if len(c.Cronname) == 0 {
            return false

      } else if len(c.Email) == 0 {
            return false

      } else if len(c.Frequency) == 0 {
            return false

      } else if len(c.Ipaddress) == 0 {
            return false

      } else if len(c.Lastruntime) == 0 {
            return false

      } else {
            return true
      }
}
