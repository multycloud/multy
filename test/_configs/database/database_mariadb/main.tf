resource "aws_db_subnet_group" "example_db_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "example-db"
  }

  name        = "example-db"
  description = "Managed by Multy"
  subnet_ids  = [
    aws_subnet.subnet1_aws.id,
    aws_subnet.subnet2_aws.id,
  ]
}
resource "aws_db_instance" "example_db_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "exampledb"
  }

  allocated_storage      = 10
  engine                 = "mariadb"
  engine_version         = "10.2"
  username               = "multyadmin"
  password               = "multy$Admin123!"
  instance_class         = "db.t2.micro"
  identifier             = "example-db"
  skip_final_snapshot    = true
  db_subnet_group_name   = aws_db_subnet_group.example_db_aws.name
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.example_db_aws.id]
}
resource "aws_security_group" "example_db_aws" {
  tags        = { "Name" = "example-db" }
  vpc_id      = aws_vpc.vn_aws.id
  name        = "example-db"
  description = "Default security group of example-db"
  ingress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
  provider = "aws.us-east-1"
}
resource "aws_subnet" "subnet1_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "subnet1"
  }

  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "us-east-1a"
}
resource "aws_subnet" "subnet2_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "subnet2"
  }

  cidr_block        = "10.0.1.0/24"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "us-east-1b"
}
resource "aws_route_table" "rt_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "db-rt"
  }

  vpc_id = aws_vpc.vn_aws.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vn_aws.id
  }
}
resource "aws_route_table_association" "subnet1_aws" {
  provider       = "aws.us-east-1"
  subnet_id      = aws_subnet.subnet1_aws.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_route_table_association" "subnet2_aws" {
  provider       = "aws.us-east-1"
  subnet_id      = aws_subnet.subnet2_aws.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_vpc" "vn_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "db-vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "vn_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "db-vn"
  }

  vpc_id = aws_vpc.vn_aws.id
}
resource "aws_default_security_group" "vn_aws" {
  provider = "aws.us-east-1"
  tags     = {
    "Name" = "db-vn"
  }

  vpc_id = aws_vpc.vn_aws.id

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
}
resource "azurerm_mariadb_server" "example_db_azure" {
  resource_group_name          = azurerm_resource_group.rg1.name
  name                         = "example-db"
  location                     = "eastus"
  administrator_login          = "multyadmin"
  administrator_login_password = "multy$Admin123!"
  sku_name                     = "GP_Gen5_2"
  storage_mb                   = 10240
  version                      = "10.2"
  ssl_enforcement_enabled      = false
}
resource "azurerm_mariadb_virtual_network_rule" "example_db_azure0" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db0"
  server_name         = azurerm_mariadb_server.example_db_azure.name
  subnet_id           = azurerm_subnet.subnet1_azure.id
}
resource "azurerm_mariadb_virtual_network_rule" "example_db_azure1" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db1"
  server_name         = azurerm_mariadb_server.example_db_azure.name
  subnet_id           = azurerm_subnet.subnet2_azure.id
}
resource "azurerm_mariadb_firewall_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "public"
  server_name         = azurerm_mariadb_server.example_db_azure.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_virtual_network" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-rt"
  location            = "eastus"

  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "azurerm_subnet_route_table_association" "subnet1_azure" {
  subnet_id      = azurerm_subnet.subnet1_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_subnet_route_table_association" "subnet2_azure" {
  subnet_id      = azurerm_subnet.subnet2_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_route_table" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"

  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "eastus"
}
provider "aws" {
  region = "us-east-1"
  alias  = "us-east-1"
}
provider "azurerm" {
  features {}
}

