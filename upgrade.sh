#!/bin/bash
cd $(dirname $0)

sudo git pull
sudo git status
read -p "Press [Enter] key to upgrade gocron on this branch"

sudo service gocron-front stop
sudo service gocron-back stop

sudo rm /usr/local/bin/gocro*
sudo cp ./bin/gocro* /usr/local/bin/

sudo service gocron-front start
sudo service gocron-back start

sudo service gocron-front status
sudo service gocron-back status
