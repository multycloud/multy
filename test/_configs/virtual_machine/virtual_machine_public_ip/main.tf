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
resource "aws_route_table_association" "rta_aws-1" {
  provider       = "aws.eu-west-1"
  subnet_id      = aws_subnet.subnet_aws-1.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_route_table_association" "rta_aws-2" {
  provider       = "aws.eu-west-1"
  subnet_id      = aws_subnet.subnet_aws-2.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_route_table_association" "rta_aws-3" {
  provider       = "aws.eu-west-1"
  subnet_id      = aws_subnet.subnet_aws-3.id
  route_table_id = aws_route_table.rt_aws.id
}
resource "aws_subnet" "subnet_aws-1" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "subnet-1"
  }

  cidr_block        = "10.0.2.0/25"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1a"
}

resource "aws_subnet" "subnet_aws-2" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "subnet-2"
  }

  cidr_block        = "10.0.2.128/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1b"
}

resource "aws_subnet" "subnet_aws-3" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "subnet-3"
  }

  cidr_block        = "10.0.2.192/26"
  vpc_id            = aws_vpc.example_vn_aws.id
  availability_zone = "eu-west-1c"
}
resource "aws_key_pair" "vm_aws" {
  provider = "aws.eu-west-1"
  tags     = {
    "Name" = "test-vm"
  }

  key_name   = "vm_aws_multy"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP"
}
resource "aws_iam_role" "vm_aws" {
  provider           = "aws.eu-west-1"
  tags               = { "Name" = "test-vm" }
  name               = "multy-vm-vm_aws-role"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_instance_profile" "vm_aws" {
  provider = "aws.eu-west-1"
  name     = "multy-vm-vm_aws-role"
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
  subnet_id                   = aws_subnet.subnet_aws-1.id
  #  user_data_base64            = "#!/bin/bash -xe\\nsudo su; yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl enable httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on AWS</h1>\\\" > /var/www/html/index.html"
  user_data_base64            = "ZWNobyAnSGVsbG8gV29ybGQn"
  key_name                    = aws_key_pair.vm_aws.key_name
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
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
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"

  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
    public_ip_address_id          = azurerm_public_ip.vm_azure.id
    primary                       = true
  }
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "test-vm"
  location            = "northeurope"
  allocation_method   = "Static"
  sku                 = "Standard"
}
resource "azurerm_linux_virtual_machine" "vm_azure" {
  zone                  = "1"
  resource_group_name   = azurerm_resource_group.rg1.name
  name                  = "test-vm"
  computer_name         = "testvm"
  location              = "northeurope"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]
  #  custom_data           = "#!/bin/bash -xe\\nsudo su\\n yum update -y; yum install -y httpd.x86_64; systemctl start httpd.service; systemctl status httpd.service; touch /var/www/html/index.html; echo \\\"<h1>Hello from Multy on Azure</h1>\\\" > /var/www/html/index.html;"
  custom_data           = "ZWNobyAnSGVsbG8gV29ybGQn"

  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }

  admin_username = "adminuser"

  admin_ssh_key {
    username   = "adminuser"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  identity {
    type = "SystemAssigned"
  }
  disable_password_authentication = true
}
data "aws_ami" "vm_aws" {
  provider    = "aws.eu-west-1"
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
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}


provider "azurerm" {
  features {}
}


resource "google_compute_subnetwork" "subnet_gcp" {
  name                     = "subnet"
  ip_cidr_range            = "10.0.2.0/24"
  network                  = google_compute_network.example_vn_gcp.id
  private_ip_google_access = true
  provider                 = "google.europe-west1"
  project                  = "multy-project"
}

resource "google_service_account" "test-vm-vmgcp-sa-dvl7" {
  project      = "multy-project"
  account_id   = "test-vm-vmgcp-sa-dvl7"
  display_name = "Service Account for VM test-vm"
  provider     = "google.europe-west1"
}

resource "google_compute_instance" "vm_gcp" {
  name         = "test-vm"
  machine_type = "e2-standard-2"
  zone         = "europe-west1-c"
  tags         = ["subnet-subnet"]
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }
  network_interface {
    subnetwork = google_compute_subnetwork.subnet_gcp.self_link
    access_config {
      network_tier = "STANDARD"
    }
  }
  metadata = {
    "ssh-keys"  = "adminuser:ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCwSjgjEIKewBWACaOVGg4qsGSHhIeteCmbtn4/DpL0yugLf5c/K/RJVQOKG+dVXVfWD3oAb4JY8jvkZdVACcuocoCewrIEHXxZJmGehxgCeUG8HZ+14mODosUOUYCe3kKCWU2SnUhjX+8x6btxqDOEhtghN3qR52kjm/OUw0Ap43weR1sdkJwtUz7CAXzdCxEKj16R0SY/dNn3uIISPetqm7vqy0ecMJdasbj/X6IAKeiZHe5UKtmOCGYMLwYfKqsrEnzk5rCfa3PK0iYoPv8AB3ocpONcuBshJyDZWhaFBhBrs5SGrWcF34wckD37SNtRZJt+Fuaxe8MpqUVueGTgFViKokCxCfbTnKWRbdGXpSfS6Q0OSvZTWkUEy5ZjxsA03LT4Bcbzq19sABdbyrcEMdv8bq0fhNyGJcGYNJr2uC4+J7irXAM/TuFje4CpJ0G+J3gCrQ2BUeWOYBdjfeP+LckgVXP+TMcEEe4iq5B9psyIS7o58KeNQdFH9jQteIE= joao@Joaos-MBP",
    "user-data" = "echo 'Hello World'"
  }
  provider = "google.europe-west1"
  project  = "multy-project"
  service_account {
    email  = google_service_account.test-vm-vmgcp-sa-dvl7.email
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_network" "example_vn_gcp" {
  name                            = "example_gcp"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
  project                         = "multy-project"
}

resource "google_compute_route" "rt_gcp-0" {
  name             = "test-rt-0"
  dest_range       = "0.0.0.0/0"
  network          = google_compute_network.example_vn_gcp.id
  priority         = 1000
  tags             = ["subnet-subnet"]
  next_hop_gateway = "default-internet-gateway"
  provider         = "google.europe-west1"
  project          = "multy-project"
}

provider "google" {
  region = "europe-west1"
  alias  = "europe-west1"
}
