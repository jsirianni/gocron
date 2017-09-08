package main

import (
      "fmt"
      "os/user"
      //"io"
      "net/http"
      "io/ioutil"
      "gopkg.in/yaml.v2"
      "database/sql"
      _ "github.com/lib/pq"
      "time"
      "strconv"
      "strings"
)

// Structs
type Config struct {
      Dbfqdn string
      Dbport string
      Dbuser string
      Dbpass string
      Dbdatabase string
      Smtpserver string
      Smtpport string
      Smtpaddress string
      Smtppassword string
}

type Cron struct {
      cronName string
      account string
      email string
      ipAddress string
      cronTime string
      tolerance string
      lastRunTime string  // Unix timestamp
}


// Global Vars
var config Config


func main() {
      // Read config file from user's home directory
      user, err := user.Current()
      yamlFile, err := ioutil.ReadFile(user.HomeDir + "/.config/gocron/.config.yml")
      if err != nil {
            panic(err)
      }
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            panic(err)
      }

      // Start HTTP server
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}



// Parse GET Request parameters
// TODO Add validation
func cronStatus(w http.ResponseWriter, r *http.Request) {
      // Build variables and assign values to cronJob
      var cronJob Cron
      var currentTime = int(time.Now().Unix())
      var socket = strings.Split(r.RemoteAddr, ":")
      cronJob.cronName = r.URL.Query().Get("cronName")
      cronJob.account = r.URL.Query().Get("account")
      cronJob.email = r.URL.Query().Get("email")
      cronJob.ipAddress = socket[0]
      cronJob.cronTime = r.URL.Query().Get("time")
      cronJob.tolerance = r.URL.Query().Get("tolerance")
      cronJob.lastRunTime = strconv.Itoa(currentTime)


      updateDatabase(cronJob)
}


// Insert or update a cron entry in the database
func updateDatabase(c Cron) {
      var connectionString = "postgres://" + config.Dbuser + ":" + config.Dbpass + "@" + config.Dbfqdn + "/gocron?sslmode=disable"


      // Check the database for the existing cron (primary key)
      // If the cron exists, update its lastRun column with the current time
      db, err := sql.Open("postgres", connectionString)
      if err != nil {
            panic(err)
      }

      _, err = db.Query("SELECT cronName FROM gocron WHERE account = " + c.account)

      // If the cron does not exist, create a table entry and record lastRun with the current time
      // Send an email alert notifying the user that the entry has been made

      db.Close()
}

// Check for missed cron updateDatabase
func checkCronStatus() {
      // TODO
      // Check the database for entries that have
      // not ran at their scheduled time + their tolerance
      //
      // Example: 1_*_*_*_* 30 should run at least every 1.5 hours
      //          Every hour with 30 minutes of tolerance
      //
      // Example: 0_19_*_*_* 120 should run at 7pm every day
      //          with 2 hours of tolerance (7-9pm)
      //
      // TODO Run this function every 10 minutes ??
      //
      // Send email alerts for any entries that have not checked in on time
}

// Send emails
func alert(recipient string, subject string, message string) {
      // TODO
      // Send an email alert

      // TODO
      // Add optional slack alerts
}


// Handle errors that do not require the program to stop
func checkError(err error) {
  if err != nil {
    fmt.Printf("ERROR: " + err.Error() + "\n\n")
  }
}
