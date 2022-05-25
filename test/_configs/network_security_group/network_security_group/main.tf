resource "aws_vpc" "example_vn_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "example_vn"
  }

  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "example_vn"
  }

  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  provider = "aws.eu-west-1"
  tags     = {
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
resource "aws_security_group" "nsg2_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "test-nsg2"
  }

  name        = "test-nsg2"
  description = "Managed by Multy"
  vpc_id      = aws_vpc.example_vn_aws.id

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
}
resource "aws_route_table" "rt_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "test-rt"
  }

  vpc_id = aws_vpc.example_vn_aws.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "subnet1_aws" {
  provider       = "aws.eu-west-1"
  subnet_id      = aws_subnet.subnet1_aws.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_subnet" "subnet1_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "subnet1"
  }

  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_iam_role" "vm_aws" {
  provider           = "aws.eu-west-1"
  tags               = { "Name" = "test-vm" }
  name               = "iam_for_vm_vm_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_instance_profile" "vm_aws" {
  provider = "aws.eu-west-1"
  name     = "iam_for_vm_vm_aws"
  role     = aws_iam_role.vm_aws.name
}
resource "aws_instance" "vm_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "test-vm"
  }

  ami                         = data.aws_ami.vm_aws.id
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.subnet1_aws.id
  user_data_base64            = "ZWNobyAnSGVsbG8gV29ybGQn"
  #  user_data_base64            = "#!/bin/bash -xe\\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on AWS</h1>\\\" > /var/www/html/index.html"
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
}
resource "aws_iam_role" "vm2_aws" {
  provider           = "aws.eu-west-1"
  tags               = { "Name" = "test-vm2" }
  name               = "iam_for_vm_vm2_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_instance_profile" "vm2_aws" {
  provider = "aws.eu-west-1"
  name     = "iam_for_vm_vm2_aws"
  role     = aws_iam_role.vm2_aws.name
}
resource "aws_instance" "vm2_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "test-vm2"
  }

  ami                         = data.aws_ami.vm2_aws.id
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.subnet1_aws.id
  #  user_data_base64            = "#!/bin/bash -xe\\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on AWS</h1>\\\" > /var/www/html/index.html"
  user_data_base64            = "ZWNobyAnSGVsbG8gV29ybGQn"
  vpc_security_group_ids      = [aws_security_group.nsg2_aws.id]
  iam_instance_profile        = aws_iam_instance_profile.vm2_aws.id
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
resource "azurerm_network_security_group" "nsg2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-nsg2"
  location            = "northeurope"

  security_rule {
    name                       = "0"
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
    name                       = "1"
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
    name                       = "2"
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
    name                       = "3"
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
resource "azurerm_subnet_route_table_association" "subnet1_azure" {
  subnet_id      = azurerm_subnet.subnet1_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet1_azure.id
    public_ip_address_id          = azurerm_public_ip.vm_azure.id
    primary                       = true
  }
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  allocation_method   = "Static"
}
resource "random_password" "vm_azure" {
  length  = 16
  special = true
  upper   = true
  lower   = true
  number  = true
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm"
  computer_name         = "testvm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]
  #  custom_data           = "#!/bin/bash -xe\\nsudo su\\n yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on Azure</h1>\\\" > /var/www/html/index.html"
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
    sku       = "16.04-LTS"
    version   = "latest"
  }
  identity {
    type = "SystemAssigned"
  }
  disable_password_authentication = false
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
resource "azurerm_network_interface" "vm2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm2"
  location            = "northeurope"

  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet1_azure.id
    public_ip_address_id          = azurerm_public_ip.vm2_azure.id
    primary                       = true
  }
}
resource "azurerm_public_ip" "vm2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm2"
  location            = "northeurope"
  allocation_method   = "Static"
}
resource "azurerm_network_interface_security_group_association" "vm2_azure" {
  network_interface_id      = azurerm_network_interface.vm2_azure.id
  network_security_group_id = azurerm_network_security_group.nsg2_azure.id
}
resource "random_password" "vm2_azure" {
  length  = 16
  special = true
  upper   = true
  lower   = true
  number  = true
}
resource "azurerm_linux_virtual_machine" "vm2_azure" {
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm2"
  computer_name         = "testvm2"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm2_azure.id]
  #  custom_data           = "#!/bin/bash -xe\\nsudo su\\n yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on Azure</h1>\\\" > /var/www/html/index.html"
  custom_data           = "ZWNobyAnSGVsbG8gV29ybGQn"

  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }

  admin_username = "adminuser"
  admin_password = random_password.vm2_azure.result

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
  identity {
    type = "SystemAssigned"
  }
  disable_password_authentication = false
}
data "aws_ami" "vm_aws" {
  provider    = "aws.eu-west-1"
  owners      = ["099720109477"]
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu*-16.04-amd64-server-*"]
  }
  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
data "aws_ami" "vm2_aws" {
  provider    = "aws.eu-west-1"
  owners      = ["099720109477"]
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu*-16.04-amd64-server-*"]
  }
  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
provider "aws" {
  region                      = "eu-west-1"
  alias                       = "eu-west-1"
  skip_credentials_validation = true
}


provider "azurerm" {
  features {}
}
