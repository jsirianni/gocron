package main
import (
      "fmt"
      "time"
      "flag"
      "strconv"
      "gopkg.in/gomail.v2"
      "github.com/jsirianni/gocronlib"
)



const (
      version string    = "2.0.8"
      libVersion string = gocronlib.Version
)

var (
      verbose bool     // Command line flag
      getVersion bool  // Command line flag
      config gocronlib.Config = gocronlib.GetConfig(verbose)
)


func main() {
      flag.BoolVar(&getVersion, "version", false, "Get the version and then exit")
      flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
      flag.Parse()

      if getVersion == true {
            fmt.Println("gocron-back version:", version)
            fmt.Println("gocronlib version:", libVersion)
            return
      }

      if verbose == true {
            fmt.Println("Verbose mode enabled")
            fmt.Println("gocron-back version:", version)
            fmt.Println("gocronlib version:", libVersion)
            fmt.Println("Using check interval:", config.Interval)
      }

      timer()
}


// Function calls checkCronStatus() on a set interval
func timer() {
      for {
            time.Sleep((time.Duration(config.Interval) * time.Second))
            gocronlib.CronLog("Checking for missed jobs.", verbose)
            checkCronStatus()
      }
}


func checkCronStatus() {
      var (
            subject string  // Subject used in alerts
            message string  // Message used in alerts
            result bool     // Handles Insert Database responses
            query string    // Queries to be sent to database functions
      )

      rows, status := gocronlib.QueryDatabase("SELECT * FROM gocron;", verbose)
      defer rows.Close()
      if status == false {
            gocronlib.CronLog("Failed to perform SELECT ALL", verbose)
            return
      }

      for rows.Next() {
            // Assign row results to a Cron struct
            var cron gocronlib.Cron
            rows.Scan(&cron.Cronname,
                        &cron.Account,
                        &cron.Email,
                        &cron.Ipaddress,
                        &cron.Frequency,
                        &cron.Lastruntime,
                        &cron.Alerted,
                        &cron.Site)

            var (
                  updateFail string = "Failed to update row for " + cron.Cronname
                  currentTime = int(time.Now().Unix())
                  lastRunTime, _ = strconv.Atoi(cron.Lastruntime)
                  frequency, _ = strconv.Atoi(cron.Frequency)
            )


            // If job not checked in on time
            if (currentTime - lastRunTime) > frequency {

                  // Mark row as alerted if not already true
                  if cron.Alerted != true {
                        query = "UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

                        // Perform the query
                        rows, result = gocronlib.QueryDatabase(query, verbose)
                        defer rows.Close()
                        if result == false {
                              gocronlib.CronLog(updateFail, verbose)

                        }

                        // Query was successful - Trigger alert
                        subject = cron.Cronname + ": " + cron.Account + " failed to check in" + "\n"
                        message = "The cronjob " + cron.Cronname + " for account " + cron.Account + " has not checked in on time"
                        alert(cron, subject, message)
                        gocronlib.CronLog(subject, verbose)


                  // If 'alerted' already  true
                  } else {
                        gocronlib.CronLog("Alert for " + cron.Cronname + ": " + cron.Account +
                              " has been supressed. Already alerted", verbose)
                  }


            // If checked in on time but previously not (alerted == true)
            } else if ((currentTime - lastRunTime) < frequency) && cron.Alerted == true {
                  query = "UPDATE gocron SET alerted = false " +
                          "WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

                  rows, result = gocronlib.QueryDatabase(query, verbose)
                  defer rows.Close()
                  if result == false {
                        gocronlib.CronLog(updateFail, verbose)

                  }

                  // Query was successful - Trigger alert
                  subject = cron.Cronname + ": " + cron.Account + " is back online" + "\n"
                  message = "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
                  alert(cron, subject, message)
                  gocronlib.CronLog(subject, verbose)


            } else {
                  subject = cron.Cronname + ": " + cron.Account + " is online" + "\n"
                  gocronlib.CronLog(subject, verbose)

            }
      }
}


func alert(cron gocronlib.Cron, subject string, message string) {

      var recipient string = cron.Email
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

      gocronlib.CronLog("Alert for " + cron.Cronname + " sent to " + recipient, verbose)
}
