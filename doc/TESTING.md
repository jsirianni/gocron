# Testing

## CI

Google Cloud Build is configured to run on every commit

## Vagrant
A Vagrant configuration is provided. It will install and configure
gocron + postgres.

Deploy with Vagrant. Note, you must set `SLACK_HOOK_URL` and `SLACK_CHANNEL` variables.
```
env SLACK_HOOK_URL='https://hooks.slack.com/services/myhook' SLACK_CHANNEL=mychannel vagrant up
```
```
vagrant ssh
sudo docker service ls
sudo docker ps

curl -v "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600"
```
