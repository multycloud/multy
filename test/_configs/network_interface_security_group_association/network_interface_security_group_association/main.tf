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
resource "aws_eip" "ip_aws" {
  tags     = { "Name" = "test-ip" }
  provider = "aws.eu-west-1"
}
resource "azurerm_public_ip" "ip_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-ip"
  location            = "northeurope"
  allocation_method   = "Static"
  sku                 = "Standard"
}
resource "aws_eip_association" "nic_aws" {
  allocation_id        = aws_eip.ip_aws.id
  network_interface_id = aws_network_interface.nic_aws.id
  provider             = "aws.eu-west-1"
}
resource "aws_network_interface" "nic_aws" {
  tags      = { "Name" = "test-nic" }
  subnet_id = aws_subnet.subnet_aws-1.id
  provider  = "aws.eu-west-1"
}
resource "azurerm_network_interface" "nic_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-nic"
  location            = "northeurope"
  ip_configuration {
    name                          = "external-test-nic"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    public_ip_address_id          = azurerm_public_ip.ip_azure.id
    primary                       = true
  }
}
resource "aws_network_interface_sg_attachment" "nic_nsg_association_aws" {
  network_interface_id = aws_network_interface.nic_aws.id
  security_group_id    = aws_security_group.nsg2_aws.id
  provider             = "aws.eu-west-1"
}
resource "azurerm_network_interface_security_group_association" "nic_nsg_association_azure" {
  network_security_group_id = azurerm_network_security_group.nsg2_azure.id
  network_interface_id      = azurerm_network_interface.nic_azure.id
}
resource "aws_security_group" "nsg2_aws" {
  tags        = { "Name" = "test-nsg2" }
  vpc_id      = aws_vpc.example_vn_aws.id
  name        = "test-nsg2"
  description = "Managed by Multy"
  ingress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/16"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "tcp"
    from_port   = 443
    to_port     = 443
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["10.0.0.0/16"]
  }
  provider = "aws.eu-west-1"
}
resource "azurerm_network_security_group" "nsg2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-nsg2"
  location            = "northeurope"
  security_rule {
    name                       = "0-Inbound"
    protocol                   = "Tcp"
    priority                   = 120
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "22-22"
    destination_address_prefix = "*"
    direction                  = "Inbound"
  }
  security_rule {
    name                       = "0-Outbound"
    protocol                   = "Tcp"
    priority                   = 120
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "22-22"
    destination_address_prefix = "*"
    direction                  = "Outbound"
  }
  security_rule {
    name                       = "1-Inbound"
    protocol                   = "Tcp"
    priority                   = 140
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "443-443"
    destination_address_prefix = "*"
    direction                  = "Inbound"
  }
  security_rule {
    name                       = "1-Outbound"
    protocol                   = "Tcp"
    priority                   = 140
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "443-443"
    destination_address_prefix = "*"
    direction                  = "Outbound"
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
