resource "aws_db_subnet_group" "example_db_aws" {
  tags       = { "Name" = "example-db" }
  name       = "example-db"
  subnet_ids = ["${aws_subnet.subnet1_aws.id}", "${aws_subnet.subnet2_aws.id}"]
}
resource "aws_db_instance" "example_db_aws" {
  tags                 = { "Name" = "exampledb" }
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
resource "aws_vpc" "example_vn_aws" {
  tags                 = { "Name" = "example_vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
  vpc_id = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  tags   = { "Name" = "example_vn" }
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
  tags   = { "Name" = "test-nsg2" }
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
    protocol    = "tcp"
    from_port   = 4000
    to_port     = 4000
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
    protocol    = "tcp"
    from_port   = 4000
    to_port     = 4000
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
  tags   = { "Name" = "test-rt" }
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
resource "aws_route_table_association" "rta2_aws" {
  subnet_id      = "${aws_subnet.subnet2_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_route_table_association" "rta3_aws" {
  subnet_id      = "${aws_subnet.subnet1_aws.id}"
  route_table_id = "${aws_route_table.rt_aws.id}"
}
resource "aws_subnet" "subnet1_aws" {
  tags              = { "Name" = "subnet1" }
  cidr_block        = "10.0.1.0/24"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "us-east-1a"
}
resource "aws_subnet" "subnet2_aws" {
  tags              = { "Name" = "subnet2" }
  cidr_block        = "10.0.2.0/24"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "us-east-1b"
}
resource "aws_subnet" "subnet3_aws" {
  tags       = { "Name" = "subnet3" }
  cidr_block = "10.0.3.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
resource "aws_key_pair" "vm_aws" {
  tags       = { "Name" = "test-vm" }
  key_name   = "vm_multy"
  public_key = file("./ssh_key.pub")
}
resource "aws_instance" "vm_aws" {
  tags                        = { "Name" = "test-vm" }
  ami                         = "ami-04ad2567c9e3d7893"
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = "${aws_subnet.subnet3_aws.id}"
  user_data_base64            = base64encode("${templatefile("init.sh", { "db_host" = "${aws_db_instance.example_db_aws.address}", "db_password" = "multy$Admin123!", "db_username" = "${aws_db_instance.example_db_aws.username}" })}")
  key_name                    = aws_key_pair.vm_aws.key_name
}
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "eastus"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "eastus"
  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_resource_group" "nsg-rg" {
  name     = "nsg-rg"
  location = "eastus"
}
resource "azurerm_network_security_group" "nsg2_azure" {
  resource_group_name = azurerm_resource_group.nsg-rg.name
  name                = "test-nsg2"
  location            = "eastus"
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
  security_rule {
    name                       = "6"
    protocol                   = "tcp"
    priority                   = 160
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "4000-4000"
    destination_address_prefix = "*"
    direction                  = "Inbound"
  }
  security_rule {
    name                       = "7"
    protocol                   = "tcp"
    priority                   = 160
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "4000-4000"
    destination_address_prefix = "*"
    direction                  = "Outbound"
  }
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "test-rt"
  location            = "eastus"
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
resource "azurerm_subnet_route_table_association" "rta2_azure" {
  subnet_id      = "${azurerm_subnet.subnet2_azure.id}"
  route_table_id = "${azurerm_route_table.rt_azure.id}"
}
resource "azurerm_subnet_route_table_association" "rta3_azure" {
  subnet_id      = "${azurerm_subnet.subnet1_azure.id}"
  route_table_id = "${azurerm_route_table.rt_azure.id}"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet" "subnet3_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet3"
  address_prefixes     = ["10.0.3.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "eastus"
  allocation_method   = "Static"
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "eastus"
  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = "${azurerm_subnet.subnet3_azure.id}"
    public_ip_address_id          = azurerm_public_ip.vm_azure.id
    primary                       = true
  }
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  resource_group_name   = azurerm_resource_group.vm-rg.name
  name                  = "test-vm"
  location              = "eastus"
  size                  = "Standard_B1ls"
  network_interface_ids = ["${azurerm_network_interface.vm_azure.id}"]
  custom_data           = base64encode("${templatefile("init.sh", { "db_host" = "${aws_db_instance.example_db_aws.address}", "db_password" = "multy$Admin123!", "db_username" = "${aws_db_instance.example_db_aws.username}" })}")
  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }
  admin_username = "adminuser"
  admin_ssh_key {
    username   = "adminuser"
    public_key = file("./ssh_key.pub")
  }
  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7_9-gen2"
    version   = "latest"
  }
  disable_password_authentication = true
}
resource "azurerm_resource_group" "vm-rg" {
  name     = "vm-rg"
  location = "eastus"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "eastus"
}
provider "aws" {
  region = "us-east-1"
}
provider "azurerm" {
  features {
  }
}
output "aws_endpoint" {
  value = "http://${aws_instance.vm_aws.public_ip}:4000"
}
output "azure_endpoint" {
  value = "http://${azurerm_public_ip.vm_azure.ip_address}:4000"
}
