resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "test-storage-9999919"
}
resource "aws_s3_bucket_object" "source_code_aws" {
  bucket = aws_s3_bucket.obj_storage_aws.id
  key    = "source_code.zip"
  acl    = "public-read"
  source = "source_dir/aws_code.zip"
}
resource "aws_lambda_function" "test_aws" {
  tags =  {
    Name = "multyfunobj"
  }

  function_name = "multyfunobj"
  role          = aws_iam_role.iam_for_lambda_test.arn
  runtime       = "python3.9"
  handler       = "lambda_function.lambda_handler"
  s3_bucket     = aws_s3_bucket.obj_storage_aws.id
  s3_key        = aws_s3_bucket_object.source_code_aws.key
}
resource "aws_iam_role" "iam_for_lambda_test" {
  tags =  {
    Name = "iam_for_lambda_test"
  }

  name               = "iam_for_lambda_test"
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
  storage_container_name = azurerm_storage_container.obj_storage_azure_public.name
  type                   = "Block"
  source                 = "source_dir/azure_code.zip"
}
resource "azurerm_resource_group" "st-rg" {
  name     = "st-rg"
  location = "northeurope"
}
resource "azurerm_function_app" "test_azure" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "multyfunobj"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.obj_storage_azure.name
  storage_account_access_key = azurerm_storage_account.obj_storage_azure.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.test_azure.id
  os_type                    = "linux"

  app_settings =  {
    WEBSITE_RUN_FROM_PACKAGE = "${azurerm_storage_blob.source_code_azure.url}"
  }
}
resource "azurerm_app_service_plan" "test_azure" {
  resource_group_name = azurerm_resource_group.fun-rg.name
  name                = "multyfunobjsvpl3ub6"
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
