resource "aws_s3_bucket" "obj_storage_aws" {
  provider = "aws.eu-west-1"
  bucket   = "teststorage12384761234"
}
resource "aws_s3_bucket_versioning" "obj_storage_aws" {
  provider = "aws.eu-west-1"
  bucket   = aws_s3_bucket.obj_storage_aws.id
  versioning_configuration {
    status = "Enabled"
  }
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name             = azurerm_resource_group.rg1.name
  name                            = "teststorage12384761234"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "GZRS"
  allow_nested_items_to_be_public = true
  blob_properties {
    versioning_enabled = true
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
resource "google_storage_bucket" "object_storage_gcp" {
  name                        = "teststorage12384761234"
  project                     = "multy-project"
  uniform_bucket_level_access = false
  versioning {
    enabled = true
  }
  location      = "europe-west1"
  provider      = "google.europe-west1"
  force_destroy = true
}
provider "google" {
  region = "europe-west1"
  alias  = "europe-west1"
}
