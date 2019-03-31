package libgocron
import (
    "os"
    "fmt"
)


// CronLog writes messages to syslog and (optionally) to standard out
func CronLog(message string) {
    fmt.Println(message)
}

func LogError(err error) {
    fmt.Fprintln(os.Stderr, err.Error())
}
