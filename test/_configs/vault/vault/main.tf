data azurerm_client_config example_azure {}
resource "azurerm_key_vault" "example_azure" {
  resource_group_name = azurerm_resource_group.rg1.name
  name                = "dev-test-multy"
  location            = "uksouth"
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.example_azure.tenant_id

  access_policy {
    tenant_id = data.azurerm_client_config.example_azure.tenant_id
    object_id = data.azurerm_client_config.example_azure.object_id

    certificate_permissions = []
    secret_permissions      = [
      "List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"
    ]
    key_permissions = []
  }
}
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "uksouth"
}
provider "aws" {
  region = "us-west-1"
  alias  = "us-west-1"
}
provider "azurerm" {
  features {}
}
