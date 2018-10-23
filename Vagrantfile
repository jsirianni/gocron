$script = <<-SCRIPT
sudo apt-get update
sudo apt-get install golang -y
sudo apt-get install postgresql -y

sudo apt-get install docker.io -y
sudo curl -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo docker swarm init

sudo echo "listen_addresses='*'" | tee -a /etc/postgresql/10/main/postgresql.conf
sudo echo "host all all 0.0.0.0/0 trust" | tee -a /etc/postgresql/10/main/pg_hba.conf
sudo service postgresql restart

sudo -u postgres createuser gocron
sudo -u postgres createdb gocron
sudo -u postgres -H -- psql -c "alter user gocron with encrypted password 'password'"
sudo -u postgres -H -- psql -c "grant all privileges on database gocron to gocron"

export GOPATH=/home/vagrant
mkdir $GOPATH/src
sudo cp -r /gocron $GOPATH/src/
sudo chown -R vagrant:vagrant $GOPATH
cd $GOPATH/src/gocron
sudo mkdir /etc/gocron
sudo cp example.config.yml /etc/gocron/config.yml
go get
./build_linux.sh
cd docker
sudo docker stack deploy gocron --compose-file docker-compose.yml
SCRIPT



Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-18.04"
  config.vm.synced_folder "./", "/gocron"
  config.vm.provision "shell", inline: $script
end
