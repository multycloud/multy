package object_storage

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

const AzureResourceName = "azurerm_storage_account"

// azurerm_storage_account
type AzureStorageAccount struct {
	*common.AzResource     `hcl:",squash" default:"name=azurerm_storage_account"`
	AccountTier            string `hcl:"account_tier"`
	AccountReplicationType string `hcl:"account_replication_type"`
	AllowBlobPublicAccess  bool   `hcl:"allow_nested_items_to_be_public"`
}

func (r AzureStorageAccount) GetResourceName() string {
	return fmt.Sprintf("azurerm_storage_account.%s.name", r.ResourceId)
}
