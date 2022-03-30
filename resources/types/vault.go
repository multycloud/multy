package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

type Vault struct {
	*resources.CommonResourceParams
	Name string `hcl:"name"`
}

type AzureClientConfig struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=azurerm_client_config"`
}

func (r *Vault) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		return []output.TfBlock{}, nil
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
			}}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *Vault) GetVaultId(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AZURE:
		return fmt.Sprintf("%s.%s.id", vault.AzureResourceName, r.GetTfResourceId(cloud)), nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}

func (r *Vault) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, r.CommonResourceParams.Validate(ctx, cloud)...)
	return errs
}

func (r *Vault) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return "", nil
	case common.AZURE:
		return vault.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
	return "", nil
}
