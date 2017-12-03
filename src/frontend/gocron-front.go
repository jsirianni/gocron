package main
import (
      "os"
      "time"
      "strings"
      "strconv"
      "net/http"
      "gocronlib"
)


const version string     = "2.0.2"
const libVersion string  = gocronlib.Version

const socket string      = ":8080"
const errorResp string   = "Internal Server Error\n"
const contentType string = "plain/text"

var verbose bool  = false       // Flag enabling / disabling verbosity
var args []string = os.Args     // Command line arguments


func main() {
      // Parse arguments
      if len(os.Args) > 1 {
            // Return the current version
            if strings.Contains(args[1], "--version") {
                  println("gocron-front version: " + version)
                  println("gocronlib version: " + libVersion)
                  os.Exit(0)
            }
            // When enabled, all logging will also print to screen
            if strings.Contains(args[1], "--verbose") {
                  verbose = true
                  gocronlib.CronLog("gocron started with --verbose.", verbose)
            }
      }

      // Start the web server on port 8080
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(socket, nil)
}


// Validate the request and then pass to updateDatabase()
func cronStatus(resp http.ResponseWriter, req *http.Request) {
      var currentTime int = int(time.Now().Unix())
      var socket = strings.Split(req.RemoteAddr, ":")
      var c gocronlib.Cron
      var method string = ""

      switch req.Method {
      case "GET":
            method = "GET"
            c.Cronname = req.URL.Query().Get("cronname")
            c.Account = req.URL.Query().Get("account")
            c.Email = req.URL.Query().Get("email")
            c.Frequency = req.URL.Query().Get("frequency")
            c.Lastruntime = strconv.Itoa(currentTime)
            c.Ipaddress = socket[0]

      case "POST":
            gocronlib.CronLog("POST not yet supported: " + c.Ipaddress, verbose)
            return

      default:
            // Log an error and do not respond
            gocronlib.CronLog("Incoming request from " + c.Ipaddress + " is not a GET or POST.", verbose)
            return
      }

      if validateParams(c) == true {
            if updateDatabase(c) == true {
                  returnCreated(resp)

            } else {
                  returnServerError(resp)
            }

      } else {
            returnNotFound(resp)
            gocronlib.CronLog(method + " from " + c.Ipaddress + " not valid. Dropping.", verbose)
      }
}


// Return a 201 Created
func returnCreated(resp http.ResponseWriter) {
      resp.Header().Set("Content-Type", contentType)
      resp.WriteHeader(http.StatusCreated)
}


// Return a 500 Server Error
func returnServerError(resp http.ResponseWriter) {
      resp.Header().Set("Content-Type", contentType)
      resp.WriteHeader(http.StatusInternalServerError)
      resp.Write([]byte(errorResp))
}


// Return 404 Not Found
func returnNotFound(resp http.ResponseWriter) {
      resp.WriteHeader(http.StatusNotFound)
}


func updateDatabase(c gocronlib.Cron) bool {
      // Build the database query
      var query string
      var result bool

      // Insert and update if already exist
      query = "INSERT INTO gocron " +
              "(cronname, account, email, ipaddress, frequency, lastruntime, alerted) " +
              "VALUES ('" +
              c.Cronname + "','" + c.Account + "','" + c.Email + "','" +
              c.Ipaddress + "','" + c.Frequency + "','" + c.Lastruntime + "','" + "false" + "') " +
              "ON CONFLICT (cronname, account) DO UPDATE " +
              "SET email = " + "'" + c.Email + "'," + "ipaddress = " + "'" + c.Ipaddress + "'," +
              "frequency = " + "'" + c.Frequency + "'," + "lastruntime = " + "'" + c.Lastruntime + "';"

      // Execute query
      _, result = gocronlib.QueryDatabase(query, verbose)
      if result == true {
            gocronlib.CronLog("Heartbeat from " + c.Cronname + ": " + c.Account + " \n", verbose)
            return true

      } else {
            return false
      }
}


// Function validates SQL variables
func validateParams(c gocronlib.Cron) bool {
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
