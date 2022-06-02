package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	"github.com/multycloud/multy/validate"
)

var vaultMetadata = resources.ResourceMetadata[*resourcespb.VaultArgs, *Vault, *resourcespb.VaultResource]{
	CreateFunc:        CreateVault,
	UpdateFunc:        UpdateVault,
	ReadFromStateFunc: VaultFromState,
	ExportFunc:        func(r *Vault, _ *resources.Resources) (*resourcespb.VaultArgs, bool, error) { return r.Args, true, nil },
	ImportFunc:        NewVault,
	AbbreviatedName:   "kv",
}

type Vault struct {
	resources.ResourceWithId[*resourcespb.VaultArgs]
}

func (r *Vault) GetMetadata() resources.ResourceMetadataInterface {
	return &vaultMetadata
}

func CreateVault(resourceId string, args *resourcespb.VaultArgs, others *resources.Resources) (*Vault, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		rgId, err := NewRg("kv", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return nil, err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewVault(resourceId, args, others)
}

func UpdateVault(resource *Vault, vn *resourcespb.VaultArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func VaultFromState(resource *Vault, _ *output.TfState) (*resourcespb.VaultResource, error) {
	return &resourcespb.VaultResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name: resource.Args.Name,
	}, nil
}

func NewVault(resourceId string, args *resourcespb.VaultArgs, _ *resources.Resources) (*Vault, error) {
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
					ResourceGroupName: GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
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
