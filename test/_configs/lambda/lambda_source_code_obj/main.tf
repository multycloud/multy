resource "aws_lambda_function" "function2_aws" {
  tags = {
    "Name" = "publicmultyfun"
  }

  function_name = "publicmultyfun"
  role          = aws_iam_role.iam_for_lambda_function2.arn
  runtime       = "python3.9"
  handler       = "lambda_function.lambda_handler"
  s3_bucket     = aws_s3_bucket.obj_storage_aws.id
  s3_key        = aws_s3_bucket_object.public_source_code_aws.key
}
resource "aws_iam_role" "iam_for_lambda_function2" {
  tags = {
    "Name" = "iam_for_lambda_function2"
  }

  name               = "iam_for_lambda_function2"
  assume_role_policy =  "{\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"lambda.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "function2_aws" {
  role       = aws_iam_role.iam_for_lambda_function2.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
resource "aws_api_gateway_rest_api" "function2_aws" {
  tags = {
    "Name" = "publicmultyfun"
  }

  name        = "publicmultyfun"
  description = ""
}
resource "aws_api_gateway_resource" "function2_proxy" {
  rest_api_id = aws_api_gateway_rest_api.function2_aws.id
  parent_id   = aws_api_gateway_rest_api.function2_aws.root_resource_id
  path_part   = "{proxy+}"
}
resource "aws_api_gateway_method" "function2_proxy" {
  rest_api_id   = aws_api_gateway_rest_api.function2_aws.id
  resource_id   = aws_api_gateway_resource.function2_proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "function2_proxy" {
  rest_api_id             = aws_api_gateway_rest_api.function2_aws.id
  resource_id             = aws_api_gateway_method.function2_proxy.resource_id
  http_method             = aws_api_gateway_method.function2_proxy.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function2_aws.invoke_arn
}
resource "aws_api_gateway_method" "function2_proxy_root" {
  rest_api_id   = aws_api_gateway_rest_api.function2_aws.id
  resource_id   = aws_api_gateway_rest_api.function2_aws.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "function2_proxy_root" {
  rest_api_id             = aws_api_gateway_rest_api.function2_aws.id
  resource_id             = aws_api_gateway_method.function2_proxy_root.resource_id
  http_method             = aws_api_gateway_method.function2_proxy_root.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.function2_aws.invoke_arn
}
resource "aws_api_gateway_deployment" "function2_aws" {
  rest_api_id = aws_api_gateway_rest_api.function2_aws.id
  stage_name  = "api"

  depends_on = [
    aws_api_gateway_integration.function2_proxy,
    aws_api_gateway_integration.function2_proxy_root,
  ]
}
resource "aws_lambda_permission" "function2_aws" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "publicmultyfun"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.function2_aws.execution_arn}/*/*"
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
    "WEBSITE_RUN_FROM_PACKAGE" = "${azurerm_storage_blob.public_source_code_azure.url}"
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