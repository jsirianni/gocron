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
      "gopkg.in/gomail.v2"
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
      cronname string
      account string
      email string
      ipaddress string
      frequency string
      tolerance string
      lastruntime string  // Unix timestamp
      alerted bool        // set to true if an alert has already been thrown
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

      go timer()

      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}


func timer() {
      for {
            time.Sleep((20 * time.Second))
            log("Checking for missed jobs.")
            go checkCronStatus()
      }
}


func cronStatus(w http.ResponseWriter, r *http.Request) {
      var cronJob Cron
      var currentTime = int(time.Now().Unix())
      var socket = strings.Split(r.RemoteAddr, ":")

      cronJob.cronname = r.URL.Query().Get("cronname")
      cronJob.account = r.URL.Query().Get("account")
      cronJob.email = r.URL.Query().Get("email")
      cronJob.ipaddress = socket[0]
      cronJob.frequency = r.URL.Query().Get("frequency")
      cronJob.tolerance = r.URL.Query().Get("tolerance")
      cronJob.lastruntime = strconv.Itoa(currentTime)

      go updateDatabase(cronJob)
}


func updateDatabase(c Cron) {
      var query string
      query = "INSERT INTO gocron (cronname, account, email, ipaddress, frequency, tolerance, lastruntime) VALUES ('" +
             c.cronname + "','" +
             c.account + "','" +
             c.email + "','" +
             c.ipaddress + "','" +
             c.frequency + "','" +
             c.tolerance + "','" +
             c.lastruntime + "') " +
             "ON CONFLICT (cronname, account) DO UPDATE " +
             "SET email = " + "'" + c.email + "'," +
             "ipaddress = " + "'" + c.ipaddress + "'," +
             "frequency = " + "'" + c.frequency + "'," +
             "tolerance = " + "'" + c.tolerance + "'," +
             "lastruntime = " + "'" + c.lastruntime + "'" +
             ";"

      go log("Cron update from " + c.account + " at " + c.ipaddress + "\n" +
            "Job: " + c.cronname + "\n" +
            "Time: " + c.lastruntime + "\n" + query)

      db, err := sql.Open("postgres", databaseString())
      defer db.Close()
      if err != nil {
            checkError(err)
            panic(err)
      }

      _, err = db.Exec(query)
      if err != nil {
            checkError(err)
            panic(err)
      }
}


func checkCronStatus() {
      db, err := sql.Open("postgres", databaseString())
      defer db.Close()
      if err != nil {
            checkError(err)
            panic(err)
      }

      rows, err := db.Query("SELECT * FROM gocron;")
      defer rows.Close()
      if err != nil {
            checkError(err)
      }

      for rows.Next() {
            var c Cron
            rows.Scan(&c.cronname, &c.account, &c.email, &c.ipaddress, &c.frequency, &c.tolerance, &c.lastruntime, &c.alerted)

            var currentTime = int(time.Now().Unix())
            var lastRunTime, _ = strconv.Atoi(c.lastruntime)
            var frequency, _ = strconv.Atoi(c.frequency)
            var tolerance, _ = strconv.Atoi(c.tolerance)
            var maxTime = frequency + tolerance

            if (currentTime - lastRunTime) > maxTime {
                  log(c.cronname + " for account " + c.account + " has not checked in on time")
                  if c.alerted != true {
                        alert(c.email, c)
                        db.Exec("UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")

                  } else {
                        log("Alert for " + c.cronname + ": " + c.account + " has been supressed. Already alerted." )
                  }

            } else {
                  log("Job: " + c.cronname + ": " + c.account + " has checked in recently.")
                  db.Exec("UPDATE gocron SET alerted = false " +
                          "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
            }
      }
}


func alert(recipient string, c Cron) {
      var port, _ = strconv.Atoi(config.Smtpport)
      var d = gomail.NewDialer(config.Smtpserver, port, config.Smtpaddress, config.Smtppassword)
      var m = gomail.NewMessage()
      var subject = "Cron failed to run: " + c.cronname + "\n"
      var message = "The cronjob " + c.cronname + " for account " + c.account + " has not checked in on time."

      m.SetHeader("From", config.Smtpaddress)
      m.SetHeader("To", recipient)
      m.SetHeader("Subject", subject)
      m.SetBody("text/html", message)

      if err := d.DialAndSend(m); err != nil {
            checkError(err)
      }

      log("Alert for " + c.cronname + " sent to " + recipient)

      // TODO
      // Add optional slack alerts
}


func checkError(err error) {
  if err != nil {
    log("Error: \n" + err.Error())
  }
}


func log(message string) {
      fmt.Println("\n" + message)
}


func databaseString() string {
      var connectionString string
      connectionString = "postgres://" +
      config.Dbuser + ":" +
      config.Dbpass + "@" +
      config.Dbfqdn +
      "/gocron" +
      "?sslmode=disable"

      return connectionString
}
