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
	"gocron/util/httphelper"
)


// StartFrontend starts the gocron frontend server
func (g Gocron) StartFrontend(frontendPort string) {
	log.Message("gocron-front version: " + Version)
	log.Message("starting web server on port: " + frontendPort)

	http.HandleFunc("/", g.incomingCron)
	http.HandleFunc("/healthcheck", g.frontEndHealthCheck)
	http.ListenAndServe(":"+frontendPort, nil)
}


// return http status 200 if connection to database is healthy
func (g Gocron) frontEndHealthCheck(resp http.ResponseWriter, req *http.Request) {
    r := strings.Split(req.RemoteAddr, ":")[0]
	log.Message("healthcheck from: " + r)
	err := g.testDatabaseConnection()
	if err != nil {
		log.Error(err)
		httphelper.ReturnServerError(resp, "a connection to the database could not be validated", true)
	} else {
		httphelper.ReturnOk(resp)
	}
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

	case "POST":
		method = "POST"

		// read the request body into a byte array
		payload, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
		 	log.Message(err.Error())
		}

		if err := json.Unmarshal(payload, &c); err != nil {
			log.Message(err.Error())
		}
		c.Lastruntime = currentTime
		c.Ipaddress = socket[0]


	default:
		// Log an error and do not respond
		log.Message("Incoming request from "+c.Ipaddress+" is not a GET or POST.")
		return
	}

	if c.ValidateParams() == true {
		if g.updateDatabase(c) == true {
			httphelper.ReturnCreated(resp)

		} else {
			httphelper.ReturnServerError(resp, "failed to update the database", false)
		}

	} else {
		httphelper.ReturnNotFound(resp)
		log.Message(method+" from "+c.Ipaddress+" not valid. Dropping.")
	}
}


// ValidateParams validates SQL variables
func (c Cron) ValidateParams() bool {
	if err := c.CheckLength(); err != nil {
		return false
	}
	return true
}


// CheckLength validates that parameters are present
func (c Cron) CheckLength() error {
	var m string
	if len(c.Account) == 0 {
		m = m + "Account, "
	}
	if len(c.Cronname) == 0 {
		m = m + "Cronname, "
	}
	if len(c.Email) == 0 {
		m = m + "Email, "
	}
	if c.Frequency < 1 {
		m = m + "Frequency, "
	}
	if len(c.Ipaddress) == 0 {
		m = m + "Ipaddress, "
	}
	if c.Lastruntime == -1 {
		m = m + "Lastruntime"
	}

	if len(m) != 0 {
		return errors.New("Length check failed for parameters: ")
	}
	return nil
}


// Convert a String to an int and return it
// If -1 returns, validation will fail
func stringToInt(x string) int {
    y, err := strconv.Atoi(x)
    if err != nil {
        log.Error(err)
        log.Error(errors.New("failed to convert int to string. Probably a bad GET"))
        return -1
    }

    return y
}
