#!/bin/bash

# This script will install and configure go cron
# with postgresql. Run as the user that will be
# running the service as the config file will live in
# ~/.config/gocron.

# This script has been dested on Debian 9 with postgres 9.6
# but may work with Debian 8 and Ubuntu 16.04.

# Get command line args or replace values in this script
user=gocron
password=$1

# Install postgres
apt-get update
apt-get install postgresql-9.6 -y
systemctl enable postgresql
systemctl start postgresql

# Configure postgres
sudo -u postgres psql -c "CREATE DATABASE gocron;"
sudo -u postgres psql -c "CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency varchar, tolerance int, lastruntime varchar, alerted boolean, PRIMARY KEY(cronname, account));"
sudo -u postgres psql -c "CREATE USER $user WITH PASSWORD '$password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON gocron TO $user;"

# Get gocron binary
