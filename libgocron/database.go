package libgocron
import (
    "os"
    "strconv"
    "errors"


	"database/sql"
    _ "github.com/lib/pq" // github.com/lib/pq is required by database/sql

)


const missedJobs = "SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;"
const revivedJobs = "SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;"


// Function handles database queries
func queryDatabase(query string) (*sql.Rows, error) {
    conn := "postgres://" + config.Dbuser + ":" + config.Dbpass + "@" + config.Dbfqdn + "/gocron" + "?sslmode=" + "disable"
    db, err := sql.Open("postgres", conn)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    return db.Query(query)
}


func updateDatabase(c Cron) bool {
	frequency   := strconv.Itoa(c.Frequency)
	lastruntime := strconv.Itoa(c.Lastruntime)
	site        := strconv.FormatBool(c.Site)


	// Insert and update if already exist
	query := "INSERT INTO gocron " +
		"(cronname, account, email, ipaddress, frequency, lastruntime, alerted, site) " +
		"VALUES ('" +
		c.Cronname + "','" + c.Account + "','" + c.Email + "','" + c.Ipaddress + "','" +
		frequency + "','" + lastruntime + "','" + "false" + "','" + site + "') " +
		"ON CONFLICT (cronname, account) DO UPDATE " +
		"SET email = " + "'" + c.Email + "'," + "ipaddress = " + "'" + c.Ipaddress + "'," +
		"frequency = " + "'" + frequency + "'," + "lastruntime = " + "'" + lastruntime + "', " +
		"site = " + "'" + site + "';"

	// Execute query
	rows, err := queryDatabase(query)
	if err != nil {
        LogError(err)
        return false
	}
    defer rows.Close()

    CronLog("Heartbeat from "+c.Cronname+": "+c.Account+" \n")
    return true
}


// Creates the gocron database table, if it does not exist
func createGocronTable() error {
    query := "CREATE TABLE IF NOT EXISTS gocron(cronName varchar, " +
        "account varchar, email varchar, ipaddress varchar, " +
        "frequency integer, lastruntime integer, alerted boolean, " +
        "site boolean, PRIMARY KEY(cronname, account));"
    _, err := queryDatabase(query)
    if err != nil {
        LogError(err)
        LogError(errors.New("Table 'gocron' is missing. Creation failed. Validate permissions in the config."))
        os.Exit(1)
    }

    return err
}
