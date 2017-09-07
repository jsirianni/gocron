package main

import (
      "fmt"
      "io"
      "net/http"
)


// TODO
// Build yaml config file that contains the following
// - database connection string (host, port, user, pass)


// Start HTTP Server
func main() {
      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}


// Parse GET Request parameters
func cronStatus(w http.ResponseWriter, r *http.Request) {
      fmt.Println("GET parameters:", r.URL.Query())
      io.WriteString(w, "Thanks!")
}


// Insert or update a cron entry in the database
func updateDatabase(cron string, email string, time string, tolerance int32) {
      // TODO
      // Read in config file
      //
      // Check the database for the existing cron (primary key)
      // If the cron exists, update its lastRun column with the current time
      //
      // If the cron does not exist, create a table entry and record lastRun with the current time
      // Send an email alert notifying the user that the entry has been made

}

// Check for missed cron updateDatabase
func checkCronStatus() {
      // TODO
      // Check the database for entries that have
      // not ran at their scheduled time + their tolerance
      //
      // Example: 1_*_*_*_* 30 should run at least every 1.5 hours
      //          Every hour with 30 minutes of tolerance
      //
      // Example: 0_19_*_*_* 120 should run at 7pm every day
      //          with 2 hours of tolerance (7-9pm)
      //
      // TODO Run this function every 10 minutes ??
      //
      // Send email alerts for any entries that have not checked in on time
}

// Send emails
func alert(recipient string, subject string, message string) {
      // TODO
      // Send an email alert

      // TODO
      // Add optional slack alerts
}
