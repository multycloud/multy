#!/bin/bash -xe

{
export USERS_S3_BUCKET_NAME='${s3_bucket_name}'

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

git clone https://github.com/multycloud/multy.git
cd multy
go get -u google.golang.org/protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
go install golang.org/x/lint/golint@latest
go mod tidy
make run

} |& tee -a logs.txt