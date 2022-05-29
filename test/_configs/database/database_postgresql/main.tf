resource "azurerm_postgresql_server" "example_db_azure" {
  resource_group_name              = azurerm_resource_group.rg1.name
  name                             = "example-db"
  location                         = "eastus"
  administrator_login              = "multyadmin"
  administrator_login_password     = "multy$Admin123!"
  sku_name                         = "GP_Gen5_2"
  storage_mb                       = 10240
  version                          = "11"
  ssl_enforcement_enabled          = false
  ssl_minimal_tls_version_enforced = "TLSEnforcementDisabled"
}
resource "azurerm_postgresql_virtual_network_rule" "example_db_azure0" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db0"
  server_name         = azurerm_postgresql_server.example_db_azure.name
  subnet_id           = azurerm_subnet.subnet1_azure.id
}
resource "azurerm_postgresql_virtual_network_rule" "example_db_azure1" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "example-db1"
  server_name         = azurerm_postgresql_server.example_db_azure.name
  subnet_id           = azurerm_subnet.subnet2_azure.id
}
resource "azurerm_postgresql_firewall_rule" "example_db_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "public"
  server_name         = azurerm_postgresql_server.example_db_azure.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
resource "azurerm_subnet" "subnet1_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet1"
  address_prefixes     = ["10.0.0.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_subnet" "subnet2_azure" {
  resource_group_name  = azurerm_resource_group.rg1.name
  name                 = "subnet2"
  address_prefixes     = ["10.0.1.0/24"]
  virtual_network_name = azurerm_virtual_network.vn_azure.name
  service_endpoints    = ["Microsoft.Sql"]
}
resource "azurerm_virtual_network" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"
  address_space       = ["10.0.0.0/16"]
}
resource "azurerm_route_table" "rt_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-rt"
  location            = "eastus"

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
resource "azurerm_subnet_route_table_association" "subnet2_azure" {
  subnet_id      = azurerm_subnet.subnet2_azure.id
  route_table_id = azurerm_route_table.rt_azure.id
}
resource "azurerm_route_table" "vn_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "db-vn"
  location            = "eastus"

  route {
    name           = "local"
    address_prefix = "0.0.0.0/0"
    next_hop_type  = "VnetLocal"
  }
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "eastus"
}
provider "azurerm" {
  features {}
}
