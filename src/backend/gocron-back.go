package main
import (
      "os"
      "time"
      "strings"
      "strconv"
      "gopkg.in/gomail.v2"
      "database/sql"; _ "github.com/lib/pq";
      "gocronlib"
)


const version string   = "2.0.0"
const selectAll string = "SELECT * FROM gocron;"

var verbose bool  = false    // Flag enabling / disabling verbosity
var checkInt int  = 300      // Time in seconds to check for missed jobs
var args []string = os.Args  // Command line arguments


func main() {
      // Parse arguments
      if len(os.Args) > 1 {

            // Return the current version
            if strings.Contains(args[1], "--version") {
                  println(version)
                  os.Exit(0)

            // When enabled, all logging will also print to screen
            } else if strings.Contains(args[1], "--verbose") {
                  verbose = true
                  gocronlib.CronLog("gocron started with --verbose.", verbose)
                  return

            } else {
                  return
            }
      }

      // Start the timer thread
      timer()
}


func timer() {
      for {
            // Check for missed jobs
            time.Sleep((time.Duration(checkInt) * time.Second))
            gocronlib.CronLog("Checking for missed jobs.", verbose)
            checkCronStatus()
      }
}


func checkCronStatus() {
      var subject string
      var message string

      db, err := sql.Open("postgres", gocronlib.DatabaseString(gocronlib.GetConfig(verbose)))
      if err != nil {
            gocronlib.CheckError(err, verbose)
      }
      defer db.Close()

      rows, err := db.Query(selectAll)
      if err != nil {
            gocronlib.CheckError(err, verbose)
      }
      defer rows.Close()

      for rows.Next() {
            var c gocronlib.Cron
            rows.Scan(&c.Cronname,
                        &c.Account,
                        &c.Email,
                        &c.Ipaddress,
                        &c.Frequency,
                        &c.Lastruntime,
                        &c.Alerted)

            var currentTime = int(time.Now().Unix())
            var lastRunTime, _ = strconv.Atoi(c.Lastruntime)
            var frequency, _ = strconv.Atoi(c.Frequency)

            // If not checked in on time
            if (currentTime - lastRunTime) > frequency {
                  subject = c.Cronname + ": " + c.Account + " failed to check in" + "\n"
                  message = "The cronjob " + c.Cronname + " for account " + c.Account +
                        " has not checked in on time"

                  // Mark row as alerted
                  if c.Alerted != true {
                        _, err = db.Exec("UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + c.Cronname + "' AND account = '" + c.Account + "';")
                        if err != nil {
                              gocronlib.CheckError(err, verbose)
                        }
                        alert(c, subject, message)
                        gocronlib.CronLog(subject, verbose)

                  // If alerted already marked true
                  } else {
                        gocronlib.CronLog("Alert for " + c.Cronname + ": " + c.Account +
                              " has been supressed. Already alerted", verbose)
                  }


            // If checked in on time but previously not (alerted == true)
            } else if ((currentTime - lastRunTime) < frequency) && c.Alerted == true {
                  _, err = db.Exec("UPDATE gocron SET alerted = false " +
                              "WHERE cronname = '" + c.Cronname + "' AND account = '" + c.Account + "';")
                  if err != nil {
                        gocronlib.CheckError(err, verbose)

                  } else {
                        subject = c.Cronname + ": " + c.Account + " is back online" + "\n"
                        message = "The cronjob " + c.Cronname + " for account " +
                              c.Account + " is back online"

                        alert(c, subject, message)
                        gocronlib.CronLog(subject, verbose)
                  }

            // Job in a good state
            } else {
                  subject = c.Cronname + ": " + c.Account + " is online" + "\n"
                  gocronlib.CronLog(subject, verbose)
            }
      }
}


func alert(c gocronlib.Cron, subject string, message string) {
      var config gocronlib.Config = gocronlib.GetConfig(verbose)
      var recipient string = c.Email
      var port, _ = strconv.Atoi(config.Smtpport)
      var d = gomail.NewDialer(config.Smtpserver, port, config.Smtpaddress, config.Smtppassword)
      var m = gomail.NewMessage()

      m.SetHeader("From", config.Smtpaddress)
      m.SetHeader("To", recipient)
      m.SetHeader("Subject", subject)
      m.SetBody("text/html", message)

      if err := d.DialAndSend(m); err != nil {
            gocronlib.CheckError(err, verbose)
      }

      gocronlib.CronLog("Alert for " + c.Cronname + " sent to " + recipient, verbose)
}
