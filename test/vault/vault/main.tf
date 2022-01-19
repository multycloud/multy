data azurerm_client_config current {}
resource "azurerm_key_vault" "example_azure" {
  resource_group_name = azurerm_resource_group.kv-rg.name
  name                = "dev"
  location            = "ukwest"
  sku_name            = "Standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
resource "azurerm_resource_group" "kv-rg" {
  name     = "kv-rg"
  location = "ukwest"
}
provider "aws" {
  region = "eu-west-2"
}
provider "azurerm" {
  features {}
}
