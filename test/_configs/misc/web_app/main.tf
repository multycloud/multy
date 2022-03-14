resource "aws_ssm_parameter" "db_host_aws" {
  name  = "/web-app-vault-test/db-host"
  type  = "SecureString"
  value = "${azurerm_mysql_server.example_db_azure.fqdn}"
}
resource "aws_ssm_parameter" "db_password_aws" {
  name  = "/web-app-vault-test/db-password"
  type  = "SecureString"
  value = "multy-Admin123!"
}
resource "aws_ssm_parameter" "db_username_aws" {
  name  = "/web-app-vault-test/db-username"
  type  = "SecureString"
  value = "multyadmin@example-db"
}
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
  password             = "multy-Admin123!"
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
resource "aws_iam_instance_profile" "vm_aws" {
  depends_on = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  name       = "iam_for_vm_vm"
  role       = aws_iam_role.vm_aws.name
}
resource "aws_iam_role" "vm_aws" {
  depends_on         = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  tags               = { "Name" = "test-vm" }
  name               = "iam_for_vm_vm"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"ec2.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
  inline_policy {
    name   = "vault_policy"
    policy = "{\"Statement\":[{\"Action\":[\"ssm:GetParameter*\"],\"Effect\":\"Allow\",\"Resource\":\"arn:aws:ssm:us-east-1:033721306154:parameter/web-app-vault-test/*\"},{\"Action\":[\"ssm:DescribeParameters\"],\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"2012-10-17\"}"
  }
}
resource "aws_key_pair" "vm_aws" {
  depends_on = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  tags       = { "Name" = "test-vm" }
  key_name   = "vm_multy"
  public_key = file("./ssh_key.pub")
}
resource "aws_instance" "vm_aws" {
  depends_on                  = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  tags                        = { "Name" = "test-vm" }
  ami                         = "ami-04ad2567c9e3d7893"
  instance_type               = "t2.nano"
  associate_public_ip_address = true
  subnet_id                   = "${aws_subnet.subnet3_aws.id}"
  user_data_base64            = base64encode("${templatefile("init.sh", { "db_host" = "${azurerm_mysql_server.example_db_azure.fqdn}", "db_password" = "multy-Admin123!", "db_username" = "multyadmin@example-db" })}")
  key_name                    = aws_key_pair.vm_aws.key_name
  iam_instance_profile        = aws_iam_instance_profile.vm_aws.id
}
resource "azurerm_resource_group" "db-rg" {
  name     = "db-rg"
  location = "eastus"
}
resource "azurerm_key_vault_secret" "db_host_azure" {
  name         = "db-host"
  key_vault_id = azurerm_key_vault.web_app_vault_azure.id
  value        = "${azurerm_mysql_server.example_db_azure.fqdn}"
}
resource "azurerm_key_vault_secret" "db_password_azure" {
  name         = "db-password"
  key_vault_id = azurerm_key_vault.web_app_vault_azure.id
  value        = "multy-Admin123!"
}
resource "azurerm_key_vault_secret" "db_username_azure" {
  name         = "db-username"
  key_vault_id = azurerm_key_vault.web_app_vault_azure.id
  value        = "multyadmin@example-db"
}
resource "azurerm_mysql_server" "example_db_azure" {
  resource_group_name          = azurerm_resource_group.db-rg.name
  name                         = "example-db"
  location                     = "eastus"
  administrator_login          = "multyadmin"
  administrator_login_password = "multy-Admin123!"
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
resource "azurerm_mysql_firewall_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.db-rg.name
  name                = "public"
  server_name         = azurerm_mysql_server.example_db_azure.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
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
resource "azurerm_resource_group" "kv-rg" {
  name     = "kv-rg"
  location = "eastus"
}
data "azurerm_client_config" "kv_ap_azure" {
}
resource "azurerm_key_vault_access_policy" "kv_ap_azure" {
  key_vault_id            = azurerm_key_vault.web_app_vault_azure.id
  tenant_id               = data.azurerm_client_config.kv_ap_azure.tenant_id
  object_id               = "${azurerm_linux_virtual_machine.vm_azure.identity[0].principal_id}"
  certificate_permissions = []
  key_permissions         = []
  secret_permissions      = ["List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"]
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
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.2.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "subnet3_azure" {
  resource_group_name  = azurerm_resource_group.vn-rg.name
  name                 = "subnet3"
  address_prefixes     = ["10.0.3.0/24"]
  virtual_network_name = azurerm_virtual_network.example_vn_azure.name
}
resource "azurerm_public_ip" "vm_azure" {
  depends_on          = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  resource_group_name = azurerm_resource_group.vm-rg.name
  name                = "test-vm"
  location            = "eastus"
  allocation_method   = "Static"
}
resource "azurerm_network_interface" "vm_azure" {
  depends_on          = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
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
  depends_on                      = [
    aws_db_subnet_group.example_db_aws, aws_db_instance.example_db_aws, azurerm_mysql_server.example_db_azure,
    azurerm_mysql_virtual_network_rule.example_db_azure0, azurerm_mysql_virtual_network_rule.example_db_azure1,
    azurerm_mysql_firewall_rule.example_db_azure
  ]
  resource_group_name             = azurerm_resource_group.vm-rg.name
  name                            = "test-vm"
  location                        = "eastus"
  size                            = "Standard_B1ls"
  network_interface_ids           = ["${azurerm_network_interface.vm_azure.id}"]
  custom_data                     = base64encode("${templatefile("azure_init.sh", { "db_host_secret_name" = "db-host", "db_password_secret_name" = "db-password", "db_username_secret_name" = "db-username", "vault_name" = "web-app-vault-test" })}")
  os_disk {
    caching              = "None"
    storage_account_type = "Standard_LRS"
  }
  admin_username                  = "adminuser"
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
  identity {
    type = "SystemAssigned"
  }
}
resource "azurerm_resource_group" "vm-rg" {
  name     = "vm-rg"
  location = "eastus"
}
resource "azurerm_resource_group" "vn-rg" {
  name     = "vn-rg"
  location = "eastus"
}
data "azurerm_client_config" "web_app_vault_azure" {
}
resource "azurerm_key_vault" "web_app_vault_azure" {
  resource_group_name = azurerm_resource_group.kv-rg.name
  name                = "web-app-vault-test"
  location            = "eastus"
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.web_app_vault_azure.tenant_id
  access_policy {
    tenant_id               = data.azurerm_client_config.web_app_vault_azure.tenant_id
    object_id               = data.azurerm_client_config.web_app_vault_azure.object_id
    certificate_permissions = []
    key_permissions         = []
    secret_permissions      = ["List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"]
  }
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
