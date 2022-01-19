package vault

import "multy-go/resources/common"

const AzureResourceName = "azurerm_key_vault"

type AzureKeyVault struct {
	common.AzResource `hcl:",squash"`
	Sku               string `hcl:"sku_name"`
	TenantId          string `hcl:"tenant_id,expr"`
	//tenant_id = data.azurerm_client_config.current.tenant_id
	//tenant_id = data.azurerm_client_config.current.tenant_id
	//object_id = data.azurerm_client_config.current.object_id
}
