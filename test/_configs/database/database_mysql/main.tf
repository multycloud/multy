resource "aws_db_subnet_group" "example_db_aws" {
  tags = {
    "Name" = "example-db"
  }

  name = "example-db"

  subnet_ids = [
    "${aws_subnet.subnet1_aws.id}",
    "${aws_subnet.subnet2_aws.id}",
  ]
}
resource "aws_db_instance" "example_db_aws" {
  tags = {
    "Name" = "exampledb"
  }

  allocated_storage    = 10
  db_name              = "exampledb"
  engine               = "mysql"
  engine_version       = "5.7"
  username             = "multyadmin"
  password             = "multy$Admin123!"
  instance_class       = "db.t2.micro"
  identifier           = "example-db"
  skip_final_snapshot  = true
  db_subnet_group_name = aws_db_subnet_group.example_db_aws.name
  publicly_accessible  = true
}
resource "aws_subnet" "subnet1_aws" {
  tags = {
    "Name" = "subnet1"
  }

  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "eu-west-1a"
}
resource "aws_subnet" "subnet2_aws" {
  tags = {
    "Name" = "subnet2"
  }

  cidr_block        = "10.0.1.0/24"
  vpc_id            = aws_vpc.vn_aws.id
  availability_zone = "eu-west-1b"
}
resource "aws_route_table" "rt_aws" {
  tags = {
    "Name" = "test-rt"
  }

  vpc_id = "${aws_vpc.vn_aws.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vn_aws.id
  }
}
resource "aws_route_table_association" "rta_aws" {
  subnet_id      = "${aws_subnet.subnet1_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_route_table_association" "rta2_aws" {
  subnet_id      = "${aws_subnet.subnet2_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_vpc" "vn_aws" {
  tags = {
    "Name" = "vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "vn_aws" {
  tags = {
    "Name" = "vn"
  }

  vpc_id = aws_vpc.vn_aws.id
}
resource "aws_default_security_group" "vn_aws" {
  tags = {
    "Name" = "vn"
  }

  vpc_id = aws_vpc.vn_aws.id

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
resource "azurerm_resource_group" "db-rg" {
  name     = "db-rg"
  location = "northeurope"
}
resource "azurerm_mysql_server" "example_db_azure" {
  resource_group_name          = azurerm_resource_group.db-rg.name
  name                         = "example-db"
  location                     = "northeurope"
  administrator_login          = "multyadmin"
  administrator_login_password = "multy$Admin123!"
  sku_name                     = "GP_Gen5_2"
  storage_mb                   = 10240
  version                      = "5.7"
  ssl_enforcement_enabled      = false
}
resource "azurerm_mysql_virtual_network_rule" "example_db_azure0" {
  resource_group_name = azurerm_resource_group.db-rg.name
  name                = "example-db0"
  server_name         = azurerm_mysql_server.example_db_azure.name
  subnet_id           = "${azurerm_subnet.subnet1_azure.id}"
}
resource "azurerm_mysql_virtual_network_rule" "example_db_azure1" {
  resource_group_name = azurerm_resource_group.db-rg.name
  name                = "example-db1"
  server_name         = azurerm_mysql_server.example_db_azure.name
  subnet_id           = "${azurerm_subnet.subnet2_azure.id}"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_virtual_network" "vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "vn"
  location            = "northeurope"
  address_space       = ["10.0.0.0/16"]
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
  subnet_id      = "${azurerm_subnet.subnet1_azure.id}"
  route_table_id = "${azurerm_route_table.rt_azure.id}"
}
resource "azurerm_subnet_route_table_association" "rta2_azure" {
  subnet_id      = "${azurerm_subnet.subnet2_azure.id}"
  route_table_id = "${azurerm_route_table.rt_azure.id}"
}
resource "azurerm_route_table" "vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "vn"
  location            = "northeurope"

  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
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
