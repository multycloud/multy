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
	*output.TerraformDataSource `hcl:",squash" default:"name=azurerm_client_config"`
}

func (r *Vault) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		return []output.TfBlock{}
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			AzureClientConfig{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.GetTfResourceId(cloud)}},
			vault.AzureKeyVault{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
					Name:              r.Name,
					ResourceGroupName: rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				},
				Sku:      "standard",
				TenantId: fmt.Sprintf("data.azurerm_client_config.%s.tenant_id", r.GetTfResourceId(cloud)),
				AccessPolicy: []vault.AzureKeyVaultAccessPolicyInline{{
					TenantId: fmt.Sprintf(
						"data.azurerm_client_config.%s.tenant_id", r.GetTfResourceId(cloud),
					),
					ObjectId: fmt.Sprintf(
						"data.azurerm_client_config.%s.object_id", r.GetTfResourceId(cloud),
					),
					AzureKeyVaultPermissions: &vault.AzureKeyVaultPermissions{
						CertificatePermissions: []string{},
						KeyPermissions:         []string{},
						SecretPermissions:      []string{"List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"},
					},
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

func (r *Vault) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
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

//func (r *Vault) GetCloudId(cloud common.CloudProvider) string {
//	switch cloud {
//	case common.AWS:
//		return fmt.Sprintf("%s.%s.arn", common.GetResourceName(), r.ResourceId)
//	default:
//		validate.LogInternalError("unknown cloud %s", cloud)
//	}
//	return ""
//}
