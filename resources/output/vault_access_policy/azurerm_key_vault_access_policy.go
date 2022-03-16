package vault_access_policy

import (
	"multy-go/resources/common"
	"multy-go/resources/output/vault"
)

type AzureKeyVaultAccessPolicy struct {
	*common.AzResource                     `hcl:",squash" default:"name=azurerm_key_vault_access_policy"`
	KeyVaultId                             string `hcl:"key_vault_id,expr"`
	*vault.AzureKeyVaultAccessPolicyInline `hcl:",squash"`
}
