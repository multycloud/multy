#!/bin/bash -xe

{
export USERS_S3_BUCKET_NAME='${s3_bucket_name}'
export MULTY_DB_CONN_STRING='${db_connection}'
export MULTY_API_ENDPOINT='${api_endpoint}'

mkdir -p "$HOME/.terraform.d/plugin-cache"
echo plugin_cache_dir = \"$HOME/.terraform.d/plugin-cache\" > "$HOME/.terraformrc"
sudo apt-get update -y && sudo apt-get install -y gnupg software-properties-common curl
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update -y
sudo apt-get -y install terraform git make protobuf-compiler mysql-client awscli

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

sudo -E aws s3 cp s3://multy-internal/certs "/etc/letsencrypt/live/${MULTY_API_ENDPOINT}/" --recursive --exclude="*" --include "fullchain.pem" --include "privkey.pem"

git clone https://github.com/multycloud/multy.git
cd multy
go get -u google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
go install golang.org/x/lint/golint@latest
go mod tidy
make run PORT=443
} |& tee -a logs.txt