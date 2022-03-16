package vault

import "multy-go/resources/common"

const AzureResourceName = "azurerm_key_vault"

type AzureKeyVault struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_key_vault"`
	Sku                string                            `hcl:"sku_name"`
	TenantId           string                            `hcl:"tenant_id,expr"`
	AccessPolicy       []AzureKeyVaultAccessPolicyInline `hcl:"access_policy,blocks"`
}

type AzureKeyVaultAccessPolicyInline struct {
	TenantId                  string `hcl:"tenant_id,expr"`
	ObjectId                  string `hcl:"object_id,expr"`
	*AzureKeyVaultPermissions `hcl:",squash"`
}

type AzureKeyVaultPermissions struct {
	CertificatePermissions []string `hcl:"certificate_permissions"`
	KeyPermissions         []string `hcl:"key_permissions"`
	SecretPermissions      []string `hcl:"secret_permissions"`
}
