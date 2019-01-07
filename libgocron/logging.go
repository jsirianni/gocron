package libgocron
import (
    "fmt"
    "os/exec"
)


// Function writes messages to syslog and (optionally) to standard out
func CronLog(message string, verbose bool) {
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
func CheckError(err error, verbose bool) {
      if err != nil {
            CronLog("Error: \n" + err.Error(), verbose)
      }
}
