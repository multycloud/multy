<p align="center">
<a href="https://multy.dev?utm_source=github.com">
    <img src="https://multy.dev/assets/multy_logo_horizontal.jpg" width="250">
</a>
</p>


<h3 align="center">
Multy is the easiest way to deploy multi cloud infrastructure
</h3>


<br>
<br>
<p align="center">
<img src="https://multy.dev/assets/multy-diagram.png" width="500">
</p>

<br>

<div align="center">

[![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=flat-square&logo=terraform&logoColor=white)](https://registry.terraform.io/providers/multycloud/multy/latest/docs)
[![Discord](https://img.shields.io/badge/Multy-%237289DA.svg?style=flat-square&logo=discord&logoColor=white)](https://discord.gg/rgaKXY4tCZ)
[![Build Status](https://img.shields.io/github/workflow/status/multycloud/multy/go_test?label=tests&style=flat-square)](https://github.com/multycloud/multy/actions)
</div>

# What is Multy?

**Multy** is an open-source tool that makes your infrastructure portable using a cloud-agnostic API.
You write your cloud-agnostic configuration once and Multy deploys it to the clouds you choose.

With Multy, you don't need to worry about how resources behave differently in the different clouds providers.
We abstract the nuances of each cloud so that moving your infrastructure between clouds is done by simply changing
the `cloud` parameter.

# Example

Let's try to deploy a simple virtual machine into AWS and Azure using
the [Multy Terraform Provider](https://github.com/multycloud/terraform-provider-multy)

```hcl
resource multy_virtual_network vn {
  name       = "test_vm"
  cidr_block = "10.0.0.0/16"
  cloud      = "aws"
  location   = "eu_west_1"
}

resource multy_subnet subnet {
  name               = "test_vm"
  cidr_block         = "10.0.10.0/24"
  virtual_network_id = multy_virtual_network.vn.id
}
```

By using the Multy cloud-agnostic API, we can simply change the `cloud` parameter to move a resource from one cloud to
another.

If we were to deploy this using the respective cloud terraform providers, we would first need to understand how
resources such as `aws_vpc` and `azurerm_virtual_network` behave and how they differ. Then we would need to define the
same infrastructure configuration twice, one for AWS and another for Azure.


<details><summary>This is the equivalent terraform configuration</summary>
<p>

```hcl
// terraform: 132 lines
resource "aws_vpc" "example_vn_aws" {
  tags = {
    Name = "example_vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags = {
    Name = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags = {
    Name = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id

  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }
}
resource "aws_subnet" "subnet_aws" {
  tags = {
    Name = "subnet"
  }

  cidr_block = "10.0.2.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_instance" "vm_aws" {
  tags = {
    Name = "test-vm"
  }

  ami           = "ami-09d4a659cdd8677be"
  instance_type = "t2.nano"
  subnet_id     = aws_subnet.subnet_aws.id
}
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "northeurope"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "northeurope"

  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    primary                       = true
  }
}
resource "random_password" "vm_azure" {
  length  = 16
  special = true
  upper   = true
  lower   = true
  number  = true
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  resource_group_name   = azurerm_resource_group.vm-rg.name
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls1"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]

  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }

  admin_username = "adminuser"
  admin_password = random_password.vm_azure.result

  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7_9-gen2"
    version   = "latest"
  }

  disable_password_authentication = false
}
resource "azurerm_resource_group" "vm-rg" {
  name     = "vm-rg"
  location = "northeurope"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
```

</p>
</details>

With Multy, **you write once, and deploy anywhere**.

## Getting started

1. Install Terraform - [see guide](https://learn.hashicorp.com/tutorials/terraform/install-cli#install-terraform), e.g.:

- Brew (Homebrew/Mac OS): `brew tap hashicorp/tap && brew install hashicorp/tap/terraform`
- Choco (Chocolatey/Windows): `choco install terraform`
- Debian (Ubuntu/Linux):
  ```
  sudo apt-get update && sudo apt-get install -y gnupg software-properties-common curl
  curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
  sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
  sudo apt-get update && sudo apt-get install terraform
  ```

2. Create an account with AWS or Azure and expose
   its [authentication credentials via environment variables](https://docs.multy.dev/getting-started#3-generate-access-keys)

3. Write your configuration file, for example a file named `main.tf` with the following content:

    ```hcl
    terraform {
      required_providers {
        multy = {
          source = "multycloud/multy"
        }
      }
    }
    
    provider "multy" {
      aws = {} # this will look for aws credentials
    }
    resource "multy_virtual_network" "vn" {
      cloud = "aws"
    
      name       = "multy_vn"
      cidr_block = "10.0.0.0/16"
      location   = "eu_west_1"
    }
    ```

4. Run `terraform init` and then `terraform apply`
5. Run `terraform destroy`

For a more detailed guide, see our official [getting started guide](https://docs.multy.dev/getting-started).

## Contributing

We love contributors! If you're interested in contributing, take a look at our [Contributing guide](./CONTRIBUTING.md)
. <br>
Join our [discord](https://discord.gg/rgaKXY4tCZ) channel to participate in live discussions or ask for support.

Repo overview: [OVERVIEW.md](./.github/overview.md)

Terraform Provider
Repo: [https://github.com/multycloud/terraform-provider-multy](https://github.com/multycloud/terraform-provider-multy?ref=multy-gh-repo)

Discord Channel: [https://discord.gg/rgaKXY4tCZ](https://discord.gg/rgaKXY4tCZ)

## Roadmap

Have a look at our [roadmap](https://github.com/orgs/multycloud/projects/4) to know the latest features released and
what we're focusing on short and long term.
You can also vote for a specific feature you want or participate in
the [discussions](https://github.com/multycloud/multy/discussions).

## FAQ

### Why build with Multy?

Multy was born after realising how difficult it is to run the same infrastructure across multiple clouds. While
providers such as AWS and Azure share the same set of core services, the small differences in how each service works
make it difficult to configure your infrastructure to run in the same way.

This is the problem that Multy aims to tackle. We created a single interface to deploy resources that have the same
behaviour regardless of the cloud provider.

### Can I use Multy for free?

Multy is available as a free and open-source tool, so you can download it directly and run it locally.

We also offer a managed solution that hosts the server for you. Managed Multy is currently offered as a free service.
You can request an API key by visiting our [website](https://multy.dev).

### Why not use the cloud specific Terraform providers?

While Terraform and its providers are great for deploying any resource into any cloud, it puts all the burden on the
infrastructure teams when it comes to understanding each provider and defining the resources. This flexibility can be
seen as an advantage, however, when it comes to multi-cloud, this considerably slows down teams that are looking to move
fast with deployments.

By abstracting the common resources across major cloud providers, users are able to deploy the same resources on AWS and
Azure without re-writing any infrastructure code.

### I want to use cloud managed resources (i.e. Amazon S3 / Azure Key Vault), is Multy for me?

Absolutely! The goal with Multy is to allow you to leverage cloud managed services and remain free to move your
infrastructure.
Not every resource will be supported, but we aim to support the most popular managed resources such as managed
databases, object storage and vault.

Let us know what services you would like to be supported by creating an Issue on
the [Issues section](https://github.com/multycloud/multy/issues).

### Why should I be locked in to Multy?

Multy is an open-source tool that can be run locally and free. If at some point you want to move off Multy, you can
export your infrastructure configuration as Terraform and use it independently.

## License

This repository is available under [Apache 2.0](./LICENSE).