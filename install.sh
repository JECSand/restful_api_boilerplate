#!/bin/bash
export RepoName=`pwd`
sudo yum install -y wget git-all

### Install and Configure Go v1.10.3 ###
cd ~
export HomeDir=`pwd`
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
sudo tar -C $HomeDir -xzf go1.10.3.linux-amd64.tar.gz

export PATH=$PATH:$RepoName/bin
export PATH=$PATH:$RepoName/pkg
export PATH=$PATH:$RepoName/src
export GOPATH=$RepoName

echo 'export PATH=$PATH:'$HOME'/go/bin' >> $HomeDir/.bash_profile
echo 'export PATH=$PATH:'$HOME'/go/pkg' >> $HomeDir/.bash_profile
echo 'export PATH=$PATH:'$HOME'/go/src' >> $HomeDir/.bash_profile
echo 'export PATH=$PATH:'$RepoName'/bin' >> $HomeDir/.bash_profile
echo 'export PATH=$PATH:'$RepoName'/pkg' >> $HomeDir/.bash_profile
echo 'export PATH=$PATH:'$RepoName'/src' >> $HomeDir/.bash_profile
echo 'export GOPATH='$RepoName >> $HomeDir/.bash_profile

source $HomeDir/.bash_profile

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

mkdir logs

rm -f ~/go1.10.3.linux-amd64.tar.gz

### Install and Configure Local Testing Mongo ###
echo '[MongoDB]
name=MongoDB Repository
baseurl=http://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/4.0/x86_64/
gpgcheck=0
enabled=1' | sudo tee -a /etc/yum.repos.d/mongodb.repo

sudo yum install -y mongodb-org
sudo systemctl start mongod.service
sudo systemctl enable mongod.service

### Setup systemd service ###
sudo sh setup_service.sh $RepoName