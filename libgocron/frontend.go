package libgocron
import (
	"net/http"
    "strconv"
	"strings"
	"time"
	"io/ioutil"
    "encoding/json"
	"errors"

	"gocron/util/log"
)


// StartFrontend starts the gocron frontend server
func (g Gocron) StartFrontend(frontendPort string) {

	//if v == true {
	//	CronLog("verbose mode enabled")
	log.CronLog("gocron-front version: " + Version)
	log.CronLog("starting web server on port: " + frontendPort)
	//}

	http.HandleFunc("/", g.incomingCron)
	http.HandleFunc("/healthcheck", frontEndHealthCheck)
	http.ListenAndServe(":"+frontendPort, nil)
}


// return http status 200 if connection to database is healthy
func frontEndHealthCheck(resp http.ResponseWriter, req *http.Request) {
    r := strings.Split(req.RemoteAddr, ":")[0]
	//if verbose == true {
	log.CronLog("healthcheck from: " + r)
	//}
	returnOk(resp)
}


// Validate the request and then pass to updateDatabase()
func (g Gocron) incomingCron(resp http.ResponseWriter, req *http.Request) {
	var (
		currentTime int = int(time.Now().Unix())
		socket          = strings.Split(req.RemoteAddr, ":")
		c           Cron
		method      string
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
		 	log.CronLog(err.Error())
		}

		if err := json.Unmarshal(payload, &c); err != nil {
			log.CronLog(err.Error())
		}
		c.Lastruntime = currentTime
		c.Ipaddress = socket[0]


	default:
		// Log an error and do not respond
		log.CronLog("Incoming request from "+c.Ipaddress+" is not a GET or POST.")
		return
	}

	if c.ValidateParams() == true {
		if g.updateDatabase(c) == true {
			returnCreated(resp)

		} else {
			returnServerError(resp)
		}

	} else {
		returnNotFound(resp)
		log.CronLog(method+" from "+c.Ipaddress+" not valid. Dropping.")
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


// ValidateParams validates SQL variables
func (c Cron) ValidateParams() bool {

	valid  := false // Flag determines the return value

	if c.CheckLength() == true {
		valid = true
	}

	/*if verbose == true {
		if valid == true {
			CronLog("Parameters from "+c.Ipaddress+" passed validation")
			return true
		}

		CronLog("Parameters from "+c.Ipaddress+" failed validation!")
		return false
	}*/

	return valid
}


// CheckLength validates that parameters are present
func (c Cron) CheckLength() bool {
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
        log.LogError(err)
        log.LogError(errors.New("failed to convert int to string. Probably a bad GET"))
        return -1
    }

    return y
}
