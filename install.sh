#!/bin/bash
export RepoName=`pwd`
export OSName=`awk -F= '/^NAME/{print $2}' /etc/os-release`
export APIUser=$USER

### Run updates according to the OS the API is running on ###
if [ "$OSName" = "CentOS Linux" ]
 then
    sudo yum -y update && sudo yum -y upgrade
    sudo yum install -y wget git-all
else
    sudo apt-get -y update && sudo apt-get -y upgrade
    sudo apt-get install -y git-all wget
fi

### Install and Configure Go v1.10.3 ###
cd ~
export HomeDir=`pwd`
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
sudo tar -C $HomeDir -xzf go1.10.3.linux-amd64.tar.gz

export PATH=$PATH:$RepoName/bin
export PATH=$PATH:$RepoName/pkg
export PATH=$PATH:$RepoName/src
export GOPATH=$RepoName

# Configure the Go Variables for both ubuntu and centos as needed
if [ "$OSName" = "CentOS Linux" ]
 then
    echo 'export PATH=$PATH:'$HOME'/go/bin' >> $HomeDir/.bash_profile
    echo 'export PATH=$PATH:'$HOME'/go/pkg' >> $HomeDir/.bash_profile
    echo 'export PATH=$PATH:'$HOME'/go/src' >> $HomeDir/.bash_profile
    echo 'export PATH=$PATH:'$RepoName'/bin' >> $HomeDir/.bash_profile
    echo 'export PATH=$PATH:'$RepoName'/pkg' >> $HomeDir/.bash_profile
    echo 'export PATH=$PATH:'$RepoName'/src' >> $HomeDir/.bash_profile
    echo 'export GOPATH='$RepoName >> $HomeDir/.bash_profile
    source $HomeDir/.bash_profile
else
    echo 'export PATH=$PATH:'$HOME'/go/bin' >> $HomeDir/.profile
    echo 'export PATH=$PATH:'$HOME'/go/pkg' >> $HomeDir/.profile
    echo 'export PATH=$PATH:'$HOME'/go/src' >> $HomeDir/.profile
    echo 'export PATH=$PATH:'$RepoName'/bin' >> $HomeDir/.profile
    echo 'export PATH=$PATH:'$RepoName'/pkg' >> $HomeDir/.profile
    echo 'export PATH=$PATH:'$RepoName'/src' >> $HomeDir/.profile
    echo 'export GOPATH='$RepoName >> $HomeDir/.profile
    source $HomeDir/.profile
fi

### Get APIs Go Dependencies ###
cd $RepoName

go get "github.com/gorilla/mux"
go get "github.com/gorilla/handlers"
go get "go.mongodb.org/mongo-driver/mongo"
go get "go.mongodb.org/mongo-driver/mongo/options"
go get "go.mongodb.org/mongo-driver/bson"
go get "go.mongodb.org/mongo-driver/bson/primitive"
go get "github.com/dgrijalva/jwt-go"
go get "github.com/gofrs/uuid"
go get "golang.org/x/crypto/bcrypt"

#mkdir logs

rm -f ~/go1.10.3.linux-amd64.tar.gz

### Install and Configure Local Testing Mongo ###
# Run updates according to the OS the API is running on
if [ "$OSName" = "CentOS Linux" ]
 then
    echo '[MongoDB]
name=MongoDB Repository
baseurl=http://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/4.0/x86_64/
gpgcheck=0
enabled=1' | sudo tee -a /etc/yum.repos.d/mongodb.repo
    sudo yum install -y mongodb-org
    sudo systemctl start mongod.service
    sudo systemctl enable mongod.service
else
    sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
    echo 'deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse' | sudo tee /etc/apt/sources.list.d/mongodb.list
    sudo apt -y update && sudo apt -y install mongodb-org
    sudo systemctl enable mongod
    sudo systemctl start mongod
fi

### Setup systemd service ###
if [ "$OSName" = "CentOS Linux" ]
 then
   sudo sh $RepoName/install/setup_service_centos.sh $RepoName $APIUser
else
  sudo sh $RepoName/install/setup_service_ubuntu.sh $RepoName $APIUser
fi
