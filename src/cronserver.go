package main
import (
      "fmt"
      "time"
      "strings"
      "strconv"
      "os/user"
      "net/http"
      "io/ioutil"
      "gopkg.in/yaml.v2"
      "database/sql"; _ "github.com/lib/pq";
)


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

var config Config


func main() {
      user, err := user.Current()

      yamlFile, err := ioutil.ReadFile(user.HomeDir + "/.config/gocron/.config.yml")
      if err != nil {
            panic(err)
      }

      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            panic(err)
      }

      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}



func cronStatus(w http.ResponseWriter, r *http.Request) {
      var currentTime = int(time.Now().Unix())
      var socket = strings.Split(r.RemoteAddr, ":")

      var cronJob Cron
      cronJob.cronName = r.URL.Query().Get("cronName")
      cronJob.account = r.URL.Query().Get("account")
      cronJob.email = r.URL.Query().Get("email")
      cronJob.ipAddress = socket[0]
      cronJob.cronTime = r.URL.Query().Get("time")
      cronJob.tolerance = r.URL.Query().Get("tolerance")
      cronJob.lastRunTime = strconv.Itoa(currentTime)


      go updateDatabase(cronJob)
}


func updateDatabase(c Cron) {
      var connectionString string
      var query string

      connectionString = "postgres://" +
      config.Dbuser + ":" +
      config.Dbpass + "@" +
      config.Dbfqdn +
      "/gocron" +
      "?sslmode=disable"

      query = "INSERT INTO gocron (cronname, account, email, ipaddress, crontime, tolerance, lastruntime) VALUES ('" +
             c.cronName + "','" +
             c.account + "','" +
             c.email + "','" +
             c.ipAddress + "','" +
             c.cronTime + "','" +
             c.tolerance + "','" +
             c.lastRunTime + "') " +
             "ON CONFLICT (cronname, account) DO UPDATE " +
             "SET lastruntime = " + "'" + c.lastRunTime + "'" +
             ";"

      go log("Cron update from " + c.account + " at " + c.ipAddress + "\n" +
      "Job: " + c.cronName + "\n" +
      "Time: " + c.lastRunTime + "\n" + query + "\n")

      db, err := sql.Open("postgres", connectionString)
      defer db.Close()
      if err != nil {
            checkError(err)
            panic(err)
      }
      _, err = db.Exec(query)
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
    log(err.Error())
  }
}

// log the output
func log(message string) {
      fmt.Println(message)
}
