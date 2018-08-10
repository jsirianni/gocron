#!/bin/bash
cd $(dirname $0)

sudo git pull
sudo git status
read -p "Press [Enter] key to upgrade gocron on this branch"

sudo service gocron-front stop
sudo service gocron-back stop

sudo rm /usr/local/bin/gocro*
cd /usr/local/bin
sudo wget https://github.com/jsirianni/gocron/releases/download/v3.0.4/gocron-back
sudo wget https://github.com/jsirianni/gocron/releases/download/v3.0.4/gocron-front
sudo chmod +x /usr/local/bin/gocron-*

sudo service gocron-front restart
sudo service gocron-back restart

sudo service gocron-front status
sudo service gocron-back status
