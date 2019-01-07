$script = <<-SCRIPT

# install packages
sudo sudo apt-get update && apt-get install golang postgresql docker.io -y

# install docker compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# init docker swarm
sudo docker swarm init

# configure postgres
sudo echo "listen_addresses='*'" | tee -a /etc/postgresql/10/main/postgresql.conf
sudo echo "host all all 0.0.0.0/0 trust" | tee -a /etc/postgresql/10/main/pg_hba.conf
sudo service postgresql restart
sudo -u postgres createuser gocron
sudo -u postgres createdb gocron
sudo -u postgres -H -- psql -c "alter user gocron with encrypted password 'password'"
sudo -u postgres -H -- psql -c "grant all privileges on database gocron to gocron"

# build the image
cd /gocron
sudo docker build -t gocron:latest .

# run the stack
cd /gocron/docker
sudo docker stack deploy gocron --compose-file docker-compose.yml

# end script
SCRIPT



Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.synced_folder "./", "/gocron"
  config.vm.provision "shell", inline: $script
end
