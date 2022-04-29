resource "local_file" "file2_private_azure" {
  content_base64 = "PGgxPkhpPC9oMT4="
  filename       = "./multy/local/file2_private_azure"
}

resource "local_file" "file1_public_azure" {
  content_base64 = "PGgxPkhpPC9oMT4="
  filename       = "./multy/local/file1_public_azure"
}

resource "aws_s3_object" "file1_public_aws" {
  bucket         = aws_s3_bucket.obj_storage_aws.id
  key            = "index.html"
  acl            = "public-read"
  content_base64 = "PGgxPkhpPC9oMT4="
  content_type   = "text/html"
}
resource "aws_s3_object" "file2_private_aws" {
  bucket         = aws_s3_bucket.obj_storage_aws.id
  key            = "index_private.html"
  acl            = "private"
  content_base64 = "PGgxPkhpPC9oMT4="
  content_type   = "text/html"
}
resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "test-storage-9999919"
}
resource "azurerm_storage_blob" "file1_public_azure" {
  name                   = "index.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_public.name
  type                   = "Block"
  source                 = local_file.file1_public_azure.filename
  content_type           = "text/html"
}
resource "azurerm_storage_blob" "file2_private_azure" {
  name                   = "index_private.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_private.name
  type                   = "Block"
  source                 = local_file.file2_private_azure.filename
  content_type           = "text/html"
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name             = azurerm_resource_group.rg1.name
  name                            = "teststorage9999919"
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
resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
