package object_storage

import (
	"fmt"
	"multy-go/resources/common"
)

const AzureResourceName = "azurerm_storage_account"

// azurerm_storage_account
type AzureStorageAccount struct {
	common.AzResource      `hcl:",squash"`
	AccountTier            string `hcl:"account_tier"`
	AccountReplicationType string `hcl:"account_replication_type"`
	AllowBlobPublicAccess  bool   `hcl:"allow_blob_public_access"`
}

func (r AzureStorageAccount) GetResourceName() string {
	return fmt.Sprintf("azurerm_storage_account.%s.name", r.ResourceId)
}
