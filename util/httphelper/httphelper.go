package httphelper
import (
    "net/http"
)

// Return 200 OK
func ReturnOk(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
}


// Return a 201 Created
func ReturnCreated(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusCreated)
}


// Return a 500 Server Error
func ReturnServerError(resp http.ResponseWriter, message string, json bool) {
	resp.WriteHeader(http.StatusInternalServerError)
    if len(message) != 0 {
        if json == true {
            resp.Write([]byte("{\"error\":\"" + message + "\"}"))
        } else {
            resp.Write([]byte(message))      
        }
    }
}


// Return 404 Not Found
func ReturnNotFound(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
}
