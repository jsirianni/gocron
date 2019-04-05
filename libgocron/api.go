package libgocron
import (
	"errors"
	"net/http"
	"strings"
	"encoding/json"
	"database/sql"
	_ "expvar" // expvar will expose metrics without being called

	"gocron/util/log"
	"gocron/util/httphelper"

	"github.com/gorilla/mux"
)


// Api runs gocron's rest api
func (g Gocron) Api(backendPort string) {
	log.Message("starting backend api on port: " + backendPort)

	r := mux.NewRouter()

	// mux routes for gocron rest api
	r.HandleFunc("/healthcheck", g.healthcheckAPI).Methods("GET")
    r.HandleFunc("/version", g.versionAPI).Methods("GET")
	r.HandleFunc("/cron/missed", g.getCronsMissedAPI).Methods("GET")
	r.HandleFunc("/cron/{account}", g.getCronsByAccountAPI).Methods("GET")
	r.HandleFunc("/cron", g.getCronsAPI).Methods("GET")

	// expvar runtime  metrics
	r.Handle("/debug/vars", http.DefaultServeMux)

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


func (g Gocron) versionAPI(resp http.ResponseWriter, req *http.Request) {
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

    a.Crons = getCronsFromRows(rows)
    a.Count = len(a.Crons)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(a)
	httphelper.ReturnOk(resp)
}

func (g Gocron) getCronsByAccountAPI(resp http.ResponseWriter, req *http.Request) {
	var a AccountCrons
	a.Account = getAccountFromRequest(req)

    rows, err := queryDatabase(g, "SELECT * FROM gocron WHERE account = '" + a.Account + "';")
    if err != nil {
        log.Error(err)
		httphelper.ReturnServerError(resp, "", true)
        return
    }

	a.Crons = getCronsFromRows(rows)
    a.Count = len(a.Crons)

	if a.Count < 1 {
		httphelper.ReturnNotFound(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(a)
	httphelper.ReturnOk(resp)
}

func (g Gocron) getCronsMissedAPI(resp http.ResponseWriter, req *http.Request) {
	var a AllCrons

	rows, err := queryDatabase(g, missedJobs)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query while attempting to build a summary: " + missedJobs))
		httphelper.ReturnServerError(resp, "", true)
		return
	}

	a.Crons = getCronsFromRows(rows)
	a.Count = len(a.Crons)

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(a)
	httphelper.ReturnOk(resp)
}

func getAccountFromRequest(req *http.Request) string {
   vars := mux.Vars(req)
   return  vars["account"]
}

func getCronsFromRows(r *sql.Rows) []Cron {
	var crons []Cron
	for r.Next() {
		var c Cron
		r.Scan(&c.Cronname,
			&c.Account,
			&c.Email,
			&c.Ipaddress,
			&c.Frequency,
			&c.Lastruntime,
			&c.Alerted,
			&c.Site)
		crons = append(crons, c)
	}
	return crons
}
