package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/vault"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

type Vault struct {
	*resources.CommonResourceParams
	Name string `hcl:"name"`
}

type AzureClientConfig struct {
	common.AzResource `hcl:",squash"`
}

func (r *Vault) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []any {
	if cloud == common.AWS {
		return []any{}
	} else if cloud == common.AZURE {
		config := output.DataSourceWrapper{R: AzureClientConfig{AzResource: common.AzResource{
			ResourceName: "azurerm_client_config",
			ResourceId:   r.GetTfResourceId(cloud),
		}}}

		return []any{
			config,
			vault.AzureKeyVault{
				AzResource: common.AzResource{
					ResourceName:      vault.AzureResourceName,
					ResourceId:        r.GetTfResourceId(cloud),
					Name:              r.Name,
					ResourceGroupName: rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				},
				Sku:      "standard",
				TenantId: fmt.Sprintf("data.azurerm_client_config.%s.tenant_id", r.GetTfResourceId(cloud)),
				AccessPolicy: []vault.AccessPolicy{{
					TenantId:               fmt.Sprintf("data.azurerm_client_config.%s.tenant_id", r.GetTfResourceId(cloud)),
					ObjectId:               fmt.Sprintf("data.azurerm_client_config.%s.object_id", r.GetTfResourceId(cloud)),
					CertificatePermissions: []string{},
					KeyPermissions:         []string{},
					SecretPermissions:      []string{"List", "Get", "Set", "Delete"},
				}},
			}}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *Vault) GetVaultId(cloud common.CloudProvider) string {
	switch cloud {
	case common.AZURE:
		return fmt.Sprintf("%s.%s.id", vault.AzureResourceName, r.GetTfResourceId(cloud))
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}

func (r *Vault) Validate(ctx resources.MultyContext) {
	return
}

func (r *Vault) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return ""
	case common.AZURE:
		return vault.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
