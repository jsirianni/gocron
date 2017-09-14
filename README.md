# gocron
Service that monitors the status of your cron jobs

The goal of this service is to receive an email alert when a cronjob does
not run after a predetermined amount of time.

Email alerts are sent one time and then suppressed. Alerts are re-triggered only if the job checks in again, and then misses its next run.

### Usage
Send a GET request with the following parameters in the query string
- cronname
- account
- email
- frequency (seconds)
- tolerance (seconds)

Test with
`curl -v "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600&tolerance=120"`

Append to an existing crontab entry with:
`&& curl -v "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600&tolerance=120"`


The above examples will notify the server every hour with a tolerance of 2 minutes. If the job does not check in within 2 hours and 2 minutes, an email alert is sent. Future notifications are suppressed until the job checks in again.

### Notes
The main purpose of this project is to gain familiarity with golang. If you have improvements or suggestions, please feel free to file an issue or open a pull request.
