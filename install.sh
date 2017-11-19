#!/bin/bash
cd $(dirname $0)

# Get gocron binary
sudo mkdir /usr/local/bin
sudo cp ./bin/gocron /usr/local/bin
sudo chmod +x /usr/local/bin/gocron

# Get gocron config and configure it
# "/etc/gocron/.config.yml"
sudo mkdir -p /etc/gocron/
sudo cp ./src/example.config.yml /etc/gocron/config.yml
sudo chmod 600 /etc/gocron/config.yml
sudo vi /etc/gocron/config.yml

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
sudo systemctl enable gocron.service && sudo systemctl start gocron.service
sudo systemctl status gocron.service
