data "archive_file" "super_long_function_aws" {
  type        = "zip"
  source_dir  = "source_dir/aws"
  output_path = ".multy/tmp/super_long_function_aws.zip"
}
resource "aws_lambda_function" "super_long_function_aws" {
  tags = {
    "Name" = "super_long_function"
  }

  function_name    = "super_long_function"
  role             = aws_iam_role.iam_for_lambda_super_long_function.arn
  filename         = ".multy/tmp/super_long_function_aws.zip"
  source_code_hash = data.archive_file.super_long_function_aws.output_base64sha256
  runtime          = "python3.9"
  handler          = "lambda_function.lambda_handler"
}
resource "aws_iam_role" "iam_for_lambda_super_long_function" {
  tags = {
    "Name" = "iam_for_lambda_super_long_function"
  }

  name               = "iam_for_lambda_super_long_function"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"lambda.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "super_long_function_aws" {
  role       = aws_iam_role.iam_for_lambda_super_long_function.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
resource "aws_api_gateway_rest_api" "super_long_function_aws" {
  tags = {
    "Name" = "super_long_function"
  }

  name        = "super_long_function"
  description = ""
}
resource "aws_api_gateway_resource" "super_long_function_proxy" {
  rest_api_id = aws_api_gateway_rest_api.super_long_function_aws.id
  parent_id   = aws_api_gateway_rest_api.super_long_function_aws.root_resource_id
  path_part   = "{proxy+}"
}
resource "aws_api_gateway_method" "super_long_function_proxy" {
  rest_api_id   = aws_api_gateway_rest_api.super_long_function_aws.id
  resource_id   = aws_api_gateway_resource.super_long_function_proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "super_long_function_proxy" {
  rest_api_id             = aws_api_gateway_rest_api.super_long_function_aws.id
  resource_id             = aws_api_gateway_method.super_long_function_proxy.resource_id
  http_method             = aws_api_gateway_method.super_long_function_proxy.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.super_long_function_aws.invoke_arn
}
resource "aws_api_gateway_method" "super_long_function_proxy_root" {
  rest_api_id   = aws_api_gateway_rest_api.super_long_function_aws.id
  resource_id   = aws_api_gateway_rest_api.super_long_function_aws.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "super_long_function_proxy_root" {
  rest_api_id             = aws_api_gateway_rest_api.super_long_function_aws.id
  resource_id             = aws_api_gateway_method.super_long_function_proxy_root.resource_id
  http_method             = aws_api_gateway_method.super_long_function_proxy_root.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.super_long_function_aws.invoke_arn
}
resource "aws_api_gateway_deployment" "super_long_function_aws" {
  rest_api_id = aws_api_gateway_rest_api.super_long_function_aws.id
  stage_name  = "api"

  depends_on = [
    aws_api_gateway_integration.super_long_function_proxy,
    aws_api_gateway_integration.super_long_function_proxy_root,
  ]
}
resource "aws_lambda_permission" "super_long_function_aws" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "super_long_function"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.super_long_function_aws.execution_arn}/*/*"
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
  resource_group_name             = azurerm_resource_group.fun-rg.name
  name                            = "superlongfunit2xstaciiav"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
  blob_properties {
    versioning_enabled = false
  }
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
    command = "az functionapp deployment source config-zip -g ${self.resource_group_name} -n ${self.name} --src ${data.archive_file.super_long_function_azure.output_path}"
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
