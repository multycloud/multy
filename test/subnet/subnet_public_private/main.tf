resource "aws_db_subnet_group" "example_db_aws" {
  tags = {
    Name = "example-db"
  }

  name = "example-db"

  subnet_ids = [
    "${aws_subnet.subnet1_aws.id}",
    "${aws_subnet.subnet2_aws.id}",
  ]
}
resource "aws_db_instance" "example_db_aws" {
  tags                 = {
    Name = "exampledb"
  }
  db_name              = "exampledb"
  allocated_storage    = 10
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
resource "aws_security_group" "nsg2_aws" {
  tags = {
    Name = "test-nsg2"
  }

  vpc_id = "${aws_vpc.example_vn_aws.id}"

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
}
resource "aws_route_table" "rt_aws" {
  tags = {
    Name = "test-rt"
  }

  vpc_id = "${aws_vpc.example_vn_aws.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.example_vn_aws.id
  }
}
resource "aws_route_table_association" "rta_aws" {
  subnet_id      = "${aws_subnet.subnet3_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_subnet" "subnet1_aws" {
  tags = {
    Name = "subnet1"
  }

  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_subnet" "subnet2_aws" {
  tags = {
    Name = "subnet2"
  }

  cidr_block        = "10.0.2.0/24"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
}
resource "aws_subnet" "subnet3_aws" {
  tags = {
    Name = "subnet3"
  }

  cidr_block = "10.0.3.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_instance" "vm_aws" {
  tags = {
    Name = "test-vm"
  }

  ami                         = "ami-09d4a659cdd8677be"
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = "${aws_subnet.subnet3_aws.id}"
  user_data_base64            = "IyEvYmluL2Jhc2ggLXhlCnN1ZG8gc3U7IHl1bSB1cGRhdGUgLXk7IHl1bSBpbnN0YWxsIC15IGh0dHBkLng4Nl82NDsgc3lzdGVtY3RsIHN0YXJ0IGh0dHBkLnNlcnZpY2U7IHN5c3RlbWN0bCBlbmFibGUgaHR0cGQuc2VydmljZTsgdG91Y2ggL3Zhci93d3cvaHRtbC9pbmRleC5odG1sOyBlY2hvICI8aDE+SGVsbG8gZnJvbSBNdWx0eSBvbiBBV1M8L2gxPiIgPiAvdmFyL3d3dy9odG1sL2luZGV4Lmh0bWw="
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
resource "azurerm_resource_group" "nsg-rg" {
  name     = "nsg-rg"
  location = "northeurope"
}
resource "azurerm_network_security_group" "nsg2_azure" {
  resource_group_name = azurerm_resource_group.nsg-rg.name
  name                = "test-nsg2"
  location            = "northeurope"

  security_rule {
    name                       = "0"
    protocol                   = "tcp"
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
    protocol                   = "tcp"
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
    protocol                   = "tcp"
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
    protocol                   = "tcp"
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
    protocol                   = "tcp"
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
    protocol                   = "tcp"
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
  subnet_id      = "${azurerm_subnet.subnet3_azure.id}"
  route_table_id = "${azurerm_route_table.rt_azure.id}"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet_route_table_association" "subnet1_azure" {
  subnet_id      = "${azurerm_subnet.subnet1_azure.id}"
  route_table_id = "${azurerm_route_table.example_vn_azure.id}"
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet_route_table_association" "subnet2_azure" {
  subnet_id      = "${azurerm_subnet.subnet2_azure.id}"
  route_table_id = "${azurerm_route_table.example_vn_azure.id}"
}
resource "azurerm_subnet" "subnet3_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet3"
  address_prefixes     = ["10.0.3.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.subnet3_azure.id}"
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
  size                  = "Standard_B1ls"
  network_interface_ids = ["${azurerm_network_interface.vm_azure.id}"]
  custom_data           = "IyEvYmluL2Jhc2ggLXhlCnN1ZG8gc3UKIHl1bSB1cGRhdGUgLXk7IHl1bSBpbnN0YWxsIC15IGh0dHBkLng4Nl82NDsgc3lzdGVtY3RsIHN0YXJ0IGh0dHBkLnNlcnZpY2U7IHN5c3RlbWN0bCBzdGF0dXMgaHR0cGQuc2VydmljZTsgdG91Y2ggL3Zhci93d3cvaHRtbC9pbmRleC5odG1sOyBlY2hvICI8aDE+SGVsbG8gZnJvbSBNdWx0eSBvbiBBenVyZTwvaDE+IiA+IC92YXIvd3d3L2h0bWwvaW5kZXguaHRtbA=="

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
