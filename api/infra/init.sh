#!/bin/bash -xe

{
export USERS_S3_BUCKET_NAME='${s3_bucket_name}'

mkdir -p "$HOME/.terraform.d/plugin-cache"
echo plugin_cache_dir = \"$HOME/.terraform.d/plugin-cache\" > "$HOME/.terraformrc"
sudo apt-get update -y && sudo apt-get install -y gnupg software-properties-common curl
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update -y && sudo apt-get -y install terraform

sudo apt-get update -y
sudo apt-get -y install git make protobuf-compiler

wget https://golang.org/dl/go1.18.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOROOT=/usr/local/go
export GOCACHE=/root/go/cache
export GOPATH=~/.go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export HOME=/home/ubuntu
go version
go env -w GO111MODULE=on
sudo chmod -R 777 /root/

sudo snap install core
sudo snap refresh core
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot certonly --standalone --non-interactive --agree-tos -m systemalerts@multy.dev --domains api.multy.dev
sudo certbot renew --dry-run

git clone https://github.com/multycloud/multy.git
cd multy
go get -u google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
go install golang.org/x/lint/golint@latest
go mod tidy
make run PORT=443
} |& tee -a logs.txt