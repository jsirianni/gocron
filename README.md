# gocron
Service that monitors the status of your cron jobs. The goal of this service is to
receive an alert when a cronjob does not run after a predetermined amount of time.


## Architecture
Gocron is made up of several services
- gocron frontend
- gocron backend
- Postgresql

Gocron web interface is an optional component that allows the user
to view the status of all jobs. The frontend service can be scaled to any number
of nodes, if required.


## Usage
These examples will notify the server to expect a notification every hour. If the job
does not check in within one hour, an alert is sent. Future notifications are
suppressed until the job checks in again.

HTTP POST and GET are supported, however, POST is recommend. Send a request
with the following parameters:
- cronname
- account
- email
- frequency (seconds)

### POST
```
curl -v -X POST -d "{\"cronname\":\"test\",\"account\":\"test account\", \"email\":\"test@gmail.com\",\"frequency\":20}" http://172.17.0.2:8080
```

### GET
```
curl -v "172.17.0.2:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600"
```
Append to an existing crontab entry with:
```
&& curl -v "172.17.0.2:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600"
```


### Weekly Summary
The backend service binary can provide a summary of all missed jobs
```
# print a summary of current missed jobs
./gocron backend --summary

# print a summary of current missed jobs and send it via slack
./gocron backend --summary --verbose
```

Run from cron every monday at 9am
```
0 9 * * MON /usr/local/bin/gocron-back --summary --verbose >> /dev/null
```


## Installation

Docker is the default deployment method as of version `4.0.3`. Systemd is also
an option. Examples can be found in `install.sh` but are not guaranteed to be up
to data. The docker image is based on `alpine:latest`.

### Database
Postgresql must be configured:
```
CREATE DATABASE gocron;
CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency integer, lastruntime integer, alerted boolean, site boolean, PRIMARY KEY(cronname, account));
CREATE USER gocron WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON gocron TO gocron;
```

### Docker
Setup docker swarm
```
sudo docker swarm init
```
Setup environment:
```
sudo vim docker/docker.env
```
Deploy:
```
sudo docker stack deploy gocron --compose-file docker/docker-compose.yml
```


## Building

### Compile
`build_linux.sh` will compile `gocron` and build the docker image.


## Notes
The main purpose of this project is to gain familiarity with golang. If you have improvements or suggestions, please feel free to file an issue or open a pull request.
