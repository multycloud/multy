package vault

import "multy-go/resources/common"

const AzureResourceName = "azurerm_key_vault"

type AzureKeyVault struct {
	common.AzResource `hcl:",squash"`
	Sku               string         `hcl:"sku_name"`
	TenantId          string         `hcl:"tenant_id,expr"`
	AccessPolicy      []AccessPolicy `hcl:"access_policy,blocks"`
}

type AccessPolicy struct {
	TenantId               string   `hcl:"tenant_id,expr"`
	ObjectId               string   `hcl:"object_id,expr"`
	CertificatePermissions []string `hcl:"certificate_permissions"`
	KeyPermissions         []string `hcl:"key_permissions"`
	SecretPermissions      []string `hcl:"secret_permissions"`
}
