package vault_secret

import "multy-go/resources/common"

const AzureResourceName = "azurerm_key_vault_secret"

type AzureKeyVaultSecret struct {
	common.AzResource `hcl:",squash"`
	KeyVaultId        string `hcl:"key_vault_id,expr"`
	Value             string `hcl:"value"`
}
