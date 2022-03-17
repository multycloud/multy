package vault_secret

import "github.com/multycloud/multy/resources/common"

const AzureResourceName = "azurerm_key_vault_secret"

type AzureKeyVaultSecret struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_key_vault_secret"`
	KeyVaultId         string `hcl:"key_vault_id,expr"`
	Value              string `hcl:"value"`
}
