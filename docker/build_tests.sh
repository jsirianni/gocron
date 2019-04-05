#!/bin/sh

# install docker compose
curl -s -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# configure postgres
echo "listen_addresses='*'" | tee -a /etc/postgresql/10/main/postgresql.conf
echo "host all all 0.0.0.0/0 trust" | tee -a /etc/postgresql/10/main/pg_hba.conf
service postgresql restart
sudo -u postgres createuser gocron
sudo -u postgres createdb gocron
sudo -u postgres -H -- psql -c "alter user gocron with encrypted password 'password'"
sudo -u postgres -H -- psql -c "grant all privileges on database gocron to gocron"

# setup the environment
export GC_DBFQDN=`hostname -i`
export GC_DBPORT=5432
export GC_DBUSER=gocron
export GC_DBPASS=password
export GC_DBDATABASE=gocron
export GC_INTERVAL=3
export GC_SLACKHOOKURL="https://httpstat.us/200"
export GC_SLACKCHANNEL="test"

# start backend and then frontend services
echo "starting backend service. . ."
/usr/local/bin/gocron backend &> backend_log &
sleep 5
echo "starting frontend service. . ."
/usr/local/bin/gocron frontend &> frontend_log &
sleep 2

# test for 201 status code on healthcheck endpoint
echo "testing for 200 response from frontend healthcheck endpoint"
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:8080/healthcheck" -o /dev/null`
if [ "$STATUS_CODE" = "200" ];
then
   echo "PASS: frontend healthcheck returned 200" ;
else
   echo "FAIL: frontend returned ${STATUS_CODE}, expected 200" ;
   exit 1
fi

# test for 201 status code when sending a valid GET
# test for "not checked in" alert
echo "testing for 201 response from frontend"
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=1" -o /dev/null`
if [ "$STATUS_CODE" = "201" ];
then
   echo "PASS: frontend returned 201" ;
else
   echo "FAIL: frontend returned ${STATUS_CODE}, expected 201" ;
   exit 1
fi

# sleep 8 seconds to allow alert to be sent
sleep 8

# check for missed jobs
echo "checking for missed jobs via api"
curl -s localhost:3000/crons/missed | jq . || exit 1

# test for 201 status code when sending a valid GET
# test for "back online" alert
echo "testing for 201 response from frontend"
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:8080/?cronname=mycronjob&account=myaccount&email=myemail@gmail.com&frequency=3600" -o /dev/null`
if [ "$STATUS_CODE" = "201" ];
then
   echo "PASS: frontend returned 201" ;
else
   echo "FAIL: frontend returned ${STATUS_CODE}, expected 201" ;
   exit 1
fi

# sleep 12 seconds to allow alert to be sent
# and to allow for allert to be suppressed in the log
sleep 12

# test for 404 status when sending invalid query string
echo "testing for 404 response from frontend (bad GET)"
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:8080/" -o /dev/null`
if [ "$STATUS_CODE" = "404" ];
then
   echo "PASS: frontend returned 404" ;
else
   echo "FAIL: frontend returned ${STATUS_CODE}, expected 404" ;
   exit 1
fi

# test backend healthcheck
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:3000/healthcheck" -o /dev/null`
if [ "$STATUS_CODE" = "200" ];
then
   echo "PASS: backend returned 200" ;
else
   echo "FAIL: backend returned ${STATUS_CODE}, expected 200" ;
   exit 1
fi

# validate backend version api works and is valid json
echo "checking backend /version endpoint"
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:3000/version" -o /dev/null`
if [ "$STATUS_CODE" = "200" ];
then
    curl -s localhost:3000/version | jq '.version' || exit 1
    curl -s localhost:3000/version | jq '.database' || exit 1
    echo "PASS: backend /version"
else
   echo "FAIL: backend /version returned ${STATUS_CODE}, expected 200" ;
   exit 1
fi

# validate backend crons api
echo "checking backend /crons endpoint"
curl -s localhost:3000/crons | jq . || exit 1
STATUS_CODE=`curl -sL -w "%{http_code}\\n" "localhost:3000/version" -o /dev/null`
if [ "$STATUS_CODE" = "200" ];
then
    echo "PASS: backend /crons"
else
   echo "FAIL: backend /crons returned ${STATUS_CODE}, expected 200" ;
   exit 1
fi

echo "checking /crons/{account} endpoint"
curl -s localhost:3000/crons/myaccount | jq . || exit 1


# # # # # # #
# Test logs by parsing their contents
# # # # # # #
echo "checking frontend log"
grep "healthcheck from: 127.0.0.1" frontend_log | wc -l | grep 1 || exit 1
grep "Heartbeat from mycronjob: myaccount" frontend_log | wc -l | grep 2 || exit 1
grep "GET from 127.0.0.1 not valid" frontend_log | wc -l | grep 1 || exit 1
echo "PASS: frontend log"
echo "checking backend log"
grep "Checking for missed jobs" backend_log || exit 1
grep "mycronjob: myaccount failed to check in" backend_log || exit 1
grep '{"channel":"test","text":"The cronjob mycronjob for account myaccount has not checked in on time"}' backend_log || exit 1
grep "gocron success: alert for mycronjob sent" backend_log | wc -l | grep 2 || exit 1
grep "Alert for mycronjob: myaccount has been supressed. Already alerted" backend_log || exit 1
grep "mycronjob: myaccount is back online" backend_log || exit 1
grep '{"channel":"test","text":"The cronjob mycronjob for account myaccount is back online"}' backend_log || exit 1
grep "healthcheck from: 127.0.0.1" backend_log | wc -l | grep 1 || exit 1
echo "PASS: backend log"
