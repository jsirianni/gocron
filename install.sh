#!/bin/bash
cd $(dirname $0)


# Get gocron config and configure it
# "/etc/gocron/config.yml"
sudo mkdir -p /etc/gocron/
sudo cp example.config.yml /etc/gocron/config.yml
sudo chmod 600 /etc/gocron/config.yml
sudo vi /etc/gocron/config.yml


# Build systemd service
sudo touch /etc/systemd/system/gocron-front.service
tee /etc/systemd/system/gocron-front.service << EOH
[Unit]
Description=GoCron Monitoring Service - Frontend Gateway
After=network.target
[Service]
ExecStart=/usr/local/bin/gocron-front
[Install]
WantedBy=multi-user.target
EOH

sudo touch /etc/systemd/system/gocron-back.service
tee /etc/systemd/system/gocron-back.service << EOH
[Unit]
Description=GoCron Monitoring Service - Backend
After=network.target
[Service]
ExecStart=/usr/local/bin/gocron-back
[Install]
WantedBy=multi-user.target
EOH


# Get gocron binary
sudo mkdir /usr/local/bin
cd /usr/local/bin
sudo wget https://github.com/jsirianni/gocron/releases/download/3.0.2/gocron-back
sudo wget https://github.com/jsirianni/gocron/releases/download/3.0.2/gocron-front
sudo chmod +x /usr/local/bin/gocron-*


# Enable the gocron service
sudo systemctl enable gocron-front.service
sudo systemctl start gocron-front.service
sudo systemctl status gocron-front.service

sudo systemctl enable gocron-back.service
sudo systemctl start gocron-back.service
sudo systemctl status gocron-back.service
