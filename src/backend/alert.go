package main


import (
    "strconv"
	"gopkg.in/gomail.v2"

	"../gocronlib"
    "github.com/jsirianni/slacklib/slacklib"
)


// alert will send slack messages first, if enabled, and fallback
// to email alerts if the slack notification fails
func alert(cron gocronlib.Cron, subject string, message string) bool {

    // Immediately log the alert
    gocronlib.CronLog(subject, verbose)

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
        gocronlib.CronLog("gocron success: alert for " + cron.Cronname + " sent", verbose)
        return true
    } else {
        gocronlib.CronLog("gocron fail: alert for " + cron.Cronname, verbose)
        return false
    }
}


func emailAlert(cron gocronlib.Cron, subject string, message string) bool {

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
		gocronlib.CheckError(err, verbose)
		return false
	}
	return true
}


func slackAlert(cron gocronlib.Cron, subject string, message string) bool {

    var slackmessage slacklib.SlackPost

    slackmessage.Channel = config.SlackChannel
    slackmessage.Text = message

    return slacklib.BasicMessage(slackmessage, config.SlackHookUrl)
}
