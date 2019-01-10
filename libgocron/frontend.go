package libgocron
import (
	"net/http"
    "strconv"
	"strings"
	"time"
	"io/ioutil"
    "encoding/json"
)


// StartFrontend starts the gocron frontend server
func StartFrontend(c Config, frontendPort string, v bool) {

	// set the global struct "config"
	config = c

	verbose = v

	CronLog("verbose mode enabled", verbose)
	CronLog("gocron-front version: " + VERSION, verbose)
	CronLog("starting web server on port: " + frontendPort, verbose)

	http.HandleFunc("/", incomingCron)
	http.HandleFunc("/healthcheck", frontEndHealthCheck)
	http.ListenAndServe(":"+frontendPort, nil)
}


// return http status 200 if connection to database is healthy
func frontEndHealthCheck(resp http.ResponseWriter, req *http.Request) {
    remote_ip := strings.Split(req.RemoteAddr, ":")[0]
	CronLog("healthcheck from: " + remote_ip, false)
	returnOk(resp)
}


// Validate the request and then pass to updateDatabase()
func incomingCron(resp http.ResponseWriter, req *http.Request) {
	var (
		currentTime int = int(time.Now().Unix())
		socket          = strings.Split(req.RemoteAddr, ":")
		c           Cron
		method      string = ""
	)

	switch req.Method {
	case "GET":
		method = "GET"
		c.Cronname = req.URL.Query().Get("cronname")
		c.Account = req.URL.Query().Get("account")
		c.Email = req.URL.Query().Get("email")
		c.Frequency = stringToInt(req.URL.Query().Get("frequency"))
		c.Lastruntime = currentTime
		c.Ipaddress = socket[0]

		// If x = 1, set c.Site to true
		// NOTE: Depricating the site feature. It was never implemented and is causing
		// technical debt
		x, err := strconv.Atoi(req.URL.Query().Get("site"))
		if err == nil && x == 1 {
			c.Site = true
		} else {
			c.Site = false
		}

	case "POST":
		method = "POST"

		// read the request body into a byte array
		payload, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
		 	CronLog(err.Error(), verbose)
		}

		if err := json.Unmarshal(payload, &c); err != nil {
			CronLog(err.Error(), verbose)
		}
		c.Lastruntime = currentTime
		c.Ipaddress = socket[0]


	default:
		// Log an error and do not respond
		CronLog("Incoming request from "+c.Ipaddress+" is not a GET or POST.", verbose)
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
		CronLog(method+" from "+c.Ipaddress+" not valid. Dropping.", verbose)
	}
}


// Return 200 OK
func returnOk(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "plain/text")
	resp.WriteHeader(http.StatusOK)
}


// Return a 201 Created
func returnCreated(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "plain/text")
	resp.WriteHeader(http.StatusCreated)
}


// Return a 500 Server Error
func returnServerError(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "plain/text")
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte("Internal Server Error"))
}


// Return 404 Not Found
func returnNotFound(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
}


// Function validates SQL variables
func validateParams(c Cron) bool {

	var valid bool = false // Flag determines the return value

	if checkLength(c) == true {
		valid = true
	}

	if verbose == true {
		if valid == true {
			CronLog("Parameters from "+c.Ipaddress+" passed validation", verbose)
			return true

		} else {
			CronLog("Parameters from "+c.Ipaddress+" failed validation!", verbose)
			return false
		}
	}

	return valid
}


// Validate that parameters are present
// Validate that ints are not -1 (failed conversion in gocronlib StringToInt())
func checkLength(c Cron) bool {
	if len(c.Account) == 0 {
		return false

	} else if len(c.Cronname) == 0 {
		return false

	} else if len(c.Email) == 0 {
		return false

	} else if c.Frequency == -1 {
		return false

	} else if len(c.Ipaddress) == 0 {
		return false

	} else if c.Lastruntime == -1 {
		return false

	} else {
		return true
	}
}


// Convert a String to an int and return it
// If -1 returns, validation will fail
func stringToInt(x string) int {
      y, err := strconv.Atoi(x)
      if err != nil {
            CheckError(err, verbose)
            CronLog("Failed to convert int to string. Probably a bad GET.", verbose)
            return -1

      } else {
            return y
      }
}
