package main
import (
      "fmt"
      "os/exec"
)



//Return a Postgres connection string
func databaseString() string {
      var connectionString string = "postgres://" +
            config.Dbuser + ":" +
            config.Dbpass + "@" +
            config.Dbfqdn +
            "/gocron" +
            "?sslmode=disable"

      return connectionString
}



// Function validates SQL variables
func validateArgs(c Cron) bool {

      // Flag determines the return value
      var valid bool = false

      // Perform validation of parameters
      if checkLength(c) == true {
            if sqlInjection(c) == true {
                  valid = true
            }
      }

      // Log result if verbose is enabled
      if verbose == true {
            if valid == true {
                  cronLog("Parameters from " + c.ipaddress + " passed validation")
                  return true

            } else {
                  cronLog("Parameters from " + c.ipaddress + " failed validation!")
                  return false
            }
      }

      // Return true or false
      return valid
}



// Validate that parameters are present
func checkLength(c Cron) bool {
      if len(c.account) == 0 {
            return false

      } else if len(c.cronname) == 0 {
            return false

      } else if len(c.email) == 0 {
            return false

      } else if len(c.frequency) == 0 {
            return false

      } else if len(c.ipaddress) == 0 {
            return false

      } else if len(c.lastruntime) == 0 {
            return false

      } else {
            return true
      }
}



// Prevent SQL injection
func sqlInjection(c Cron) bool {
      // TODO
      return true
}



// Function writes messages to syslog and (optionally) to standard out
func cronLog(message string) {
      err := exec.Command("logger", message).Run()
      if err != nil {
            fmt.Println("Failed to write to syslog")
            fmt.Println(message)
      }
      if verbose == true {
            fmt.Println(message)
      }
}



// Function passes error messages to the cronLog() function
func checkError(err error) {
      if err != nil {
            cronLog("Error: \n" + err.Error())
      }
}
