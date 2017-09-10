# gocron
GO Service that monitors the status of your cron jobs


CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency varchar, tolerance int, lastruntime varchar, PRIMARY KEY(cronname, account));

curl -v "localhost:8080/?cronname=zfs-2&account=tes21&email=alerts.teamit.localnet&frequency=60&tolerance=30"
