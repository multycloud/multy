data "archive_file" "super_long_function_aws" {
  type        = "zip"
  source_dir  = "source_dir/aws"
  output_path = ".multy/tmp/super_long_function_aws.zip"
}
resource "aws_lambda_function" "super_long_function_aws" {
  tags =  {
    Name = "super_long_function"
  }

  function_name    = "super_long_function"
  role             = aws_iam_role.iam_for_lambda_super_long_function.arn
  filename         = ".multy/tmp/super_long_function_aws.zip"
  source_code_hash = data.archive_file.super_long_function_aws.output_base64sha256
  runtime          = "python3.9"
  handler          = "lambda_function.lambda_handler"
}
resource "aws_iam_role" "iam_for_lambda_super_long_function" {
  tags =  {
    Name = "iam_for_lambda_super_long_function"
  }

  name               = "iam_for_lambda_super_long_function"
  assume_role_policy = "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}"
}
resource "azurerm_resource_group" "fun-rg" {
  name     = "fun-rg"
  location = "northeurope"
}
data "archive_file" "super_long_function_azure" {
  type        = "zip"
  source_dir  = "source_dir/azure"
  output_path = ".multy/tmp/super_long_function_azure.zip"
}
resource "azurerm_storage_account" "super_long_function_azure" {
  resource_group_name      = azurerm_resource_group.fun-rg.name
  name                     = "superlongfunit2xstaciiav"
  location                 = "northeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = false
}
resource "azurerm_function_app" "super_long_function_azure" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "superlongfunction"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.super_long_function_azure.name
  storage_account_access_key = azurerm_storage_account.super_long_function_azure.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.super_long_function_azure.id
  os_type                    = "linux"

  provisioner "local-exec" {
    command     = "az functionapp deployment source config-zip -g ${self.resource_group_name} -n ${self.name} --src ${data.archive_file.super_long_function_azure.output_path}"
  }
}
resource "azurerm_app_service_plan" "super_long_function_azure" {
  resource_group_name = azurerm_resource_group.fun-rg.name
  name                = "superlongfunit2xsvpl15kn"
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
