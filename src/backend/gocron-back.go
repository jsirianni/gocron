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
      version    string = "2.2.2"
      libVersion string = gocronlib.Version
)

var (
      verbose    bool  // Command line flag
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
            cronStatus()
      }
}


func cronStatus() {
      checkMissedJobs("SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;")
      checkRevivedJobs("SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;")
}


func checkMissedJobs(query string) {
      //query := "SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;"
      rows, status := gocronlib.QueryDatabase(query, verbose)
      defer rows.Close()
      if status == false {
            gocronlib.CronLog("Failed to perform query: " + query, verbose)
            return
      }

      for rows.Next() {
            var cron gocronlib.Cron
            rows.Scan(&cron.Cronname,
                        &cron.Account,
                        &cron.Email,
                        &cron.Ipaddress,
                        &cron.Frequency,
                        &cron.Lastruntime,
                        &cron.Alerted,
                        &cron.Site)

            var updateFail string = "Failed to update row for " + cron.Cronname

            if cron.Alerted != true {
                  subject := cron.Cronname + ": " + cron.Account + " failed to check in" + "\n"
                  message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " has not checked in on time"

                  // Only update database if alert sent successful
                  if alert(cron, subject, message) == true {
                        query = "UPDATE gocron SET alerted = true " +
                                "WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

                        rows, result := gocronlib.QueryDatabase(query, verbose)
                        defer rows.Close()
                        if result == false {
                              gocronlib.CronLog(updateFail, verbose)
                        }
                  }

            } else {
                  gocronlib.CronLog("Alert for " + cron.Cronname + ": " + cron.Account +
                        " has been supressed. Already alerted", verbose)
            }
      }
}


func checkRevivedJobs(query string) {
      //query := "SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;"
      rows, status := gocronlib.QueryDatabase(query, verbose)
      defer rows.Close()
      if status == false {
            gocronlib.CronLog("Failed to perform query: " + query, verbose)
            return
      }

      for rows.Next() {
            var cron gocronlib.Cron
            rows.Scan(&cron.Cronname,
                        &cron.Account,
                        &cron.Email,
                        &cron.Ipaddress,
                        &cron.Frequency,
                        &cron.Lastruntime,
                        &cron.Alerted,
                        &cron.Site)

            var updateFail string = "Failed to update row for " + cron.Cronname
            query = "UPDATE gocron SET alerted = false " +
                    "WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

            rows, result := gocronlib.QueryDatabase(query, verbose)
            defer rows.Close()
            if result == false {
                  gocronlib.CronLog(updateFail, verbose)

            }

            subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
            message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
            alert(cron, subject, message)
      }
}


func alert(cron gocronlib.Cron, subject string, message string) bool {

      // Immediately log the alert
      gocronlib.CronLog(subject, verbose)

      var (
            recipient string = cron.Email
            port, _ = strconv.Atoi(config.Smtpport)
            d       = gomail.NewDialer(config.Smtpserver, port, config.Smtpaddress, config.Smtppassword)
            m       = gomail.NewMessage()
      )

      m.SetHeader("From", config.Smtpaddress)
      m.SetHeader("To", recipient)
      m.SetHeader("Subject", subject)
      m.SetBody("text/html", message)

      // Failed to send alert
      if err := d.DialAndSend(m); err != nil {
            gocronlib.CheckError(err, verbose)
            return false
      }

      // Alert sent
      gocronlib.CronLog("Alert for " + cron.Cronname + " sent to " + recipient, verbose)
      return true
}
