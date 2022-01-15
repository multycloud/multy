resource "aws_eip" "example_ip_aws" {
  tags =  {
    Name = "example_ip"
  }
}
resource "aws_network_interface" "example_nic_aws" {
  tags =  {
    Name = "example_nic"
  }

  subnet_id = "${aws_subnet.example_subnet_aws.id}"
}
resource "aws_subnet" "example_subnet_aws" {
  tags =  {
    Name = "example_subnet"
  }

  cidr_block = "10.0.0.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_vpc" "example_vn_aws" {
  tags =  {
    Name = "example_vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags =  {
    Name = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags =  {
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
resource "aws_s3_bucket" "obj_storage_aws" {
  bucket   = "test-storage"
  provider = "aws.eu-west-2"
}
resource "azurerm_public_ip" "example_ip_azure" {
  resource_group_name = azurerm_resource_group.pip-rg.name
  name                = "example_ip"
  location            = "northeurope"
  allocation_method   = "Static"
}
resource "azurerm_network_interface" "example_nic_azure" {
  resource_group_name = azurerm_resource_group.nic-rg.name
  name                = "example_nic"
  location            = "northeurope"

  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.example_subnet_azure.id}"
    primary                       = true
  }
}
resource "azurerm_subnet" "example_subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "example_subnet"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "example_subnet_azure" {
  subnet_id      = "${azurerm_subnet.example_subnet_azure.id}"
  route_table_id = azurerm_route_table.example_vn_azure.id
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
resource "azurerm_resource_group" "nic-rg" {
  name     = "nic-rg"
  location = "northeurope"
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name      = azurerm_resource_group.st-rg.name
  name                     = "teststorage"
  location                 = "ukwest"
  account_tier             = "Standard"
  account_replication_type = "GZRS"
  allow_blob_public_access = true
}
resource "azurerm_storage_container" "obj_storage_azure_public" {
  name                  = "public"
  storage_account_name  = azurerm_storage_account.obj_storage_azure.name
  container_access_type = "blob"
}
resource "azurerm_storage_container" "obj_storage_azure_private" {
  name                  = "private"
  storage_account_name  = azurerm_storage_account.obj_storage_azure.name
  container_access_type = "private"
}
resource "azurerm_resource_group" "pip-rg" {
  name     = "pip-rg"
  location = "northeurope"
}
resource "azurerm_resource_group" "st-rg" {
  name     = "st-rg"
  location = "northeurope"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "aws" {
  region = "eu-west-2"
  alias  = "eu-west-2"
}
provider "azurerm" {
  features {}
}
