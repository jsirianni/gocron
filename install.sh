#!/bin/bash

# This script will install and configure go cron
# with postgresql. Run as the user that will be
# running the service as the config file will live in
# ~/.config/gocron.

# This script has been dested on Debian 9 with postgres 9.6
# but may work with Debian 8 and Ubuntu 16.04.


# Adjustable values
user=gocron
password=$1


# Install postgres
sudo apt-get update
sudo apt-get install postgresql-9.6 wget nano -y
sudo systemctl enable postgresql
sudo systemctl start postgresql


# Configure postgres
sudo -u postgres psql -c "CREATE DATABASE gocron;"
sudo -u postgres psql -c "CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency varchar, tolerance int, lastruntime varchar, alerted boolean, PRIMARY KEY(cronname, account));"
sudo -u postgres psql -c "CREATE USER $user WITH PASSWORD '$password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON gocron TO $user;"


# Get gocron binary
sudo mkdir /usr/local/bin
wget -O /usr/local/bin/gocron https://github.com/jsirianni/gocron/blob/master/bin/cronserver?raw=true


# Get gocron config and configure it
sudo mkdir -p ~/.config/gocron
wget -O ~/.config/gocron/.config.yml https://raw.githubusercontent.com/jsirianni/gocron/master/src/.example.config.yml
sudo chmod 600 ~/.config/gocron/.config.yml
sudo nano ~/.config/gocron/.config.yml


# Build systemd service
sudo touch /etc/systemd/system/gocron.service

echo "[Unit]" > /etc/systemd/system/gocron.service
echo "Description=GOCron Monitoring Service" >> /etc/systemd/system/gocron.service
echo "After=network.target" >> /etc/systemd/system/gocron.service

echo "[Service]" >> /etc/systemd/system/gocron.service
echo "ExecStart=/usr/local/bin/gocron" >> /etc/systemd/system/gocron.service

echo "[Install]" >> /etc/systemd/system/gocron.service
echo "WantedBy=multi-user.target" >> /etc/systemd/system/gocron.service


# Enable the gocron service
sudo systemctl enable gocron.service
sudo systemctl start gocron.service
sudo systemctl status gocron.service
