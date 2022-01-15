resource "aws_s3_bucket_object" "file1_public_aws" {
  bucket       = aws_s3_bucket.obj_storage_aws.id
  key          = "index.html"
  acl          = "public-read"
  content      = "<h1>Hi from AWS</h1>"
  content_type = "text/html"
}
resource "aws_s3_bucket_object" "file2_private_aws" {
  bucket       = aws_s3_bucket.obj_storage_aws.id
  key          = "index_private.html"
  content      = "<h1>Hi</h1>"
  content_type = "text/html"
}
resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "test-storage-9999919"
}
resource "azurerm_storage_blob" "file1_public_azure" {
  name                   = "index.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_public.name
  type                   = "Block"
  source_content         = "<h1>Hi from Azure</h1>"
  content_type           = "text/html"
}
resource "azurerm_storage_blob" "file2_private_azure" {
  name                   = "index_private.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_private.name
  type                   = "Block"
  source_content         = "<h1>Hi</h1>"
  content_type           = "text/html"
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name      = azurerm_resource_group.st-rg.name
  name                     = "teststorage9999919"
  location                 = "northeurope"
  account_tier             = "Standard"
  account_replication_type = "GZRS"
  allow_blob_public_access = true
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
