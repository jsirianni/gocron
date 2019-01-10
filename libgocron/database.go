package libgocron
import (
    "os"
    "strconv"

	"database/sql"; _ "github.com/lib/pq";
)


const missedJobs = "SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;"
const revivedJobs = "SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;"


// Function handles database queries
// Returns false if bad query
func queryDatabase(query string) (*sql.Rows, bool) {
      var (
            db *sql.DB
            rows *sql.Rows
            err error
            status bool
      )

      db, err = sql.Open("postgres", "postgres://" + config.Dbuser + ":" + config.Dbpass + "@" + config.Dbfqdn + "/gocron" + "?sslmode=" + "disable")
      defer db.Close()
      if err != nil {
            CheckError(err)
      }

      rows, err = db.Query(query)
      if err != nil {
            CheckError(err)
            status = false
      } else {
            status = true
      }

      // Return query result and status
      return rows, status
}


func updateDatabase(c Cron) bool {
	var (
		query  string
		result bool

		frequency   string = strconv.Itoa(c.Frequency)
		lastruntime string = strconv.Itoa(c.Lastruntime)
		site        string = strconv.FormatBool(c.Site)
	)

	// Insert and update if already exist
	query = "INSERT INTO gocron " +
		"(cronname, account, email, ipaddress, frequency, lastruntime, alerted, site) " +
		"VALUES ('" +
		c.Cronname + "','" + c.Account + "','" + c.Email + "','" + c.Ipaddress + "','" +
		frequency + "','" + lastruntime + "','" + "false" + "','" + site + "') " +
		"ON CONFLICT (cronname, account) DO UPDATE " +
		"SET email = " + "'" + c.Email + "'," + "ipaddress = " + "'" + c.Ipaddress + "'," +
		"frequency = " + "'" + frequency + "'," + "lastruntime = " + "'" + lastruntime + "', " +
		"site = " + "'" + site + "';"

	// Execute query
	rows, result := queryDatabase(query)
	defer rows.Close()
	if result == true {
		CronLog("Heartbeat from "+c.Cronname+": "+c.Account+" \n")
		return true

	} else {
		return false
	}
}


// Creates the gocron database table, if it does not exist
// Returns false if not successful, else true
func createGocronTable() bool {
    query := "CREATE TABLE IF NOT EXISTS gocron(cronName varchar, " +
        "account varchar, email varchar, ipaddress varchar, " +
        "frequency integer, lastruntime integer, alerted boolean, " +
        "site boolean, PRIMARY KEY(cronname, account));"
    _, result := queryDatabase(query)
    if result == false {
        CronLog("Table 'gocron' is missing. Creation failed. Validate permissions in the config.")
        os.Exit(1)
    }
    return result
}
