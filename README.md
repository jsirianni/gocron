# gocron
GO Service that monitors the status of your cron jobs




Send a GET request with the following parameters in the query string
- cronname
- account
- email
- frequency (seconds)
- tolerance (seconds)

curl -v "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=60&tolerance=30"


The frequency is the amount of time the job should run (and check in). Tolerance is the amount of time added to the frequency. For example, a backup job should run every 6000 seconds, but may be given a tolerance of 1000 seconds.

Email alerts are sent one time and then suppressed. Alerts are re-triggered only if the job checks in again, and then misses its next run.

Example: Job does not check in, resulting in an alert. Future misses are ignored. Once the job checks in again, it is once again eligible for an alert if it fails to check in in the future.  
