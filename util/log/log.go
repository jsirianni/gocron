package log
import (
    "os"
    "fmt"
    "net/http"
)

// APIReq defines an incoming api call
type APIReq struct {
    Host   string
    Path   string
    Method string
}

// Message writes messages to syslog and (optionally) to standard out
func Message(message string) {
    fmt.Println(message)
}

// Error takes an error and prints it to standard error
func Error(err error) {
    fmt.Fprintln(os.Stderr, err.Error())
}

// APILog takes an http request and prints it to standard out
func APILog(req *http.Request) {
	var x APIReq
    x.Host = req.Host
    x.Path = req.URL.RequestURI()
    x.Method = req.Method
    fmt.Println(x)
}
