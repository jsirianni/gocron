package main
import (
      "time"
      "strconv"
      "gopkg.in/gomail.v2"
      "database/sql"; _ "github.com/lib/pq";
)


func timer() {
      for {
            // Check for missed jobs every five minutes
            time.Sleep((300 * time.Second))
            cronLog("Checking for missed jobs.")
            checkCronStatus()
      }
}


func checkCronStatus() {
      var subject string
      var message string

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
            rows.Scan(&c.cronname, &c.account,
                        &c.email, &c.ipaddress,
                        &c.frequency, &c.tolerance,
                        &c.lastruntime, &c.alerted)

            var currentTime = int(time.Now().Unix())
            var lastRunTime, _ = strconv.Atoi(c.lastruntime)
            var frequency, _ = strconv.Atoi(c.frequency)
            var tolerance, _ = strconv.Atoi(c.tolerance)
            var maxTime = frequency + tolerance

            // If not checked in on time
            if (currentTime - lastRunTime) > maxTime {
                  subject = c.cronname + ": " + c.account + " failed to check in" + "\n"
                  message = "The cronjob " + c.cronname + " for account " + c.account + " has not checked in on time"

                  // Mark row as alerted
                  if c.alerted != true {
                        _, err = db.Exec("UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
                        if err != nil {
                              checkError(err)
                        }
                        alert(c.email, c, subject, message)
                        cronLog(subject)

                  // If alerted already marked true
                  } else {
                        cronLog("Alert for " + c.cronname + ": " + c.account + " has been supressed. Already alerted" )
                  }

            // If checked in on time but previously not (alerted == true)
            } else if ((currentTime - lastRunTime) < maxTime) && c.alerted == true {
                  _, err = db.Exec("UPDATE gocron SET alerted = false " +
                              "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
                  if err != nil {
                        checkError(err)

                  } else {
                        subject = c.cronname + ": " + c.account + " is back online" + "\n"
                        message = "The cronjob " + c.cronname + " for account " + c.account + " is back online"
                        alert(c.email, c, subject, message)
                        cronLog(subject)
                  }

            // Job in a good state
            } else {
                  // Set alerted to false if null
                  if c.alerted != true && c.alerted != false  {
                        _, err = db.Exec("UPDATE gocron SET alerted = false " +
                                    "WHERE cronname = '" + c.cronname + "' AND account = '" + c.account + "';")
                        if err != nil {
                              checkError(err)
                        }
                  }

                  subject = c.cronname + ": " + c.account + " is online" + "\n"
                  cronLog(subject)
            }
      }
}


func alert(recipient string, c Cron, subject string, message string) {
      var port, _ = strconv.Atoi(config.Smtpport)
      var d = gomail.NewDialer(config.Smtpserver, port, config.Smtpaddress, config.Smtppassword)
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
