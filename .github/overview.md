## Overview

This repository is the brain of Multy which translates Multy resources into cloud specific infrastructure.

The Multy Engine is a GRPC server that translates a multy infrastructure resource request into cloud specific resources.

## File Structure

```bash
.
├── api # GRPC server
│    ├── aws # aws backend 
│    ├── deploy # runs terraform to refresh/plan/apply the translated config
│    ├── errors # common API errors returned 
│    ├── proto # contains all proto definitions and services exposed by the GRPC server 
│    │    ├── commonpb
│    │    ├── configpb
│    │    ├── credspb
│    │    ├── errorspb
│    │    └── resourcespb
│    ├── services # generic implementation of all resource services
│    └── util
├── cmd # command line tool to list/delete resources and start the server
├── db # database with user API keys and locks
├── encoder # translates all resources into terraform
├── flags
├── mhcl # implements custom go tag processors
├── resources
│    ├── common
│    ├── output # cloud specific resource outputs
│    └── types # contains all translation logic for every multy resource
├── test
│    ├── _configs # resource unit tests
│    └── e2e # end-to-end testing (deployment and test)
├── util 
└── validate # validation errors

```

## Technologies

- Golang (>=1.18)
- Terraform (>=1.0) - Backend for resource deployment
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

    1. Terraform - https://learn.hashicorp.com/tutorials/terraform/install-cli
    2. AWS CLI - https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

- Setup environment

  Setup your AWS account, adding your credentials locally by running `aws configure`.

  Create an Amazon S3 bucket

- Build project (and generate protos)

  ```bash
  make build
  ```

3. Run server locally

```bash
go run main.go serve --env=local
```

You can add the `--dry_run` flag when running the server. Dry run mode will work normally except it will not deploy any
resources.

4. Deploy configuration

You can find some examples of infrastructure configuration on
the [Terraform provider](https://github.com/multycloud/terraform-provider-multy/tree/main/tests)

Check the [docs](https://docs.multy.dev/getting-started) for more details.

5. Run tests

To run unit tests without running terraform, run:

```bash
go test ./test/... -v 
```

To also test that `terraform plan` works correctly on the generated configs, you can run:

```bash
go test ./test/... -v --tags=plan .
```
