resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "test-storage-9999919"
}
resource "aws_s3_bucket_object" "source_code_aws" {
  bucket       = aws_s3_bucket.obj_storage_aws.id
  key          = "source_code.zip"
  acl          = "private"
  content      = "test"
  content_type = "application/zip"
}
resource "aws_lambda_function" "test_aws" {
  tags =  {
    Name = "test_name"
  }

  function_name = "test_name"
  role          = aws_iam_role.iam_for_lambda_test_name.arn
  runtime       = "python3.9"
  handler       = "lambda_function.lambda_handler"
  s3_bucket     = aws_s3_bucket.obj_storage_aws.id
  s3_key        = aws_s3_bucket_object.source_code_aws.key
}
resource "aws_iam_role" "iam_for_lambda_test_name" {
  tags =  {
    Name = "iam_for_lambda_test_name"
  }

  name               = "iam_for_lambda_test_name"
  assume_role_policy = "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}"
}
resource "azurerm_resource_group" "fun-rg" {
  name     = "fun-rg"
  location = "northeurope"
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
resource "azurerm_storage_blob" "source_code_azure" {
  name                   = "source_code.zip"
  storage_account_name   = azurerm_storage_account.obj_storage_azure.name
  storage_container_name = azurerm_storage_container.obj_storage_azure_private.name
  type                   = "Block"
  source_content         = "test"
  content_type           = "application/zip"
}
resource "azurerm_resource_group" "st-rg" {
  name     = "st-rg"
  location = "northeurope"
}
resource "azurerm_function_app" "test" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "testname"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.obj_storage_azure.name
  storage_account_access_key = azurerm_storage_account.obj_storage_azure.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.test.id
  os_type                    = "linux"
}
resource "azurerm_app_service_plan" "test" {
  resource_group_name = azurerm_resource_group.fun-rg.name
  name                = "testservplan"
  location            = "northeurope"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
