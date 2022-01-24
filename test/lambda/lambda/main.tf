resource "aws_lambda_function" "test_aws" {
  tags = {
    Name = "test_name"
  }
  function_name    = "test_name"
  filename         = ".multy/tmp/test_name.zip"
  source_code_hash = data.archive_file.test_aws.output_base64sha256
  role             = aws_iam_role.iam_for_lambda_test_name.arn
  runtime          = "python3.9"
  handler          = "lambda_function.lambda_handler"
}
resource "aws_iam_role" "iam_for_lambda_test_name" {
  tags = {
    Name = "iam_for_lambda_test_name"
  }
  name               = "iam_for_lambda_test_name"
  assume_role_policy = "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}"
}
data "archive_file" "test_aws" {
  type        = "zip"
  source_dir  = "source_dir/aws"
  output_path = ".multy/tmp/test_name.zip"
}
resource "azurerm_resource_group" "fun-rg" {
  name     = "fun-rg"
  location = "northeurope"
}
resource "azurerm_storage_account" "test" {
  resource_group_name      = azurerm_resource_group.fun-rg.name
  name                     = "teststacct"
  location                 = "northeurope"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_blob_public_access = false
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
resource "azurerm_function_app" "test" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "testname"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.test.id
  os_type                    = "linux"
  provisioner "local-exec" {
    working_dir = "source_dir/azure"
    command = "func azure functionapp publish ${self.id}"
  }
}
provider "aws" {
  region = "eu-west-1"
}
provider "azurerm" {
  features {}
}