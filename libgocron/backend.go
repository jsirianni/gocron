package libgocron
import (
	"os"
	"time"
	"errors"
	"net/http"
	"strings"
	"encoding/json"

	"gocron/util/log"
	"gocron/util/slack"
	"gocron/util/httphelper"
)


// StartBackend calls checkCronStatus() on a set interval
func (g Gocron) StartBackend(backendPort string) error {
	log.Message("gocron-back version: " + Version)

	// create the gocron table, if not exists
	err := g.createGocronTable()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// start the backend api on a new thread
	go g.BackendAPI(backendPort)

	// backend server is just a never ending loop that checks for missed
	// jobs at the set interval
	for {
		time.Sleep((time.Duration(g.Interval) * time.Second))
		log.Message("Checking for missed jobs.")
		g.cronStatus()
	}
}

// BackendAPI is a web service that exposes the backend to
// HTTP connections
func (g Gocron) BackendAPI(backendPort string) {
	log.Message("starting backend api on port: " + backendPort)

	http.HandleFunc("/healthcheck", g.backEndHealthCheck)
	http.HandleFunc("/version", g.backendVersionAPI)
	http.HandleFunc("/crons", g.getCrons)
	http.ListenAndServe(":" + backendPort, nil)
}

func (g Gocron) backEndHealthCheck(resp http.ResponseWriter, req *http.Request) {
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

func (g Gocron) backendVersionAPI(resp http.ResponseWriter, req *http.Request) {
	var b BackendVersion
	var err error
	b.Version = Version
	b.Database.Type = "postgres"
	b.Database.Version, err = g.getDatabaseVersion()
	if err != nil {
		log.Error(err)
		httphelper.ReturnServerError(resp, err.Error(), true)
	} else {
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(b)
		httphelper.ReturnOk(resp)
	}
}

func (g Gocron) getCrons(resp http.ResponseWriter, req *http.Request) {
	v := strings.Split(string(req.URL.Path), "/")
	var b []byte
	var err  error
	if len(v) == 1 {
		b, err = g.queryAllCrons("")
	} else if len(v) == 2 {
		b, err = g.queryAllCrons(v[1])
	}

	if err != nil {
		log.Error(err)
		httphelper.ReturnServerError(resp, "", true)
	} else {
		resp.Header().Set("Content-Type", "application/json")
		json.NewEncoder(resp).Encode(b)
		httphelper.ReturnOk(resp)
	}

}


// GetSummary prints a summary to standard out
func (g Gocron) GetSummary() {
	message := "gocron summary - missed jobs:\n"

	rows, err := queryDatabase(g, missedJobs)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query while attempting to build a summary: " + missedJobs))
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)

		message = message + "Name: " + cron.Cronname  + "| Account: " + cron.Account + "\n"
	}


	// Send slack alert
	err = g.slackAlert("gocron alert summary", message)
	if err != nil {
		log.Message("GOCRON: Failed to build alert summary.")
	} else {
		log.Message(message)
	}

}


func (g Gocron) cronStatus() {
	g.checkMissedJobs(missedJobs)
	g.checkRevivedJobs(revivedJobs)
}


func (g Gocron) checkMissedJobs(query string) {
	rows, err := queryDatabase(g, query)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query: " + query))
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)


		if cron.Alerted != true {
			subject := cron.Cronname + ": " + cron.Account + " failed to check in" + "\n"
			message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " has not checked in on time"

			// Only update database if alert sent successful
			if g.alert(cron, subject, message) == true {
				query = "UPDATE gocron SET alerted = true " +
					"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

				rows, err := queryDatabase(g, query)
				defer rows.Close()
				if err != nil {
					log.Error(err)
					log.Error(errors.New("Failed to update row for " + cron.Cronname))
				}
			}

		} else {
			log.Message("Alert for "+cron.Cronname+": "+cron.Account+
				" has been supressed. Already alerted")
		}
	}
}


func (g Gocron) checkRevivedJobs(query string) {
	rows, err := queryDatabase(g, query)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query: " + query))
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)

		query = "UPDATE gocron SET alerted = false " +
			"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

		rows, err := queryDatabase(g, query)
		defer rows.Close()
		if err != nil {
			log.Message("Failed to update row for " + cron.Cronname)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		g.alert(cron, subject, message)
	}
}


func (g Gocron) alert(cron Cron, subject string, message string) bool {

    // Immediately log the alert
    log.Message(subject)

	err := g.slackAlert(subject, message)
	if err != nil {
		log.Error(errors.New("gocron fail: alert for " + cron.Cronname))
		log.Error(err)
		return false
	}

    log.Message("gocron success: alert for " + cron.Cronname + " sent")
    return true
}


func (g Gocron) slackAlert(subject string, message string) error {
	var slack slack.Slack
	slack.HookURL      = g.SlackHookURL
	slack.Post.Channel = g.SlackChannel
	slack.Post.Text    = message
	return slack.Message()
}
