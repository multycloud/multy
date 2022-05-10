## Overview

This repository is the brain of Multy which translates Multy resources into cloud specific infrastructure.

The Multy Engine is a GRPC server that translates a multy infrastructure resource request into cloud specific resources.

## File Structure

```bash
.
├── api # GRPC server
│    ├── aws # aws backend 
│    ├── converter
│    ├── deploy
│    ├── errors
│    ├── proto
│    │    ├── commonpb
│    │    ├── configpb
│    │    ├── credspb
│    │    ├── errorspb
│    │    └── resourcespb
│    ├── services
│    │    ├── database
│    │    ├── kubernetes_cluster
│    │    ├── kubernetes_node_pool
│    │    ├── lambda
│    │    ├── network_interface
│    │    ├── network_security_group
│    │    ├── object_storage
│    │    ├── object_storage_object
│    │    ├── public_ip
│    │    ├── route_table
│    │    ├── route_table_association
│    │    ├── subnet
│    │    ├── vault
│    │    ├── virtual_machine
│    │    └── virtual_network
│    └── util
├── cli
├── db # database with User API Keys and Locks
├── encoder
├── flags
├── mhcl
├── resources
│    ├── common
│    ├── output # cloud specific resource outputs
│    │    ├── database
│    │    ├── iam
│    │    ├── kubernetes_node_pool
│    │    ├── kubernetes_service
│    │    ├── lambda
│    │    ├── local_exec
│    │    ├── network_interface
│    │    ├── network_security_group
│    │    ├── object_storage
│    │    ├── object_storage_object
│    │    ├── provider
│    │    ├── public_ip
│    │    ├── route_table
│    │    ├── route_table_association
│    │    ├── subnet
│    │    ├── terraform
│    │    ├── vault
│    │    ├── vault_access_policy
│    │    ├── vault_secret
│    │    ├── virtual_machine
│    │    └── virtual_network
│    ├── resource_group
│    ├── tags
│    └── types
├── test
│    ├── _configs # resource unit tests
│    │    ├── database
│    │    ├── functions
│    │    ├── kubernetes
│    │    ├── lambda
│    │    ├── network_interface
│    │    ├── network_security_group
│    │    ├── object_storage
│    │    ├── object_storage_object
│    │    ├── public_ip
│    │    ├── resource_group
│    │    ├── route_table
│    │    ├── subnet
│    │    ├── vault
│    │    ├── vault_access_policy
│    │    ├── vault_secret
│    │    ├── virtual_machine
│    │    └── virtual_network
│    └── e2e # end-to-end testing (deployment and test)
│        ├── database
│        └── kubernetes
├── util
└── validate

```

## Technologies

- Golang (>1.18)
- Terraform - Backend for resource deployment
- MySQL - Store API keys and locks
- Amazon S3 - Store internal user configuration (WIP to remove dependency)
- GRPC

## Running locally

1. Clone repository

```bash
git clone https://github.com/multycloud/multy.git
cd multy
```

2. Setup project

- Install dependencies

MySQL

- Installation guide: https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/
- Run `./db/init.sql` script
- To create a new local Multy API Key
  run `INSERT INTO multydb.ApiKeys (ApiKey, UserId) VALUES ("test-key", "test-user");`

Terraform - https://learn.hashicorp.com/tutorials/terraform/install-cli

- Setup environment

Setup your AWS account, adding your credentials locally by running `aws configure`.

Create an Amazon S3 bucket

- Build project

```bash
make build
```

- Set environment variables

```bash
export MULTY_API_ENDPOINT="root:@tcp(localhost:3306)/multydb?parseTime=true;"
export MULTY_DB_CONN_STRING="localhost"
export USER_STORAGE_NAME=#YOUR_S3_BUCKET_NAME#
```

3. Run server

```bash
go run main.go serve 
```

You can add the `--dry_run` flag when running the server. Dry run mode will work normally except it will not deploy any
resources.

4. Deploy configuration

You can find some examples of infrastructure configuration on
the [Terraform provider](https://github.com/multycloud/terraform-provider-multy/tree/main/tests)

Check the [docs](https://docs.multy.dev/getting-started) for more details.

5. Run tests

```bash

```
