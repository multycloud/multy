resource "aws_vpc" "example_vn_aws" {
  tags = {
    "Name" = "example_vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags = {
    "Name" = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags = {
    "Name" = "example_vn"
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
resource "aws_route_table" "rt_aws" {
  tags = {
    "Name" = "test-rt"
  }

  vpc_id = aws_vpc.example_vn_aws.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "rta_aws" {
  subnet_id      = aws_subnet.subnet_aws.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_subnet" "subnet_aws" {
  tags = {
    "Name" = "subnet"
  }

  cidr_block        = "10.0.2.0/24"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
}
resource "aws_key_pair" "vm_aws" {
  tags = {
    "Name" = "test-vm"
  }

  key_name   = "vm_aws_multy"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP"
}
resource "aws_iam_role" "vm_aws" {
  tags               = { "Name" = "test-vm" }
  name               = "iam_for_vm_vm_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_instance_profile" "vm_aws" {
  name = "iam_for_vm_vm_aws"
  role = aws_iam_role.vm_aws.name
}
resource "aws_instance" "vm_aws" {
  tags = {
    "Name" = "test-vm"
  }

  ami                         = "ami-09d4a659cdd8677be"
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.subnet_aws.id
  user_data_base64            = base64encode("#!/bin/bash -xe\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on AWS</h1>\" > /var/www/html/index.html")
  key_name                    = aws_key_pair.vm_aws.key_name
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
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
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "test-rt"
  location            = "northeurope"

  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "azurerm_subnet_route_table_association" "rta_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    public_ip_address_id          = azurerm_public_ip.vm_azure.id
    primary                       = true
  }
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "northeurope"
  allocation_method   = "Static"
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  resource_group_name   = azurerm_resource_group.vm-rg.name
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]
  custom_data           = base64encode("#!/bin/bash -xe\nsudo su\n yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \"<h1>Hello from Multy on Azure</h1>\" > /var/www/html/index.html;")

  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }

  admin_username = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP"
  }

  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7_9-gen2"
    version   = "latest"
  }
  identity {
    type = "SystemAssigned"
  }
  disable_password_authentication = true
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
