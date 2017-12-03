package main
import (
      "os"
      "time"
      "strings"
      "strconv"
      "gopkg.in/gomail.v2"
      "gocronlib"
)


const version string    = "2.0.1"
const libVersion string = gocronlib.Version

var verbose bool  = false    // Flag enabling / disabling verbosity
var checkInt int  = 60       // Time in seconds to check for missed jobs
var args []string = os.Args  // Command line arguments


func main() {
      // Parse arguments
      if len(os.Args) > 1 {
            // Return the current version
            if strings.Contains(args[1], "--version") {
                  println("gocron-front version: " + version)
                  println("gocronlib version: " + libVersion)
                  os.Exit(0)
            }
            // When enabled, all logging will also print to screen
            if strings.Contains(args[1], "--verbose") {
                  verbose = true
                  gocronlib.CronLog("gocron started with --verbose.", verbose)
            }
      }

      // Run the timer
      timer()
}


// Function calls checkCronStatus() on a set interval
func timer() {
      for {
            time.Sleep((time.Duration(checkInt) * time.Second))
            gocronlib.CronLog("Checking for missed jobs.", verbose)
            checkCronStatus()
      }
}


func checkCronStatus() {
      var subject string  // Subject used in alerts
      var message string  // Message used in alerts
      var result bool     // Handles InserDatabase responses
      var query string    // Queries to be sent to database functions
      const selectAll string = "SELECT * FROM gocron;"

      // Perform a SELECT ALL
      rows, status := gocronlib.SelectDatabase(selectAll, verbose)
      if status == false {
            gocronlib.CronLog("Failed to perform SELECT ALL", verbose)
            return
      }

      // Iterate each row
      for rows.Next() {
            // Assign row results to a Cron struct
            var c gocronlib.Cron
            rows.Scan(&c.Cronname,
                        &c.Account,
                        &c.Email,
                        &c.Ipaddress,
                        &c.Frequency,
                        &c.Lastruntime,
                        &c.Alerted)

            var updateFail string = "Failed to update row for " + c.Cronname
            var currentTime = int(time.Now().Unix())
            var lastRunTime, _ = strconv.Atoi(c.Lastruntime)
            var frequency, _ = strconv.Atoi(c.Frequency)

            // If job not checked in on time
            if (currentTime - lastRunTime) > frequency {

                  // Mark row as alerted if not already true
                  if c.Alerted != true {
                        query = "UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + c.Cronname + "' AND account = '" + c.Account + "';"

                        // Perform the query
                        result = gocronlib.InsertDatabase(query, verbose)
                        if result == false {
                              gocronlib.CronLog(updateFail, verbose)
                              return
                        }

                        // Query was successful - Trigger alert
                        subject = c.Cronname + ": " + c.Account + " failed to check in" + "\n"
                        message = "The cronjob " + c.Cronname + " for account " + c.Account + " has not checked in on time"
                        alert(c, subject, message)
                        gocronlib.CronLog(subject, verbose)
                        return

                  // If alerted already marked true
                  } else {
                        gocronlib.CronLog("Alert for " + c.Cronname + ": " + c.Account +
                              " has been supressed. Already alerted", verbose)
                        return
                  }


            // If checked in on time but previously not (alerted == true)
            } else if ((currentTime - lastRunTime) < frequency) && c.Alerted == true {
                  query = "UPDATE gocron SET alerted = false " +
                          "WHERE cronname = '" + c.Cronname + "' AND account = '" + c.Account + "';"

                  result = gocronlib.InsertDatabase(query, verbose)
                  if result == false {
                        gocronlib.CronLog(updateFail, verbose)
                        return
                  }

                  // Query was successful - Trigger alert
                  subject = c.Cronname + ": " + c.Account + " is back online" + "\n"
                  message = "The cronjob " + c.Cronname + " for account " + c.Account + " is back online"
                  alert(c, subject, message)
                  gocronlib.CronLog(subject, verbose)
                  return

            // Job in a good state
            } else {
                  subject = c.Cronname + ": " + c.Account + " is online" + "\n"
                  gocronlib.CronLog(subject, verbose)
                  return
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
