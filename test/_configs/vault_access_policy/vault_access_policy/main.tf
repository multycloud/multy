resource "aws_ssm_parameter" "api_key_aws" {
  name     = "/dev-test-secret-multy/api-key"
  type     = "SecureString"
  value    = "xxx"
  provider = "aws.eu-west-1"
}
resource "azurerm_key_vault_secret" "api_key_azure" {
  name         = "api-key"
  key_vault_id = azurerm_key_vault.example_azure.id
  value        = "xxx"
}
resource "google_secret_manager_secret" "api_key_gcp" {
  project   = "multy-project"
  secret_id = "api-key"
  replication {
    user_managed {
      replicas {
        location = "europe-west1"
      }
    }
  }
  provider = "google.europe-west1"
}
resource "google_secret_manager_secret_version" "api_key_gcp" {
  secret      = google_secret_manager_secret.api_key_gcp.id
  secret_data = "xxx"
  provider    = "google.europe-west1"
}
data "azurerm_client_config" "example_azure" {
}
resource "azurerm_key_vault" "example_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "dev-test-secret-multy"
  location            = "northeurope"
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.example_azure.tenant_id
  access_policy {
    tenant_id               = data.azurerm_client_config.example_azure.tenant_id
    object_id               = data.azurerm_client_config.example_azure.object_id
    certificate_permissions = []
    key_permissions         = []
    secret_permissions      = ["List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"]
  }
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
resource "google_compute_network" "example_vn_gcp" {
  name                            = "example-gcp"
  project                         = "multy-project"
  routing_mode                    = "REGIONAL"
  description                     = "Managed by Multy"
  auto_create_subnetworks         = false
  delete_default_routes_on_create = true
  provider                        = "google.europe-west1"
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
data "azurerm_client_config" "vault_access_policy_azure" {
}
resource "azurerm_key_vault_access_policy" "vault_access_policy_azure" {
  key_vault_id            = azurerm_key_vault.example_azure.id
  tenant_id               = data.azurerm_client_config.vault_access_policy_azure.tenant_id
  object_id               = "cb47ad5c-c182-4dad-893d-10b9558e82d0"
  certificate_permissions = []
  key_permissions         = []
  secret_permissions      = ["List", "Get"]
}
resource "google_secret_manager_secret_iam_member" "vault_access_policy_gcp-api_key_gcp" {
  project   = "multy-project"
  secret_id = google_secret_manager_secret.api_key_gcp.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:test@multy-project.iam.gserviceaccount.com"
  provider  = "google.europe-west1"
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
data "aws_caller_identity" "vault_access_policy_aws" {
  provider = "aws.eu-west-1"
}
resource "aws_iam_policy" "vault_access_policy_aws" {
  tags     = { "Name" = "vault_access_policy_aws" }
  name     = "vault_access_policy_aws"
  policy   = "{\"Statement\":[{\"Action\":[\"ssm:GetParameter*\"],\"Effect\":\"Allow\",\"Resource\":\"arn:aws:ssm:eu-west-1:${data.aws_caller_identity.vault_access_policy_aws.account_id}:parameter/dev-test-secret-multy/*\"},{\"Action\":[\"ssm:DescribeParameters\"],\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"2012-10-17\"}"
  provider = "aws.eu-west-1"
}
resource "aws_iam_role_policy_attachment" "vault_access_policy_aws" {
  role       = "multy-vm-vm_aws-role"
  policy_arn = aws_iam_policy.vault_access_policy_aws.arn
  provider   = "aws.eu-west-1"
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
resource "google_service_account" "test-vm-vmgcp-sa-dvl7" {
  project      = "multy-project"
  account_id   = "test-vm-vmgcp-sa-dvl7"
  display_name = "Service Account for VM test-vm"
  provider     = "google.europe-west1"
}
resource "google_compute_instance" "vm_gcp" {
  name         = "test-vm"
  project      = "multy-project"
  machine_type = "e2-micro"
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }
  zone = "europe-west1-b"
  tags = ["subnet-subnet"]
  network_interface {
    subnetwork = google_compute_subnetwork.subnet_gcp.self_link
  }
  metadata = { "user-data" = "echo 'Hello World'" }
  service_account {
    email  = google_service_account.test-vm-vmgcp-sa-dvl7.email
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
