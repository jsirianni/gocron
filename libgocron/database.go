package libgocron
import (
    "os"
    "strconv"
    "errors"

    "gocron/util/log"

	"database/sql"
    _ "github.com/lib/pq" // github.com/lib/pq is required by database/sql

)


const missedJobs = "SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;"
const revivedJobs = "SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;"
const testConnection = "SELECT max(table_catalog) as x from information_schema.tables;"

// Function ensures a connection to the database. Errors
// should be logged by the caller.
func (g Gocron) testDatabaseConnection() error {
    _, err := queryDatabase(g, testConnection)
    return err
}

// Function handles database queries
func queryDatabase(g Gocron, query string) (*sql.Rows, error) {
    conn := "postgres://" + g.Dbuser + ":" + g.Dbpass + "@" + g.Dbfqdn + "/gocron" + "?sslmode=" + "disable"
    db, err := sql.Open("postgres", conn)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    return db.Query(query)
}


func (g Gocron) updateDatabase(c Cron) bool {
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
	rows, err := queryDatabase(g, query)
	if err != nil {
        log.Error(err)
        return false
	}
    defer rows.Close()

    log.Message("Heartbeat from "+c.Cronname+": "+c.Account+" \n")
    return true
}


// Creates the gocron database table, if it does not exist
func (g Gocron) createGocronTable() error {
    query := "CREATE TABLE IF NOT EXISTS gocron(cronName varchar, " +
        "account varchar, email varchar, ipaddress varchar, " +
        "frequency integer, lastruntime integer, alerted boolean, " +
        "site boolean, PRIMARY KEY(cronname, account));"
    _, err := queryDatabase(g, query)
    if err != nil {
        log.Error(err)
        log.Error(errors.New("table 'gocron' is missing. Creation failed. Validate permissions in the config"))
        os.Exit(1)
    }

    return err
}

// returns the database version
func (g Gocron) getDatabaseVersion() (string, error) {
    var err error
    var v   string

    db, err := sql.Open("postgres", "user=" + g.Dbuser + " password=" + g.Dbpass + " dbname=" + g.Dbdatabase + " sslmode=disable")
    if err != nil {
      return "", err
    }

    err = db.QueryRow("SELECT version();").Scan(&v)
    if err != nil {
        return "", errors.New(err.Error())
    }
    return v, nil
}
