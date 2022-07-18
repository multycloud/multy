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
  name                            = "example-gcp"
  project                         = "multy-project"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
}
resource "google_compute_firewall" "example_vn_gcp" {
  name               = "example-gcp-default-deny-egress"
  project            = "multy-project"
  network            = google_compute_network.example_vn_gcp.id
  direction          = "EGRESS"
  destination_ranges = ["0.0.0.0/0"]
  priority           = 65535
  deny {
    protocol = "all"
  }
  target_tags = ["vn-example-gcp"]
  provider    = "google.europe-west1"
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
resource "google_compute_subnetwork" "subnet_gcp" {
  name                     = "subnet"
  project                  = "multy-project"
  ip_cidr_range            = "10.0.2.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
}
resource "aws_iam_instance_profile" "vm2_aws" {
  name     = "vm2_aws-vm-role"
  role     = aws_iam_role.vm2_aws.name
  provider = "aws.eu-west-1"
}
resource "aws_iam_role" "vm2_aws" {
  tags               = { "Name" = "test-vm" }
  name               = "vm2_aws-vm-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
data "aws_ami" "vm2_aws" {
  owners      = ["136693071363"]
  most_recent = true
  filter {
    name   = "name"
    values = ["debian-10-amd64-*"]
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
  tags                 = { "Name" = "test-vm" }
  ami                  = data.aws_ami.vm2_aws.id
  instance_type        = "t2.nano"
  subnet_id            = aws_subnet.subnet_aws-3.id
  iam_instance_profile = aws_iam_instance_profile.vm2_aws.id
  provider             = "aws.eu-west-1"
}
resource "azurerm_network_interface" "vm2_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    primary                       = true
  }
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
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm2_azure.id]
  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }
  admin_username = "adminuser"
  admin_password = random_password.vm2_azure.result
  source_image_reference {
    publisher = "Debian"
    offer     = "debian-10"
    sku       = "10"
    version   = "latest"
  }
  disable_password_authentication = false
  identity {
    type = "SystemAssigned"
  }
  computer_name = "testvm"
  zone          = "3"
}
resource "google_service_account" "vm2_gcp" {
  project      = "multy-project"
  account_id   = "test-vm-vm2gcp-sa-w48h"
  display_name = "Service Account for VM test-vm"
  provider     = "google.europe-west1"
}
resource "google_compute_instance" "vm2_gcp" {
  name         = "test-vm"
  project      = "multy-project"
  machine_type = "e2-micro"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
  zone = "europe-west1-c"
  tags = ["vn-example-gcp", "subnet-subnet"]
  network_interface {
    subnetwork = google_compute_subnetwork.subnet_gcp.self_link
  }
  metadata = {}
  service_account {
    email  = google_service_account.vm2_gcp.email
    scopes = ["cloud-platform"]
  }
  provider = "google.europe-west1"
}
resource "aws_iam_instance_profile" "vm3_aws" {
  name     = "vm3_aws-vm-role"
  role     = aws_iam_role.vm3_aws.name
  provider = "aws.eu-west-1"
}
resource "aws_iam_role" "vm3_aws" {
  tags               = { "Name" = "test-vm" }
  name               = "vm3_aws-vm-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  provider           = "aws.eu-west-1"
}
data "aws_ami" "vm3_aws" {
  owners      = ["125523088429"]
  most_recent = true
  filter {
    name   = "name"
    values = ["CentOS 8.2* x86_64"]
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
resource "aws_instance" "vm3_aws" {
  tags                 = { "Name" = "test-vm" }
  ami                  = data.aws_ami.vm3_aws.id
  instance_type        = "t2.nano"
  subnet_id            = aws_subnet.subnet_aws-2.id
  iam_instance_profile = aws_iam_instance_profile.vm3_aws.id
  provider             = "aws.eu-west-1"
}
resource "azurerm_network_interface" "vm3_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    primary                       = true
  }
}
resource "random_password" "vm3_azure" {
  length  = 16
  special = true
  upper   = true
  lower   = true
  number  = true
}
resource "azurerm_linux_virtual_machine" "vm3_azure" {
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm3_azure.id]
  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }
  admin_username = "adminuser"
  admin_password = random_password.vm3_azure.result
  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOs"
    sku       = "8_2"
    version   = "latest"
  }
  disable_password_authentication = false
  identity {
    type = "SystemAssigned"
  }
  computer_name = "testvm"
  zone          = "2"
}
resource "aws_iam_instance_profile" "vm_aws" {
  name     = "vm_aws-vm-role"
  role     = aws_iam_role.vm_aws.name
  provider = "aws.eu-west-1"
}
resource "aws_iam_role" "vm_aws" {
  tags               = { "Name" = "test-vm" }
  name               = "vm_aws-vm-role"
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
  tags                 = { "Name" = "test-vm" }
  ami                  = data.aws_ami.vm_aws.id
  instance_type        = "t2.nano"
  subnet_id            = aws_subnet.subnet_aws-1.id
  user_data_base64     = "ZWNobyAnSGVsbG8gV29ybGQn"
  iam_instance_profile = aws_iam_instance_profile.vm_aws.id
  provider             = "aws.eu-west-1"
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  ip_configuration {
    name                          = "internal"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
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
  zone          = "1"
}
resource "google_service_account" "vm_gcp" {
  project      = "multy-project"
  account_id   = "test-vm-vmgcp-sa-dvl7"
  display_name = "Service Account for VM test-vm"
  provider     = "google.europe-west1"
}
resource "google_compute_instance" "vm_gcp" {
  name         = "test-vm"
  project      = "multy-project"
  machine_type = "e2-medium"
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }
  zone = "europe-west1-c"
  tags = ["vn-example-gcp", "subnet-subnet"]
  network_interface {
    subnetwork = google_compute_subnetwork.subnet_gcp.self_link
  }
  metadata = { "user-data" = "echo 'Hello World'" }
  service_account {
    email  = google_service_account.vm_gcp.email
    scopes = ["cloud-platform"]
  }
  provider = "google.europe-west1"
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
