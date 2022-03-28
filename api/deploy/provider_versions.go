package deploy

import "fmt"

const AwsProviderVersion = "4.8.0"
const AzureProviderVersion = "4.8.0"

var RequiredProviders = fmt.Sprintf(`
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "%s"
    }
	azurerm = {
      source = "hashicorp/azurerm"
      version = "%s"
	}
  }
}
`, AwsProviderVersion, AzureProviderVersion)
