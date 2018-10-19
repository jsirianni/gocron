# gocron
Service that monitors the status of your cron jobs

The goal of this service is to receive an email alert when a cronjob does
not run after a predetermined amount of time.

Email alerts are sent one time and then suppressed. Alerts are re-triggered only if the job checks in again, and then misses its next run.

## Architecture
GoCron is made up of several services
- gocron-front
- gocron-back
- gocronlib
- gocron web interface https://github.com/jsirianni/gocron-frontend

Gocron web interface is an optional component that allows the user
to view the status of all jobs.

## Usage
These examples will notify the server to expect a notification every hour. If the job
does not check in within one hour, an alert is sent. Future notifications are
suppressed until the job checks in again.

#### GET
Send a GET request with the following parameters in the query string
- cronname
- account
- email
- frequency (seconds)

Test with
```
curl -v "172.17.0.2:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600"
```
Append to an existing crontab entry with:
```
&& curl -v "172.17.0.2:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600"
```

#### POST
```
curl -v -X POST -d "{\"cronname\":\"test\",\"account\":\"test account\", \"email\":\"test@gmail.com\",\"frequency\":20}" http://172.17.0.2:8080
```



## Weekly Summary
The backend service binary can provide a summary of all missed jobs
```
# print a summary of current missed jobs
./gocron-back --summary

# print a summary of current missed jobs and send it via slack
./gocron-back --summary --verbose
```

Run from cron every monday at 9am
```
0 9 * * MON /usr/local/bin/gocron-back --summary --verbose >> /dev/null
```

## Sizing
A single Google Compute Engine f1-micro instance has been proven to handle 50,000 jobs
that check in at a rate of 10,000 jobs per 90 seconds. At this rate, the load was less than
15%. The CPU would max out every 5 minutes when the service would check the entire database
for jobs that have not checked in.

If this kind of load is expected, it is possible to run multiple `gocron-front` services
behind a load balancer and a single `gocron-back` service on a single machine. Additionally,
the database can live on a separate system entirely.


## Installing

### Postgresql must be installed and listening on localhost
```
CREATE DATABASE gocron;
CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency integer, lastruntime integer, alerted boolean, site boolean, PRIMARY KEY(cronname, account));
CREATE USER gocron WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON gocron TO gocron;
```

### Run install.sh
```
chmod +x install.sh
sudo ./install.sh
```

This script will:
- Copy the gocron binaries to `/usr/local/bin`
- Copy the config example to `/etc/gocron`
- Create a Systemd services

### Manage services
```
systemctl status gocron-front
systemctl status gocron-back
```


## Building
### Compile for linux
simply run `build_linux.sh`, which places the compuled binaries into `bin/`

### Docker
Dockerfiles are included in `src/fronend` and `src/backend`.
```
docker run -d \
  -p 8080:8080
  -v /path/to/config/dir:/etc/gocron gocron-front
```
https://hub.docker.com/r/firefoxx04/gocron-front/
https://hub.docker.com/r/firefoxx04/gocron-back/


## Notes
The main purpose of this project is to gain familiarity with golang. If you have improvements or suggestions, please feel free to file an issue or open a pull request.
