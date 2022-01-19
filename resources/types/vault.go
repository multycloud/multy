package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/vault"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

type Vault struct {
	*resources.CommonResourceParams
	Name string `hcl:"name"`
}

func (r *Vault) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	if cloud == common.AWS {
		return []interface{}{}
	} else if cloud == common.AZURE {
		return []interface{}{
			vault.AzureKeyVault{
				AzResource: common.AzResource{
					ResourceName:      vault.AzureResourceName,
					ResourceId:        r.GetTfResourceId(cloud),
					Name:              r.Name,
					ResourceGroupName: rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				},
				Sku:      "Standard",
				TenantId: "data.azurerm_client_config.current.tenant_id",
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
		return vault.AwsResourceName
	case common.AZURE:
		return vault.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
