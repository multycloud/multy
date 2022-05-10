resource "aws_ssm_parameter" "api_key_aws" {
  provider = "aws.eu-west-2"
  name     = "/dev-test-secret-multy/api-key"
  type     = "SecureString"
  value    = "xxx"
}
resource "aws_vpc" "example_vn_aws" {
  provider             = "aws.eu-west-2"
  tags                 = { "Name" = "example_vn" }
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}
resource "aws_internet_gateway" "example_vn_aws" {
  provider = "aws.eu-west-2"
  tags     = { "Name" = "example_vn" }
  vpc_id   = aws_vpc.example_vn_aws.id
}
resource "aws_default_security_group" "example_vn_aws" {
  provider = "aws.eu-west-2"
  tags     = { "Name" = "example_vn" }
  vpc_id   = aws_vpc.example_vn_aws.id
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
resource "aws_subnet" "subnet_aws" {
  provider   = "aws.eu-west-2"
  tags       = { "Name" = "subnet" }
  cidr_block = "10.0.2.0/24"
  vpc_id     = aws_vpc.example_vn_aws.id
}
data "aws_caller_identity" "vm_aws" {
  provider = "aws.eu-west-2"
}
data "aws_region" "vm_aws" {
  provider = "aws.eu-west-2"
}
resource "aws_iam_instance_profile" "vm_aws" {
  provider = "aws.eu-west-2"
  name     = "iam_for_vm_vm_aws"
  role     = aws_iam_role.vm_aws.name
}
resource "aws_iam_role" "vm_aws" {
  provider           = "aws.eu-west-2"
  tags               = { "Name" = "test-vm" }
  name               = "iam_for_vm_vm_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  inline_policy {
    name   = "vault_policy"
    policy = "{\"Statement\":[{\"Action\":[\"ssm:GetParameter*\"],\"Effect\":\"Allow\",\"Resource\":\"arn:aws:ssm:${data.aws_region.vm_aws.name}:${data.aws_caller_identity.vm_aws.account_id}:parameter/dev-test-secret-multy/*\"},{\"Action\":[\"ssm:DescribeParameters\"],\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"2012-10-17\"}"
  }
}
resource "aws_instance" "vm_aws" {
  provider                    = "aws.eu-west-2"
  tags                        = { "Name" = "test-vm" }
  ami                         = ""
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.subnet_aws.id
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
}
resource "azurerm_key_vault_secret" "api_key_azure" {
  name         = "api-key"
  key_vault_id = azurerm_key_vault.example_azure.id
  value        = "xxx"
}
data "azurerm_client_config" "example_azure" {
}
resource "azurerm_key_vault" "example_azure" {
  resource_group_name = azurerm_resource_group.kv-rg.name
  name                = "dev-test-secret-multy"
  location            = "uksouth"
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
resource "azurerm_virtual_network" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "uksouth"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "example_vn_azure" {
  resource_group_name = azurerm_resource_group.vn-rg.name
  name                = "example_vn"
  location            = "uksouth"
  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_resource_group" "kv-rg" {
  name     = "kv-rg"
  location = "uksouth"
}
data "azurerm_client_config" "kv_ap_azure" {
}
resource "azurerm_key_vault_access_policy" "kv_ap_azure" {
  key_vault_id            = azurerm_key_vault.example_azure.id
  tenant_id               = data.azurerm_client_config.kv_ap_azure.tenant_id
  object_id               = azurerm_linux_virtual_machine.vm_azure.identity[0].principal_id
  certificate_permissions = []
  key_permissions         = []
  secret_permissions      = ["List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"]
}
resource "azurerm_subnet" "subnet_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_subnet_route_table_association" "subnet_azure" {
  subnet_id      = azurerm_subnet.subnet_azure.id
  route_table_id = azurerm_route_table.example_vn_azure.id
}
resource "azurerm_public_ip" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "uksouth"
  allocation_method   = "Static"
}
resource "azurerm_network_interface" "vm_azure" {
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "uksouth"
  ip_configuration {
    name                          = "external"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.subnet_azure.id
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
  resource_group_name   = azurerm_resource_group.vm-rg.name
  name                  = "test-vm"
  computer_name         = "testvm"
  location              = "uksouth"
  size                  = "Standard_B1ls"
  network_interface_ids = [azurerm_network_interface.vm_azure.id]
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
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_resource_group" "vm-rg" {
  name     = "vm-rg"
  location = "uksouth"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "uksouth"
}
provider "aws" {
  region = "eu-west-2"
  alias  = "eu-west-2"
}


provider "azurerm" {
  features {
  }
}
