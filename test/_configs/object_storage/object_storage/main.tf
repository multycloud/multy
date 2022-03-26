resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "test-storage-12384761234"
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name             = azurerm_resource_group.st-rg.name
  name                            = "teststorage12384761234"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "GZRS"
  allow_nested_items_to_be_public = true
  blob_properties {
    versioning_enabled = false
  }
}
resource "azurerm_storage_container" "obj_storage_azure_public" {
  name                  = "public"
  storage_account_name  = azurerm_storage_account.obj_storage_azure.name
  container_access_type = "blob"
}
resource "azurerm_storage_container" "obj_storage_azure_private" {
  name                  = "private"
  storage_account_name  = azurerm_storage_account.obj_storage_azure.name
  container_access_type = "private"
}
resource "azurerm_resource_group" "st-rg" {
  name     = "st-rg"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
