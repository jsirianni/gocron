package cmd
import (
    "os"
    "strconv"

	"database/sql"; _ "github.com/lib/pq";
)


// Function handles database queries
// Returns false if bad query
func QueryDatabase(query string, verbose bool) (*sql.Rows, bool) {
      var (
            db *sql.DB
            rows *sql.Rows
            err error
            status bool
      )

      db, err = sql.Open("postgres", DatabaseString(verbose))
      defer db.Close()
      if err != nil {
            CheckError(err, verbose)
      }

      rows, err = db.Query(query)
      if err != nil {
            CheckError(err, verbose)
            status = false
      } else {
            status = true
      }

      // Return query result and status
      return rows, status
}


// Return a Postgres connection string
func DatabaseString(verbose bool) string {
      return "postgres://" + config.Dbuser + ":" + config.Dbpass + "@" + config.Dbfqdn + "/gocron" + "?sslmode=" + "disable"
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
	rows, result := QueryDatabase(query, verbose)
	defer rows.Close()
	if result == true {
		CronLog("Heartbeat from "+c.Cronname+": "+c.Account+" \n", verbose)
		return true

	} else {
		return false
	}
}


// Creates the gocron database table, if it does not exist
// Returns false if not successful, else true
func CreateGocronTable(verbose bool) bool {
    query := "CREATE TABLE IF NOT EXISTS gocron(cronName varchar, " +
        "account varchar, email varchar, ipaddress varchar, " +
        "frequency integer, lastruntime integer, alerted boolean, " +
        "site boolean, PRIMARY KEY(cronname, account));"
    _, result := QueryDatabase(query, verbose)
    if result == false {
        CronLog("Table 'gocron' is missing. Creation failed. Validate permissions in the config.", verbose)
        os.Exit(1)
    }
    return result
}
