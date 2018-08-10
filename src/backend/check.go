package main


import (
	"fmt"

	"../gocronlib"
)


const (
	missedJobs = "SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;"
	revivedJobs = "SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;"
)


func cronStatus() {
	checkMissedJobs(missedJobs)
	checkRevivedJobs(revivedJobs)
}


func checkMissedJobs(query string) {
	rows, status := gocronlib.QueryDatabase(query, verbose)
	defer rows.Close()
	if status == false {
		gocronlib.CronLog("Failed to perform query: "+query, verbose)
		return
	}

	for rows.Next() {
		var cron gocronlib.Cron
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

				rows, result := gocronlib.QueryDatabase(query, verbose)
				defer rows.Close()
				if result == false {
					gocronlib.CronLog("Failed to update row for " + cron.Cronname, verbose)
				}
			}

		} else {
			gocronlib.CronLog("Alert for "+cron.Cronname+": "+cron.Account+
				" has been supressed. Already alerted", verbose)
		}
	}
}


func checkRevivedJobs(query string) {
	rows, status := gocronlib.QueryDatabase(query, verbose)
	defer rows.Close()
	if status == false {
		gocronlib.CronLog("Failed to perform query: "+query, verbose)
		return
	}

	for rows.Next() {
		var cron gocronlib.Cron
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

		rows, result := gocronlib.QueryDatabase(query, verbose)
		defer rows.Close()
		if result == false {
			gocronlib.CronLog("Failed to update row for " + cron.Cronname, verbose)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		alert(cron, subject, message)
	}
}


func getSummary() {
	var message string = "gocron summary - missed jobs:\n"

	rows, status := gocronlib.QueryDatabase(missedJobs, verbose)
	defer rows.Close()
	if status == false {
		gocronlib.CronLog("Failed to perform query while attempting to build a summary: " + missedJobs, verbose)
		return
	}

	for rows.Next() {
		var cron gocronlib.Cron
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

	// If verbose is true, send alert
	// Useful if running from cron and not the command line
	if verbose == true {

		// build cummy cron struct
		var c gocronlib.Cron

		// Send slack alert and pass dummy cron object
		if slackAlert(c, "gocron alert summary", message) == true {
			gocronlib.CronLog(message, verbose)
			return

		} else {
			gocronlib.CronLog("GOCRON: Failed to build alert summary.", verbose)
		}
	} else {
		fmt.Printf(message)
	}
}
