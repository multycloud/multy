resource "aws_s3_bucket" "obj_storage1_aws" {
  bucket = "mty-storage-001"
}
resource "aws_s3_bucket" "obj_storage2_aws" {
  bucket = "mty-storage-002"
}
resource "aws_s3_bucket" "obj_storage3_aws" {
  bucket     = "mty-storage-003"
  depends_on = [
    aws_s3_bucket.obj_storage1_aws,
    aws_s3_bucket.obj_storage2_aws,
    azurerm_storage_account.obj_storage1_azure,
    azurerm_storage_container.obj_storage1_azure_public,
    azurerm_storage_container.obj_storage1_azure_private,
    azurerm_storage_account.obj_storage2_azure,
    azurerm_storage_container.obj_storage2_azure_public,
    azurerm_storage_container.obj_storage2_azure_private,
  ]
}
resource "azurerm_storage_account" "obj_storage1_azure" {
  resource_group_name             = azurerm_resource_group.st-rg.name
  name                            = "mtystorage001"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "GZRS"
  allow_nested_items_to_be_public = true
}
resource "azurerm_storage_container" "obj_storage1_azure_public" {
  name                  = "public"
  storage_account_name  = azurerm_storage_account.obj_storage1_azure.name
  container_access_type = "blob"
}
resource "azurerm_storage_container" "obj_storage1_azure_private" {
  name                  = "private"
  storage_account_name  = azurerm_storage_account.obj_storage1_azure.name
  container_access_type = "private"
}
resource "azurerm_storage_account" "obj_storage2_azure" {
  resource_group_name             = azurerm_resource_group.st-rg.name
  name                            = "mtystorage002"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "GZRS"
  allow_nested_items_to_be_public = true
}
resource "azurerm_storage_container" "obj_storage2_azure_public" {
  name                  = "public"
  storage_account_name  = azurerm_storage_account.obj_storage2_azure.name
  container_access_type = "blob"
}
resource "azurerm_storage_container" "obj_storage2_azure_private" {
  name                  = "private"
  storage_account_name  = azurerm_storage_account.obj_storage2_azure.name
  container_access_type = "private"
}
resource "azurerm_storage_account" "obj_storage3_azure" {
  resource_group_name             = azurerm_resource_group.st-rg.name
  name                            = "mtystorage003"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "GZRS"
  allow_nested_items_to_be_public = true
  depends_on                      = [
    aws_s3_bucket.obj_storage1_aws,
    aws_s3_bucket.obj_storage2_aws,
    azurerm_storage_account.obj_storage1_azure,
    azurerm_storage_container.obj_storage1_azure_public,
    azurerm_storage_container.obj_storage1_azure_private,
    azurerm_storage_account.obj_storage2_azure,
    azurerm_storage_container.obj_storage2_azure_public,
    azurerm_storage_container.obj_storage2_azure_private,
  ]
}
resource "azurerm_storage_container" "obj_storage3_azure_public" {
  name                  = "public"
  storage_account_name  = azurerm_storage_account.obj_storage3_azure.name
  container_access_type = "blob"
  depends_on            = [
    aws_s3_bucket.obj_storage1_aws,
    aws_s3_bucket.obj_storage2_aws,
    azurerm_storage_account.obj_storage1_azure,
    azurerm_storage_container.obj_storage1_azure_public,
    azurerm_storage_container.obj_storage1_azure_private,
    azurerm_storage_account.obj_storage2_azure,
    azurerm_storage_container.obj_storage2_azure_public,
    azurerm_storage_container.obj_storage2_azure_private,
  ]
}
resource "azurerm_storage_container" "obj_storage3_azure_private" {
  name                  = "private"
  storage_account_name  = azurerm_storage_account.obj_storage3_azure.name
  container_access_type = "private"
  depends_on            = [
    aws_s3_bucket.obj_storage1_aws,
    aws_s3_bucket.obj_storage2_aws,
    azurerm_storage_account.obj_storage1_azure,
    azurerm_storage_container.obj_storage1_azure_public,
    azurerm_storage_container.obj_storage1_azure_private,
    azurerm_storage_account.obj_storage2_azure,
    azurerm_storage_container.obj_storage2_azure_public,
    azurerm_storage_container.obj_storage2_azure_private,
  ]
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
