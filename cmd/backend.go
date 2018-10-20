package cmd
import (
	"os"
	"time"
	"strconv"
	"gopkg.in/gomail.v2"

	"github.com/spf13/cobra"
	"github.com/jsirianni/slacklib/slacklib"
)


// backendCmd represents the backend command
var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Start the gocron backend server",
	Long: "Start the gocron backend server, which alerts on missed jobs",
	Run: func(cmd *cobra.Command, args []string) {
		startBackend()
	},
}


func init() {
	rootCmd.AddCommand(backendCmd)
	backendCmd.Flags().BoolVar(&summary, "summary", false, "Get summary")
}


func startBackend() {
	// Initilize the config struct
	config = GetConfig(verbose)

	if summary == true {
		// If verbose == true, summary will send to syslog AND the configured
		// alert system
		getSummary()
		return
	}

	CronLog("verbose mode enabled", verbose)
	CronLog("gocron-back version: " + VERSION, verbose)
	CronLog("check interval: " + strconv.Itoa(config.Interval), verbose)

	if config.PreferSlack == true {
		CronLog("Prefer slack: enabled", verbose)
		CronLog("Slack channel: " + config.SlackChannel, verbose)
		CronLog("Slack hook url: " + config.SlackHookUrl, verbose)
	}


	// create the gocron table, if not exists
	if CreateGocronTable(verbose) == false {
		os.Exit(1)
	}

	timer()
}


// Function calls checkCronStatus() on a set interval
func timer() {
	for {
		time.Sleep((time.Duration(config.Interval) * time.Second))
		CronLog("Checking for missed jobs.", verbose)
		cronStatus()
	}
}


func cronStatus() {
	checkMissedJobs(missedJobs)
	checkRevivedJobs(revivedJobs)
}


func checkMissedJobs(query string) {
	rows, status := QueryDatabase(query, verbose)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query: "+query, verbose)
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

				rows, result := QueryDatabase(query, verbose)
				defer rows.Close()
				if result == false {
					CronLog("Failed to update row for " + cron.Cronname, verbose)
				}
			}

		} else {
			CronLog("Alert for "+cron.Cronname+": "+cron.Account+
				" has been supressed. Already alerted", verbose)
		}
	}
}


func checkRevivedJobs(query string) {
	rows, status := QueryDatabase(query, verbose)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query: "+query, verbose)
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

		rows, result := QueryDatabase(query, verbose)
		defer rows.Close()
		if result == false {
			CronLog("Failed to update row for " + cron.Cronname, verbose)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		alert(cron, subject, message)
	}
}


func getSummary() {
	var message string = "gocron summary - missed jobs:\n"

	rows, status := QueryDatabase(missedJobs, verbose)
	defer rows.Close()
	if status == false {
		CronLog("Failed to perform query while attempting to build a summary: " + missedJobs, verbose)
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
	var c Cron // build dummy cron stru
	if verbose == true && slackAlert(c, "gocron alert summary", message) == true {
		CronLog(message, verbose)
		return

	} else {
		CronLog("GOCRON: Failed to build alert summary.", verbose)
	}

}


// alert will send slack messages first, if enabled, and fallback
// to email alerts if the slack notification fails
func alert(cron Cron, subject string, message string) bool {

    // Immediately log the alert
    CronLog(subject, verbose)

    var result bool
	if config.PreferSlack == true {

        if slackAlert(cron, subject, message) == true {
			result = true
		} else {
			result = emailAlert(cron, subject, message)
		}

	} else {
		result = emailAlert(cron, subject, message)
	}


    if result == true {
        CronLog("gocron success: alert for " + cron.Cronname + " sent", verbose)
        return true
    } else {
        CronLog("gocron fail: alert for " + cron.Cronname, verbose)
        return false
    }
}


func emailAlert(cron Cron, subject string, message string) bool {

	var (
		recipient string = cron.Email
		port, _          = strconv.Atoi(config.Smtpport)
		d                = gomail.NewDialer(config.Smtpserver, port, config.Smtpaddress, config.Smtppassword)
		m                = gomail.NewMessage()
	)

	m.SetHeader("From", config.Smtpaddress)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	// Failed to send alert
	if err := d.DialAndSend(m); err != nil {
		CheckError(err, verbose)
		return false
	}
	return true
}


func slackAlert(cron Cron, subject string, message string) bool {
    var slackmessage slacklib.SlackPost
    slackmessage.Channel = config.SlackChannel
    slackmessage.Text = message
    return slacklib.BasicMessage(slackmessage, config.SlackHookUrl)
}
