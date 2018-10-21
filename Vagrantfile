$script = <<-SCRIPT
sudo apt-get update
sudo apt-get install golang -y
sudo apt-get install postgresql -y
sudo apt-get install docker.io -y

sudo -u postgres createuser gocron
sudo -u postgres createdb gocron
sudo -u postgres -H -- psql -c "alter user gocron with encrypted password 'password'"
sudo -u postgres -H -- psql -c "grant all privileges on database gocron to gocron"


export GOPATH=/home/vagrant
mkdir $GOPATH/src
sudo cp -r /gocron $GOPATH/src/
sudo chown -R vagrant:vagrant $GOPATH
cd $GOPATH/src/gocron
go get 
./build_linux.sh
SCRIPT



Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-18.04"
  config.vm.synced_folder "./", "/gocron"
  config.vm.provision "shell", inline: $script  
end
