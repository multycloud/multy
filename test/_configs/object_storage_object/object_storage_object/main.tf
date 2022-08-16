resource "aws_s3_object" "file1_public_aws" {
  provider       = "aws.eu-west-1"
  bucket         = aws_s3_bucket.obj_storage_aws.id
  key            = "index.html"
  acl            = "public-read"
  content_base64 = "PGgxPkhpPC9oMT4="
  content_type   = "text/html"
}
resource "aws_s3_object" "file2_private_aws" {
  provider       = "aws.eu-west-1"
  bucket         = aws_s3_bucket.obj_storage_aws.id
  key            = "index_private.html"
  acl            = "private"
  content_base64 = "PGgxPkhpPC9oMT4="
  content_type   = "text/html"
}
resource "aws_s3_bucket" "obj_storage_aws" {
  provider = "aws.eu-west-1"
  bucket   = "teststorage9999919"
}
resource "azurerm_storage_blob" "file1_public_azure" {
  name                   = "index.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_public.name
  type                   = "Block"
  source_content         = base64decode("PGgxPkhpPC9oMT4=")
  content_type           = "text/html"
}
resource "azurerm_storage_blob" "file2_private_azure" {
  name                   = "index_private.html"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_private.name
  type                   = "Block"
  source_content         = base64decode("PGgxPkhpPC9oMT4=")
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
  alias  = "eu-west-1"
}


provider "azurerm" {
  features {}
}
resource "google_storage_bucket_object" "file1_public_GCP" {
  name         = "index.html"
  bucket       = google_storage_bucket.obj_storage_GCP.name
  content      = base64decode("PGgxPkhpPC9oMT4=")
  content_type = "text/html"
  provider     = "google.europe-west1"
}
resource "google_storage_object_access_control" "file1_public_GCP" {
  object   = google_storage_bucket_object.file1_public_GCP.output_name
  bucket   = google_storage_bucket.obj_storage_GCP.name
  role     = "READER"
  entity   = "allUsers"
  provider = "google.europe-west1"
}
resource "google_storage_bucket_object" "file2_private_GCP" {
  name         = "index_private.html"
  bucket       = google_storage_bucket.obj_storage_GCP.name
  content      = base64decode("PGgxPkhpPC9oMT4=")
  content_type = "text/html"
  provider     = "google.europe-west1"
}
resource "google_storage_bucket" "obj_storage_GCP" {
  name                        = "teststorage9999919"
  project                     = "multy-project"
  uniform_bucket_level_access = false
  location                    = "europe-west1"
  provider                    = "google.europe-west1"
  force_destroy               = true
}
provider "google" {
  region = "europe-west1"
  alias  = "europe-west1"
}