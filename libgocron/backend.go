package libgocron
import (
	"os"
	"time"
	"errors"

	"gocron/util/log"

	"github.com/jsirianni/slacklib/slacklib"
)


// StartBackend calls checkCronStatus() on a set interval
func (g Gocron) StartBackend() error {
	// create the gocron table, if not exists
	err := g.createGocronTable()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// backend server is just a never ending loop that checks for missed
	// jobs at the set interval
	for {
		time.Sleep((time.Duration(g.Interval) * time.Second))
		log.Message("Checking for missed jobs.")
		g.cronStatus()
	}
}

// GetSummary prints a summary to standard out
func (g Gocron) GetSummary() {
	message := "gocron summary - missed jobs:\n"

	rows, err := queryDatabase(g, missedJobs)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query while attempting to build a summary: " + missedJobs))
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


	// Send slack alert
	if g.slackAlert("gocron alert summary", message) == true {
		log.Message(message)

	} else {
		log.Message("GOCRON: Failed to build alert summary.")
	}
}


func (g Gocron) cronStatus() {
	g.checkMissedJobs(missedJobs)
	g.checkRevivedJobs(revivedJobs)
}


func (g Gocron) checkMissedJobs(query string) {
	rows, err := queryDatabase(g, query)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query: " + query))
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
			if g.alert(cron, subject, message) == true {
				query = "UPDATE gocron SET alerted = true " +
					"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

				rows, err := queryDatabase(g, query)
				defer rows.Close()
				if err != nil {
					log.Error(err)
					log.Error(errors.New("Failed to update row for " + cron.Cronname))
				}
			}

		} else {
			log.Message("Alert for "+cron.Cronname+": "+cron.Account+
				" has been supressed. Already alerted")
		}
	}
}


func (g Gocron) checkRevivedJobs(query string) {
	rows, err := queryDatabase(g, query)
	defer rows.Close()
	if err != nil {
		log.Error(err)
		log.Error(errors.New("Failed to perform query: " + query))
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

		rows, err := queryDatabase(g, query)
		defer rows.Close()
		if err != nil {
			log.Message("Failed to update row for " + cron.Cronname)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		g.alert(cron, subject, message)
	}
}


func (g Gocron) alert(cron Cron, subject string, message string) bool {

    // Immediately log the alert
    log.Message(subject)

    result := false
	if g.slackAlert(subject, message) == true {
		result = true
	}

    if result == true {
        log.Message("gocron success: alert for " + cron.Cronname + " sent")
        return true
    }

    log.Message("gocron fail: alert for " + cron.Cronname)
    return false
}


func (g Gocron) slackAlert(subject string, message string) bool {
    var slackmessage slacklib.SlackPost
    slackmessage.Channel = g.SlackChannel
    slackmessage.Text = message
    return slacklib.BasicMessage(slackmessage, g.SlackHookURL)
}
