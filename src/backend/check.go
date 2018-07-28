package main


import (
    "github.com/jsirianni/gocronlib"
)


func cronStatus() {
	checkMissedJobs("SELECT * FROM gocron WHERE (extract(epoch from now()) - lastruntime) > frequency;")
	checkRevivedJobs("SELECT * FROM gocron WHERE alerted = true AND (extract(epoch from now()) - lastruntime) < frequency;")
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

		var updateFail string = "Failed to update row for " + cron.Cronname

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
					gocronlib.CronLog(updateFail, verbose)
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

		var updateFail string = "Failed to update row for " + cron.Cronname
		query = "UPDATE gocron SET alerted = false " +
			"WHERE cronname = '" + cron.Cronname + "' AND account = '" + cron.Account + "';"

		rows, result := gocronlib.QueryDatabase(query, verbose)
		defer rows.Close()
		if result == false {
			gocronlib.CronLog(updateFail, verbose)

		}

		subject := cron.Cronname + ": " + cron.Account + " is back online" + "\n"
		message := "The cronjob " + cron.Cronname + " for account " + cron.Account + " is back online"
		alert(cron, subject, message)
	}
}
