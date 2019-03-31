package libgocron
import (
	"os"
	"fmt"
	"time"

	"github.com/jsirianni/slacklib/slacklib"
)


// Timer calls checkCronStatus() on a set interval
func (c Config) StartBackend(v bool) {

	// create the gocron table, if not exists
	if createGocronTable() == false {
		os.Exit(1)
	}

	// backend server is just a never ending loop that checks for missed
	// jobs at the set interval
	for {
		time.Sleep((time.Duration(c.Interval) * time.Second))
		CronLog("Checking for missed jobs.")
		cronStatus()
	}
}


func (c Config) GetSummary(v bool) {
	var message string = "gocron summary - missed jobs:\n"

	rows, status := queryDatabase(missedJobs)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query while attempting to build a summary: " + missedJobs)
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)

		message = message + "Name: " + cron.Cronname  + "| Account: " + cron.Account + "\n"
	}


	// Send slack alert and pass dummy cron object
	if v == true && slackAlert("gocron alert summary", message) == true {
		CronLog(message)
		return

	} else if v == false {
		fmt.Println(message)

	} else {
		CronLog("GOCRON: Failed to build alert summary.")
	}
}


func cronStatus() {
	checkMissedJobs(missedJobs)
	checkRevivedJobs(revivedJobs)
}


func checkMissedJobs(query string) {
	rows, status := queryDatabase(query)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query: "+query)
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)


		if cron.Alerted != true {
			subject := cron.Cronname + ": " + cron.Account + " failed to check in" + "\n"
			message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " has not checked in on time"

			// Only update database if alert sent successful
			if alert(cron, subject, message) == true {
				query = "UPDATE gocron SET alerted = true " +
					"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

				rows, result := queryDatabase(query)
				defer rows.Close()
				if result == false {
					CronLog("Failed to update row for " + cron.Cronname)
				}
			}

		} else {
			CronLog("Alert for "+cron.Cronname+": "+cron.Account+
				" has been supressed. Already alerted")
		}
	}
}


func checkRevivedJobs(query string) {
	rows, status := queryDatabase(query)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query: "+query)
		return
	}

	for rows.Next() {
		var cron Cron
		rows.Scan(&cron.Cronname,
			&cron.Account,
			&cron.Email,
			&cron.Ipaddress,
			&cron.Frequency,
			&cron.Lastruntime,
			&cron.Alerted,
			&cron.Site)

		query = "UPDATE gocron SET alerted = false " +
			"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

		rows, result := queryDatabase(query)
		defer rows.Close()
		if result == false {
			CronLog("Failed to update row for " + cron.Cronname)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		alert(cron, subject, message)
	}
}


func alert(cron Cron, subject string, message string) bool {

    // Immediately log the alert
    CronLog(subject)

    var result bool = false
	if slackAlert(subject, message) == true {
		result = true
	}

	// NOTE: future alert methods will go here. Removed SMTP due to complexity

    if result == true {
        CronLog("gocron success: alert for " + cron.Cronname + " sent")
        return true
    } else {
        CronLog("gocron fail: alert for " + cron.Cronname)
        return false
    }
}


func slackAlert(subject string, message string) bool {
    var slackmessage slacklib.SlackPost
    slackmessage.Channel = config.SlackChannel
    slackmessage.Text = message
    return slacklib.BasicMessage(slackmessage, config.SlackHookURL)
}
