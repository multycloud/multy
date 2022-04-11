package deploy

import "fmt"

const (
	AwsProviderVersion   = "4.8.0"
	AzureProviderVersion = "3.0.2"
	StateRegion          = "eu-west-2"
	Bucket               = "multy-users-tfstate"
	tfState              = "terraform.tfstate"
)

func GetTerraformBlock(userId string) string {
	return fmt.Sprintf(`
terraform {
  backend "s3" {
    bucket         = "%s"
    key            = "%s/%s"
    region         = "%s"
  }
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
`, Bucket, userId, tfState, StateRegion, AwsProviderVersion, AzureProviderVersion)
}
