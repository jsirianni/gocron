// Version 1.0.0
// Debian 9 Officially supported

package main
import (
      "time"
      "strings"
      "strconv"
      "net/http"
      "io/ioutil"
      "gopkg.in/yaml.v2"
      "database/sql"; _ "github.com/lib/pq";
)


type Config struct {
      Dbfqdn string
      Dbport string
      Dbuser string
      Dbpass string
      Dbdatabase string
      Smtpserver string
      Smtpport string
      Smtpaddress string
      Smtppassword string
}

type Cron struct {
      cronname string
      account string
      email string
      ipaddress string
      frequency string
      tolerance string
      lastruntime string  // Unix timestamp
      alerted bool        // set to true if an alert has already been thrown
}

var config Config


func main() {
      yamlFile, err := ioutil.ReadFile("/etc/gocron/.config.yml")
      if err != nil {
            checkError(err)
      }
      err = yaml.Unmarshal(yamlFile, &config)
      if err != nil {
            checkError(err)
      }

      go timer()

      http.HandleFunc("/", cronStatus)
      http.ListenAndServe(":8080", nil)
}


func cronStatus(w http.ResponseWriter, r *http.Request) {
      var currentTime int = int(time.Now().Unix())
      var socket = strings.Split(r.RemoteAddr, ":")
      var cronJob Cron

      cronJob.cronname = r.URL.Query().Get("cronname")
      cronJob.account = r.URL.Query().Get("account")
      cronJob.email = r.URL.Query().Get("email")
      cronJob.frequency = r.URL.Query().Get("frequency")
      cronJob.tolerance = r.URL.Query().Get("tolerance")
      cronJob.lastruntime = strconv.Itoa(currentTime)
      cronJob.ipaddress = socket[0]

      if checkLength(cronJob) == true {
            go updateDatabase(cronJob)

      } else {
            cronLog("GET request not valid. Dropping.")
      }
}


func updateDatabase(c Cron) {
      var query string
      query = "INSERT INTO gocron (cronname, account, email, ipaddress, frequency, tolerance, lastruntime) VALUES ('" +
             c.cronname + "','" +
             c.account + "','" +
             c.email + "','" +
             c.ipaddress + "','" +
             c.frequency + "','" +
             c.tolerance + "','" +
             c.lastruntime + "') " +
             "ON CONFLICT (cronname, account) DO UPDATE " +
             "SET email = " + "'" + c.email + "'," +
             "ipaddress = " + "'" + c.ipaddress + "'," +
             "frequency = " + "'" + c.frequency + "'," +
             "tolerance = " + "'" + c.tolerance + "'," +
             "lastruntime = " + "'" + c.lastruntime + "'" +
             ";"

      db, err := sql.Open("postgres", databaseString())
      if err != nil {
            checkError(err)
      }
      defer db.Close()

      _, err = db.Exec(query)
      if err != nil {
            checkError(err)

      } else {
            cronLog("Heartbeat from " + c.cronname + ": " + c.account + " \n")
      }
}
