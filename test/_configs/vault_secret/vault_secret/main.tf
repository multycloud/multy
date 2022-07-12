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
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
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
