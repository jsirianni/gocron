package main
import (
      "os"
      "io/ioutil"
      "fmt"
      "os/exec"
      "gopkg.in/yaml.v2"
)


type Config struct {
      Dbfqdn       string
      Dbport       string
      Dbuser       string
      Dbpass       string
      Dbdatabase   string
      Smtpserver   string
      Smtpport     string
      Smtpaddress  string
      Smtppassword string
}


type Cron struct {
      cronname    string   // Name of the cronjob
      account     string   // Account the job belongs to
      email       string   // Address to send alerts to
      ipaddress   string   // Source IP address
      frequency   string   // How often a job should check in
      lastruntime string   // Unix timestamp
      alerted     bool     // set to true if an alert has already been thrown
}


const sslmode string  = "disable"   // Disable or enable ssl
const syslog string   = "logger"    // Command to write to syslog
const confPath string = "/etc/gocron/config.yml"


// Read in the config file
func getConfig() Config {
      var config Config
      yamlFile, err := ioutil.ReadFile(confPath)
      if err != nil {
            checkError(err, verbose)
            os.Exit(1)
      }

      // Set the global config var
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            checkError(err, verbose)
            os.Exit(1)
      }

      return config
}


// Return a Postgres connection string
func databaseString(config Config) string {
      var connectionString string = "postgres://" +
            config.Dbuser + ":" +
            config.Dbpass + "@" +
            config.Dbfqdn +
            "/gocron" +
            "?sslmode=" + sslmode

      return connectionString
}


// Function writes messages to syslog and (optionally) to standard out
func cronLog(message string, verbose bool) {
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
func checkError(err error, verbose bool) {
      if err != nil {
            cronLog("Error: \n" + err.Error(), verbose)
      }
}
