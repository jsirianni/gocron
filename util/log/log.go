package log
import (
    "os"
    "fmt"
)


// Message writes messages to syslog and (optionally) to standard out
func Message(message string) {
    fmt.Println(message)
}

// Error takes an error and prints it to standard error
func Error(err error) {
    fmt.Fprintln(os.Stderr, err.Error())
}
