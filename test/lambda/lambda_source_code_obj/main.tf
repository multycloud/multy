resource "aws_lambda_function" "function2_aws" {
  tags =  {
    Name = "publicmultyfun"
  }

  function_name = "publicmultyfun"
  role          = aws_iam_role.iam_for_lambda_function2.arn
  runtime       = "python3.9"
  handler       = "lambda_function.lambda_handler"
  s3_bucket     = aws_s3_bucket.obj_storage_aws.id
  s3_key        = aws_s3_bucket_object.public_source_code_aws.key
}
resource "aws_iam_role" "iam_for_lambda_function2" {
  tags =  {
    Name = "iam_for_lambda_function2"
  }

  name               = "iam_for_lambda_function2"
  assume_role_policy = "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}"
}
resource "aws_s3_bucket" "obj_storage_aws" {
  bucket = "function-storage-1722"
}
resource "aws_s3_bucket_object" "public_source_code_aws" {
  bucket = aws_s3_bucket.obj_storage_aws.id
  key    = "source_code.zip"
  acl    = "public-read"
  source = "source_dir/aws_code.zip"
}
resource "azurerm_resource_group" "fun-rg" {
  name     = "fun-rg"
  location = "northeurope"
}
resource "azurerm_function_app" "function2_azure" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "publicmultyfun"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.obj_storage_azure.name
  storage_account_access_key = azurerm_storage_account.obj_storage_azure.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.function2_azure.id
  os_type                    = "linux"

  app_settings =  {
    WEBSITE_RUN_FROM_PACKAGE = "${azurerm_storage_blob.public_source_code_azure.url}"
  }
}
resource "azurerm_app_service_plan" "function2_azure" {
  resource_group_name = azurerm_resource_group.fun-rg.name
  name                = "publicmultyfunsvpl492v"
  location            = "northeurope"
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
resource "azurerm_storage_account" "obj_storage_azure" {
  resource_group_name      = azurerm_resource_group.st-rg.name
  name                     = "functionstorage1722"
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
resource "azurerm_storage_blob" "public_source_code_azure" {
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
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}
