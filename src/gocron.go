package main
import (
      "fmt"
      "time"
      "strings"
      "strconv"
      "os/user"
      "os/exec"
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
            checkError(err)
      }
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            checkError(err)
      }

      go timer()

      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}


func timer() {
      for {
            time.Sleep((120 * time.Second))
            cronLog("Checking for missed jobs.")
            go checkCronStatus()
      }
}


func cronStatus(w http.ResponseWriter, r *http.Request) {
      var currentTime int = int(time.Now().Unix())
      var socket = strings.Split(r.RemoteAddr, ":")
      var cronJob Cron

      cronJob.cronname = r.URL.Query().Get("cronname")
      cronJob.account = r.URL.Query().Get("account")
      cronJob.email = r.URL.Query().Get("email")
      cronJob.frequency = r.URL.Query().Get("frequency")
      cronJob.tolerance = r.URL.Query().Get("tolerance")
      cronJob.lastruntime = strconv.Itoa(currentTime)
      cronJob.ipaddress = socket[0]

      if checkLength(cronJob) == true {
            go updateDatabase(cronJob)

      } else {
            cronLog("GET request not valid. Dropping.")
      }
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

      db, err := sql.Open("postgres", databaseString())
      defer db.Close()
      if err != nil {
            checkError(err)
      }

      _, err = db.Exec(query)
      if err != nil {
            checkError(err)
      }
}


func checkCronStatus() {
      db, err := sql.Open("postgres", databaseString())
      if err != nil {
            checkError(err)
      }
      defer db.Close()

      rows, err := db.Query("SELECT * FROM gocron;")
      if err != nil {
            checkError(err)
      }
      defer rows.Close()

      for rows.Next() {
            var c Cron
            rows.Scan(&c.cronname,
                        &c.account,
                        &c.email,
                        &c.ipaddress,
                        &c.frequency,
                        &c.tolerance,
                        &c.lastruntime,
                        &c.alerted)

            var currentTime = int(time.Now().Unix())
            var lastRunTime, _ = strconv.Atoi(c.lastruntime)
            var frequency, _ = strconv.Atoi(c.frequency)
            var tolerance, _ = strconv.Atoi(c.tolerance)
            var maxTime = frequency + tolerance

            if (currentTime - lastRunTime) > maxTime {
                  cronLog(c.cronname + " for account " + c.account + " has not checked in on time")
                  if c.alerted != true {
                        alert(c.email, c)
                        _, err = db.Exec("UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
                        if err != nil {
                              checkError(err)
                        }

                  } else {
                        cronLog("Alert for " + c.cronname + ": " + c.account + " has been supressed. Already alerted" )
                  }

            } else {
                  cronLog("Job: " + c.cronname + ": " + c.account + " has checked in recently.")
                  _, err = db.Exec("UPDATE gocron SET alerted = false " +
                          "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
                  if err != nil {
                        checkError(err)
                  }
            }
      }
}


func alert(recipient string, c Cron) {
      var port, _ = strconv.Atoi(config.Smtpport)

      var d = gomail.NewDialer(config.Smtpserver,
                                    port,
                                    config.Smtpaddress,
                                    config.Smtppassword)

      var subject = "Cron failed to run: " + c.cronname + "\n"
      var message = "The cronjob " + c.cronname + " for account " + c.account + " has not checked in on time"

      var m = gomail.NewMessage()
      m.SetHeader("From", config.Smtpaddress)
      m.SetHeader("To", recipient)
      m.SetHeader("Subject", subject)
      m.SetBody("text/html", message)

      if err := d.DialAndSend(m); err != nil {
            checkError(err)
      }

      cronLog("Alert for " + c.cronname + " sent to " + recipient)
}


func checkError(err error) {
  if err != nil {
    cronLog("Error: \n" + err.Error())
  }
}


func cronLog(message string) {
      err := exec.Command("logger", message).Run()
      if err != nil {
            fmt.Println("Failed to write to syslog")
            fmt.Println(message)
      }
}


func databaseString() string {
      var connectionString string = "postgres://" +
            config.Dbuser + ":" +
            config.Dbpass + "@" +
            config.Dbfqdn +
            "/gocron" +
            "?sslmode=disable"

      return connectionString
}


func checkLength(c Cron) bool {
      if len(c.account) == 0 {
            cronLog("Account is not valid")
            return false

      } else if len(c.cronname) == 0 {
            cronLog("Cronname is not valid")
            return false

      } else if len(c.email) == 0 {
            cronLog("Email is not valid")
            return false

      } else if len(c.frequency) == 0 {
            cronLog("Frequency is not valid")
            return false

      } else if len(c.ipaddress) == 0 {
            cronLog("IP Address is not valid")
            return false

      } else if len(c.lastruntime) == 0 {
            cronLog("Runtime is not valid")
            return false

      } else if len(c.tolerance) == 0 {
            cronLog("Tolerance is not valid")
            return false

      } else {
            cronLog("Validation passed.")
            return true
      }
}
