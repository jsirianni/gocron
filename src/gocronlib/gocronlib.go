package gocronlib
import (
      "os"
      "io/ioutil"
      "fmt"
      "os/exec"
      "strconv"
      "gopkg.in/yaml.v2"
      "database/sql"; _ "github.com/lib/pq";
)


const Version string  = "2.0.4"

const (
      sslmode  string = "disable"   // Disable or enable ssl
      syslog   string = "logger"    // Command to write to syslog
      confPath string = "/etc/gocron/config.yml"
)

type Config struct {
      Dbfqdn       string `yaml:"dbfqdn"`
      Dbport       string `yaml:"dbport"`
      Dbuser       string `yaml:"dbuser"`
      Dbpass       string `yaml:"dbpass"`
      Dbdatabase   string `yaml:"dbdatabase"`
      Smtpserver   string `yaml:"smtpserver"`
      Smtpport     string `yaml:"smtpport"`
      Smtpaddress  string `yaml:"smtpaddress"`
      Smtppassword string `yaml:"smtppassword"`
      Interval     int    `yaml:"interval"`
      SlackHookUrl string `yaml:"slackhookurl"`
      SlackChannel string `yaml:"slackchannel"`
      PreferSlack  bool   `yaml:"preferslack"`
}


type Cron struct {
      Cronname    string   // Name of the cronjob
      Account     string   // Account the job belongs to
      Email       string   // Address to send alerts to
      Ipaddress   string   // Source IP address
      Frequency   int      // How often a job should check in
      Lastruntime int      // Unix timestamp
      Alerted     bool     // set to true if an alert has already been thrown
      Site        bool     // Set true if service is a site (Example: Network gateway)
}


// Read in the config file
func GetConfig(verbose bool) Config {
      var config Config
      yamlFile, err := ioutil.ReadFile(confPath)
      if err != nil {
           CheckError(err, verbose)
           os.Exit(1)
      }

      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            CheckError(err, verbose)
            os.Exit(1)
      }

      return config
}


// Return a Postgres connection string
func DatabaseString(verbose bool) string {
      var c Config = GetConfig(verbose)
      return "postgres://" + c.Dbuser + ":" + c.Dbpass + "@" + c.Dbfqdn + "/gocron" + "?sslmode=" + sslmode
}


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


// Convert a String to an int and return it
// If -1 returns, validation will fail
func StringToInt(x string, verbose bool) int {
      y, err := strconv.Atoi(x)
      if err != nil {
            CheckError(err, verbose)
            CronLog("Failed to convert int to string. Probably a bad GET.", verbose)
            return -1

      } else {
            return y
      }
}


// Function writes messages to syslog and (optionally) to standard out
func CronLog(message string, verbose bool) {
      err := exec.Command(syslog, message).Run()
      if err != nil {
            fmt.Println("Failed to write to syslog")
            fmt.Println(message)
      }
      if verbose == true {
            fmt.Println(message)
      }
}


// Function passes error messages to the cronLog() function
func CheckError(err error, verbose bool) {
      if err != nil {
            CronLog("Error: \n" + err.Error(), verbose)
      }
}
