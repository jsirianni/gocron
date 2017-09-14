#!/bin/bash

# This script will install and configure go cron
# with postgresql. Run as the user that will be
# running the service as the config file will live in
# ~/.config/gocron.

# This script has been dested on Debian 9 with postgres 9.6
# but may work with Debian 8 and Ubuntu 16.04.


# Adjustable values
password=$1


# Install postgres
sudo apt-get update
sudo apt-get install postgresql-9.6 nano -y
sudo systemctl enable postgresql
sudo systemctl start postgresql


# Configure postgres
echo "CREATE DATABASE gocron;"
sudo -u postgres psql -c "CREATE DATABASE gocron;"
sudo -u postgres psql gocron -c "CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency varchar, tolerance int, lastruntime varchar, alerted boolean, PRIMARY KEY(cronname, account));"
echo "CREATE TABLE gocron(cronName varchar, account varchar, email varchar, ipaddress varchar, frequency varchar, tolerance int, lastruntime varchar, alerted boolean, PRIMARY KEY(cronname, account));"
sudo -u postgres psql gocron -c "CREATE USER gocron WITH PASSWORD '$password';"
echo "CREATE USER gocron WITH PASSWORD '$password';"
sudo -u postgres psql gocron -c "GRANT ALL PRIVILEGES ON gocron TO gocron;"
echo "GRANT ALL PRIVILEGES ON gocron TO gocron;"
sleep 10

# Pull newest Build
sudo git pull

# Get gocron binary
sudo mkdir /usr/local/bin
sudo cp ./bin/gocron /usr/local/bin
sudo chmod +x /usr/local/bin/gocron

# Get gocron config and configure it
# "/etc/gocron/.config.yml"
sudo mkdir -p /etc/gocron/
sudo cp ./src/.example.config.yml /etc/gocron/.config.yml
sudo chmod 600 /etc/gocron/.config.yml
sudo nano /etc/gocron/.config.yml
sleep 5

# Build systemd service
sudo touch /etc/systemd/system/gocron.service

echo "[Unit]" > /etc/systemd/system/gocron.service
echo "Description=GOCron Monitoring Service" >> /etc/systemd/system/gocron.service
echo "After=network.target" >> /etc/systemd/system/gocron.service

echo "[Service]" >> /etc/systemd/system/gocron.service
echo "ExecStart=/usr/local/bin/gocron" >> /etc/systemd/system/gocron.service

echo "[Install]" >> /etc/systemd/system/gocron.service
echo "WantedBy=multi-user.target" >> /etc/systemd/system/gocron.service
sleep 5

# Enable the gocron service
sudo systemctl enable gocron.service
sudo systemctl start gocron.service
sudo systemctl status gocron.service
