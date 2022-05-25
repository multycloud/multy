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
  region                      = "eu-west-2"
  alias                       = "eu-west-2"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
}
provider "azurerm" {
  features {}
}
