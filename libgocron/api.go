package libgocron
import (
	"errors"
	"net/http"
	"strings"
	"encoding/json"

	"gocron/util/log"
	"gocron/util/httphelper"

	"github.com/gorilla/mux"
)


// BackendAPI is a web service that exposes the backend to
// HTTP connections
func (g Gocron) Api(backendPort string) {
	log.Message("starting backend api on port: " + backendPort)
	r := mux.NewRouter()
    r.HandleFunc("/healthcheck", g.healthcheckAPI)
    r.HandleFunc("/version", g.backendVersionAPI)
    r.HandleFunc("/crons", g.getCronsAPI)
	http.ListenAndServe(":" + backendPort, r)
}


func (g Gocron) healthcheckAPI(resp http.ResponseWriter, req *http.Request) {
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


func (g Gocron) getCronsAPI(resp http.ResponseWriter, req *http.Request) {
    var a AllCrons

    rows, err := queryDatabase(g, "SELECT * FROM gocron;")
    if err != nil {
        log.Error(err)
		httphelper.ReturnServerError(resp, "", true)
        return
    }

    for rows.Next() {
        var c Cron
        rows.Scan(&c.Cronname,
            &c.Account,
            &c.Email,
            &c.Ipaddress,
            &c.Frequency,
            &c.Lastruntime,
            &c.Alerted,
            &c.Site)
        a.Crons = append(a.Crons, c)
    }
    a.Count = len(a.Crons)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(a)
	httphelper.ReturnOk(resp)
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
