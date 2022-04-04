package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

type Vault struct {
	resources.ResourceWithId[*resourcespb.VaultArgs]
}

func NewVault(resourceId string, args *resourcespb.VaultArgs, _ resources.Resources) (*Vault, error) {
	return &Vault{
		ResourceWithId: resources.ResourceWithId[*resourcespb.VaultArgs]{resourceId, args},
	}, nil
}

type AzureClientConfig struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=azurerm_client_config"`
}

func (r *Vault) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return []output.TfBlock{
			AzureClientConfig{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}},
			vault.AzureKeyVault{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					Name:              r.Args.Name,
					ResourceGroupName: rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
					Location:          r.GetCloudSpecificLocation(),
				},
				Sku:      "standard",
				TenantId: fmt.Sprintf("data.azurerm_client_config.%s.tenant_id", r.ResourceId),
				AccessPolicy: []vault.AzureKeyVaultAccessPolicyInline{{
					TenantId: fmt.Sprintf(
						"data.azurerm_client_config.%s.tenant_id", r.ResourceId,
					),
					ObjectId: fmt.Sprintf(
						"data.azurerm_client_config.%s.object_id", r.ResourceId,
					),
					AzureKeyVaultPermissions: &vault.AzureKeyVaultPermissions{
						CertificatePermissions: []string{},
						KeyPermissions:         []string{},
						SecretPermissions:      []string{"List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"},
					},
				}},
			}}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *Vault) GetVaultId() (string, error) {
	switch r.GetCloud() {
	case common.AZURE:
		return fmt.Sprintf("%s.%s.id", vault.AzureResourceName, r.ResourceId), nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *Vault) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}

func (r *Vault) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return "", nil
	case common.AZURE:
		return vault.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
