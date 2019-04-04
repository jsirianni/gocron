package libgocron
import (
	"os"
	"time"
	"errors"

	"gocron/util/log"
	"gocron/util/slack"

)


// StartBackend calls checkCronStatus() on a set interval
func (g Gocron) StartBackend() error {
	log.Message("gocron-back version: " + Version)

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

	err := g.slackAlert(subject, message)
	if err != nil {
		log.Error(errors.New("gocron fail: alert for " + cron.Cronname))
		log.Error(err)
		return false
	}

    log.Message("gocron success: alert for " + cron.Cronname + " sent")
    return true
}


func (g Gocron) slackAlert(subject string, message string) error {
	var slack slack.Slack
	slack.HookURL      = g.SlackHookURL
	slack.Post.Channel = g.SlackChannel
	slack.Post.Text    = message
	return slack.Message()
}
