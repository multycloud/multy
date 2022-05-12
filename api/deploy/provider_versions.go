package deploy

import (
	"fmt"
	"github.com/multycloud/multy/flags"
	"os"
	"path"
)

const (
	awsProviderVersion    = "4.8.0"
	azureProviderVersion  = "3.0.2"
	randomProviderVersion = "3.1.3"
	localProviderVersion  = "2.2.2"
	tfStateRegion         = "eu-west-2"
	tfState               = "terraform.tfstate"
)

func GetTerraformBlock(userId string) (string, error) {
	if flags.Environment == flags.Local {
		return getTerraformBlock(getLocalStateBlock(userId))
	}
	userStorageName, exists := os.LookupEnv("USER_STORAGE_NAME")
	if !exists {
		return "", fmt.Errorf("USER_STORAGE_NAME not found")
	}
	return getTerraformBlock(getRemoteStateBlock(userId, userStorageName))
}

func getLocalStateBlock(userId string) string {
	p := path.Join(os.TempDir(), userId, tfState)
	return fmt.Sprintf(
		`backend "local" {
    path = "%s"
}`, p)
}
func getRemoteStateBlock(userId, userStorageName string) string {
	return fmt.Sprintf(`backend "s3" {
    bucket         = "%s"
    key            = "%s/%s"
    region         = "%s"
    profile        = "multy"
  }`, userId, userStorageName, tfState, tfStateRegion)
}

func getTerraformBlock(providerBlock string) (string, error) {
	return fmt.Sprintf(`
terraform {
  %s
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
`, providerBlock, awsProviderVersion, azureProviderVersion, randomProviderVersion, localProviderVersion), nil
}
