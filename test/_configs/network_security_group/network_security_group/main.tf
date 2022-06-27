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
resource "google_compute_network" "example_vn_gcp" {
  name                            = "example-vn"
  project                         = "multy-project"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
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
    protocol    = "tcp"
    from_port   = 8000
    to_port     = 8000
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
    protocol    = "tcp"
    from_port   = 8001
    to_port     = 8002
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
  security_rule {
    name                       = "4"
    protocol                   = "Tcp"
    priority                   = 150
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "8000-8000"
    destination_address_prefix = "*"
    direction                  = "Inbound"
  }
  security_rule {
    name                       = "5"
    protocol                   = "Tcp"
    priority                   = 160
    access                     = "Allow"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_port_range     = "8001-8002"
    destination_address_prefix = "*"
    direction                  = "Outbound"
  }
}
resource "google_compute_firewall" "nsg_gcp-default-deny-egress" {
  name               = "test-nsg-default-deny-egress"
  project            = "multy-project"
  network            = google_compute_network.example_vn_gcp.id
  direction          = "EGRESS"
  destination_ranges = ["0.0.0.0/0"]
  priority           = 65535
  deny {
    protocol = "all"
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-i-0" {
  name          = "test-nsg-i-0"
  project       = "multy-project"
  network       = google_compute_network.example_vn_gcp.id
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  priority      = 120
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-e-0" {
  name               = "test-nsg-e-0"
  project            = "multy-project"
  network            = google_compute_network.example_vn_gcp.id
  direction          = "EGRESS"
  destination_ranges = ["0.0.0.0/0"]
  priority           = 120
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-i-1" {
  name          = "test-nsg-i-1"
  project       = "multy-project"
  network       = google_compute_network.example_vn_gcp.id
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  priority      = 140
  allow {
    protocol = "tcp"
    ports    = ["443"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-e-1" {
  name               = "test-nsg-e-1"
  project            = "multy-project"
  network            = google_compute_network.example_vn_gcp.id
  direction          = "EGRESS"
  destination_ranges = ["0.0.0.0/0"]
  priority           = 140
  allow {
    protocol = "tcp"
    ports    = ["443"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-i-2" {
  name          = "test-nsg-i-2"
  project       = "multy-project"
  network       = google_compute_network.example_vn_gcp.id
  direction     = "INGRESS"
  source_ranges = ["0.0.0.0/0"]
  priority      = 150
  allow {
    protocol = "tcp"
    ports    = ["8000"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
}
resource "google_compute_firewall" "nsg_gcp-e-3" {
  name               = "test-nsg-e-3"
  project            = "multy-project"
  network            = google_compute_network.example_vn_gcp.id
  direction          = "EGRESS"
  destination_ranges = ["0.0.0.0/0"]
  priority           = 160
  allow {
    protocol = "tcp"
    ports    = ["8001-8002"]
  }
  target_tags = ["nsg-test-nsg"]
  provider    = "google.europe-west1"
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
resource "google_compute_route" "rt_gcp-0" {
  name             = "test-rt-0"
  project          = "multy-project"
  dest_range       = "0.0.0.0/0"
  network          = google_compute_network.example_vn_gcp.id
  priority         = 1000
  tags             = ["subnet-subnet1"]
  next_hop_gateway = "default-internet-gateway"
  provider         = "google.europe-west1"
}
resource "aws_route_table_association" "rta_aws-1" {
  subnet_id      = aws_subnet.subnet1_aws-1.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-2" {
  subnet_id      = aws_subnet.subnet1_aws-2.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "aws_route_table_association" "rta_aws-3" {
  subnet_id      = aws_subnet.subnet1_aws-3.id
  route_table_id = aws_route_table.rt_aws.id
  provider       = "aws.eu-west-1"
}
resource "azurerm_subnet_route_table_association" "subnet1_azure" {
  subnet_id      = azurerm_subnet.subnet1_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "aws_subnet" "subnet1_aws-1" {
  tags              = { "Name" = "subnet1-1" }
  cidr_block        = "10.0.1.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet1_aws-2" {
  tags              = { "Name" = "subnet1-2" }
  cidr_block        = "10.0.1.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
  provider          = "aws.eu-west-1"
}
resource "aws_subnet" "subnet1_aws-3" {
  tags              = { "Name" = "subnet1-3" }
  cidr_block        = "10.0.1.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
  provider          = "aws.eu-west-1"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "google_compute_subnetwork" "subnet1_gcp" {
  name                     = "subnet1"
  project                  = "multy-project"
  ip_cidr_range            = "10.0.1.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
}
resource "aws_iam_instance_profile" "vm2_aws" {
  name     = "multy-vm-vm2_aws-role"
  role     = aws_iam_role.vm2_aws.name
  provider = "aws.eu-west-1"
}
resource "aws_iam_role" "vm2_aws" {
  tags               = { "Name" = "test-vm2" }
  name               = "multy-vm-vm2_aws-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
data "aws_ami" "vm2_aws" {
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
resource "aws_instance" "vm2_aws" {
  tags                        = { "Name" = "test-vm2" }
  ami                         = data.aws_ami.vm2_aws.id
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.subnet1_aws-2.id
  user_data_base64            = "ZWNobyAnSGVsbG8gV29ybGQn"
  vpc_security_group_ids      = [aws_security_group.nsg2_aws.id]
  iam_instance_profile        = aws_iam_instance_profile.vm2_aws.id
  provider                    = "aws.eu-west-1"
}
resource "azurerm_public_ip" "vm2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm2"
  location            = "northeurope"
  allocation_method   = "Static"
  sku                 = "Standard"
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
  zone                  = "2"
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm2"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm2_azure.id]
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
    sku       = "18.04-LTS"
    version   = "latest"
  }
  disable_password_authentication = false
  identity {
    type = "SystemAssigned"
  }
  computer_name = "testvm2"
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
  subnet_id                   = aws_subnet.subnet1_aws-1.id
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
    subnet_id                     = azurerm_subnet.subnet1_azure.id
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
provider "google" {
  region = "europe-west1"
  alias  = "europe-west1"
}
