#!/bin/bash
cd $(dirname $0)

sudo git pull
sudo git status
read -p "Press [Enter] key to upgrade gocron on this branch"

sudo service gocron-front stop
sudo service gocron-back stop

sudo rm /usr/local/bin/gocro*
sudo cp ./bin/gocro* /usr/local/bin/

sudo service gocron-front restart
sudo service gocron-back restart

sudo service gocron-front status
sudo service gocron-back status
