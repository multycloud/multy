resource "aws_db_subnet_group" "example_db_aws" {
  tags        = { "Name" = "example-db" }
  name        = "example-db"
  description = "Managed by Multy"
  subnet_ids  = [
    aws_subnet.private_subnet_aws-1.id, aws_subnet.private_subnet_aws-2.id, aws_subnet.private_subnet_aws-3.id
  ]
  provider = "aws.eu-west-1"
}
resource "aws_security_group" "example_db_aws" {
  tags        = { "Name" = "example-db" }
  vpc_id      = aws_vpc.example_vn_aws.id
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
  provider = "aws.eu-west-1"
}
resource "aws_db_instance" "example_db_aws" {
  tags                   = { "Name" = "exampledb" }
  allocated_storage      = 10
  engine                 = "mysql"
  engine_version         = "5.7"
  username               = "multyadmin"
  password               = "multy$Admin123!"
  instance_class         = "db.t2.micro"
  identifier             = "example-db"
  skip_final_snapshot    = true
  db_subnet_group_name   = aws_db_subnet_group.example_db_aws.name
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.example_db_aws.id]
  provider               = "aws.eu-west-1"
}
resource "azurerm_mysql_server" "example_db_azure" {
  resource_group_name              = azurerm_resource_group.rg1.name
  name                             = "example-db"
  location                         = "northeurope"
  administrator_login              = "multyadmin"
  administrator_login_password     = "multy$Admin123!"
  sku_name                         = "GP_Gen5_2"
  storage_mb                       = 10240
  version                          = "5.7"
  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
}
resource "azurerm_mysql_virtual_network_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db"
  server_name         = azurerm_mysql_server.example_db_azure.name
  subnet_id           = azurerm_subnet.private_subnet_azure.id
}
resource "azurerm_mysql_firewall_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "public"
  server_name         = azurerm_mysql_server.example_db_azure.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
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
resource "aws_security_group" "nsg2_aws" {
  tags        = { "Name" = "test-nsg2" }
  vpc_id      = aws_vpc.example_vn_aws.id
  name        = "test-nsg2"
  description = "Managed by Multy"
  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }
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
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
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
    name                       = "0"
    protocol                   = "Tcp"
    priority                   = 100
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "80-80"
    destination_address_prefix = "*"
    direction                  = "Inbound"
  }
  security_rule {
    name                       = "1"
    protocol                   = "Tcp"
    priority                   = 100
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "80-80"
    destination_address_prefix = "*"
    direction                  = "Outbound"
  }
  security_rule {
    name                       = "2"
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
    name                       = "3"
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
    name                       = "4"
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
    name                       = "5"
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
resource "aws_subnet" "private_subnet_aws-1" {
  tags              = { "Name" = "private-subnet-1" }
  cidr_block        = "10.0.2.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "private_subnet_aws-2" {
  tags              = { "Name" = "private-subnet-2" }
  cidr_block        = "10.0.2.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "private_subnet_aws-3" {
  tags              = { "Name" = "private-subnet-3" }
  cidr_block        = "10.0.2.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "private_subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "private-subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet_route_table_association" "private_subnet_azure" {
  subnet_id      = azurerm_subnet.private_subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "aws_subnet" "public_subnet_aws-1" {
  tags              = { "Name" = "public-subnet-1" }
  cidr_block        = "10.0.3.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "public_subnet_aws-2" {
  tags              = { "Name" = "public-subnet-2" }
  cidr_block        = "10.0.3.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "public_subnet_aws-3" {
  tags              = { "Name" = "public-subnet-3" }
  cidr_block        = "10.0.3.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "public_subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "public-subnet"
  address_prefixes     = ["10.0.3.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
resource "aws_route_table" "rt_aws" {
  tags   = { "Name" = "test-rt" }
  vpc_id = aws_vpc.example_vn_aws.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
  provider = "aws.eu-west-1"
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-rt"
  location            = "northeurope"
  route {
    name           = "internet"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "Internet"
  }
}
resource "aws_route_table_association" "rta_aws-1" {
  subnet_id      = aws_subnet.public_subnet_aws-1.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-2" {
  subnet_id      = aws_subnet.public_subnet_aws-2.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-3" {
  subnet_id      = aws_subnet.public_subnet_aws-3.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "azurerm_subnet_route_table_association" "public_subnet_azure" {
  subnet_id      = azurerm_subnet.public_subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "aws_iam_instance_profile" "vm_aws" {
  name     = "multy-vm-vm_aws-role"
  role     = aws_iam_role.vm_aws.name
  provider = "aws.eu-west-1"
}
resource "aws_iam_role" "vm_aws" {
  tags               = { "Name" = "test-vm" }
  name               = "multy-vm-vm_aws-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
data "aws_ami" "vm_aws" {
  owners      = ["099720109477"]
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu*-18.04-amd64-server-*"]
  }
  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  provider = "aws.eu-west-1"
}
resource "aws_instance" "vm_aws" {
  tags                        = { "Name" = "test-vm" }
  ami                         = data.aws_ami.vm_aws.id
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public_subnet_aws-1.id
  user_data_base64            = "ZWNobyAnSGVsbG8gV29ybGQn"
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
  provider                    = "aws.eu-west-1"
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  allocation_method   = "Static"
  sku                 = "Standard"
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.public_subnet_azure.id
    public_ip_address_id          = azurerm_public_ip.vm_azure.id
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
  zone                  = "1"
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]
  custom_data           = "ZWNobyAnSGVsbG8gV29ybGQn"
  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }
  admin_username = "adminuser"
  admin_password = random_password.vm_azure.result
  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  disable_password_authentication = false
  identity {
    type = "SystemAssigned"
  }
  computer_name = "testvm"
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}
provider "azurerm" {
  features {
  }
}
