resource "aws_vpc" "example_vn_aws" {
  tags                 = { "Name" = "example_vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  provider             = "aws.eu-west-1"
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags     = { "Name" = "example_vn" }
  vpc_id   = aws_vpc.example_vn_aws.id
  provider = "aws.eu-west-1"
}
resource "aws_default_security_group" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
  vpc_id = aws_vpc.example_vn_aws.id
  ingress {
    protocol  = "-1"
    from_port = 0
    to_port   = 0
    self      = true
  }
  egress {
    protocol  = "-1"
    from_port = 0
    to_port   = 0
    self      = true
  }
  provider = "aws.eu-west-1"
}
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example_vn"
  location            = "northeurope"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example_vn"
  location            = "northeurope"
  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "aws_network_interface" "nic_aws" {
  tags      = { "Name" = "test-nic" }
  subnet_id = aws_subnet.subnet_aws-2.id
  provider  = "aws.eu-west-1"
}
resource "azurerm_network_interface" "nic_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-nic"
  location            = "northeurope"
  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    primary                       = true
  }
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
resource "aws_subnet" "subnet_aws-1" {
  tags              = { "Name" = "subnet-1" }
  cidr_block        = "10.0.2.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet_aws-2" {
  tags              = { "Name" = "subnet-2" }
  cidr_block        = "10.0.2.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet_aws-3" {
  tags              = { "Name" = "subnet-3" }
  cidr_block        = "10.0.2.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}
provider "azurerm" {
  features {
  }
}
