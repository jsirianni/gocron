# gocron
GO Service that monitors the status of your cron jobs


CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipAddress varchar, cronTime varchar, tolerance int, lastRunTime varchar, PRIMARY KEY(cronName, account));
