package deploy

import (
	"fmt"
	"os"
)

const (
	AwsProviderVersion    = "4.8.0"
	AzureProviderVersion  = "3.0.2"
	RandomProviderVersion = "3.1.3"
	LocalProviderVersion  = "2.2.2"
	StateRegion           = "eu-west-2"
	tfState               = "terraform.tfstate"
)

func GetTerraformBlock(userId string) (string, error) {
	userStorageName, exists := os.LookupEnv("USER_STORAGE_NAME")
	if !exists {
		return "", fmt.Errorf("USER_STORAGE_NAME not found")
	}

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
    random = {
      source  = "hashicorp/random"
      version = "%s"
    }
    local = {
      source  = "hashicorp/local"
      version = "%s"
    }
  }
}
`, userStorageName, userId, tfState, StateRegion, AwsProviderVersion, AzureProviderVersion, RandomProviderVersion, LocalProviderVersion), nil
}
