package main
import (
      "fmt"
      "time"
      "flag"
      "strings"
      "strconv"
      "net/http"
      "github.com/jsirianni/gocronlib"
)


const (
      version string     = "2.0.10"
      libVersion string  = gocronlib.Version
      errorResp string   = "Internal Server Error"
      contentType string = "plain/text"
)

var ( // Flags set in main()
      port string
      verbose bool
      getVersion bool
)


func main() {
      flag.BoolVar(&getVersion, "version", false, "Get the version and then exit")
      flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
      flag.StringVar(&port, "p", "8080", "Listening port for the web server")
      flag.Parse()

      if getVersion == true {
            fmt.Println("gocron-front version: " + version)
            fmt.Println("gocronlib version: " + libVersion)
            return
      }

      if verbose == true {
            fmt.Println("Verbose mode enabled")
            fmt.Println("gocron-front version: " + version)
            fmt.Println("gocronlib version: " + libVersion)
            fmt.Println("Starting web server on port: " + port)
      }

      // Start the web server
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":" + port, nil)
}


// Validate the request and then pass to updateDatabase()
func cronStatus(resp http.ResponseWriter, req *http.Request) {
      var (
            currentTime int = int(time.Now().Unix())
            socket = strings.Split(req.RemoteAddr, ":")
            c gocronlib.Cron
            method string = ""
      )

      switch req.Method {
      case "GET":
            method = "GET"
            c.Cronname    = req.URL.Query().Get("cronname")
            c.Account     = req.URL.Query().Get("account")
            c.Email       = req.URL.Query().Get("email")
            c.Frequency   = req.URL.Query().Get("frequency")
            c.Lastruntime = strconv.Itoa(currentTime)
            c.Ipaddress   = socket[0]

            // If x = 1, set c.Site to true
            x, err  := strconv.Atoi(req.URL.Query().Get("site"))
            if err == nil && x == 1 {
                  c.Site = true
            } else {
                  c.Site = false
            }

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
      var (
            query string
            result bool
      )

      // Insert and update if already exist
      query = "INSERT INTO gocron " +
              "(cronname, account, email, ipaddress, frequency, lastruntime, alerted, site) " +
              "VALUES ('" +
              c.Cronname + "','" + c.Account + "','" + c.Email + "','" + c.Ipaddress + "','" +
              c.Frequency + "','" + c.Lastruntime + "','" + "false" + "','" + strconv.FormatBool(c.Site) + "') " +
              "ON CONFLICT (cronname, account) DO UPDATE " +
              "SET email = " + "'" + c.Email + "'," + "ipaddress = " + "'" + c.Ipaddress + "'," +
              "frequency = " + "'" + c.Frequency + "'," + "lastruntime = " + "'" + c.Lastruntime + "', " +
              "site = " + "'" + strconv.FormatBool(c.Site) + "';"

      // Execute query
      rows, result := gocronlib.QueryDatabase(query, verbose)
      defer rows.Close()
      if result == true {
            gocronlib.CronLog("Heartbeat from " + c.Cronname + ": " + c.Account + " \n", verbose)
            return true

      } else {
            return false
      }
}


// Function validates SQL variables
func validateParams(c gocronlib.Cron) bool {

      var valid bool = false  // Flag determines the return value

      if checkLength(c) == true {
            valid = true
      }

      if verbose == true {
            if valid == true {
                  gocronlib.CronLog("Parameters from " + c.Ipaddress + " passed validation", verbose)
                  return true

            } else {
                  gocronlib.CronLog("Parameters from " + c.Ipaddress + " failed validation!", verbose)
                  return false
            }
      }

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
