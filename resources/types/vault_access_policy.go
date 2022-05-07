package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	"github.com/multycloud/multy/resources/output/vault_access_policy"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/validate"
)

var vaultAccessPolicyMetadata = resources.ResourceMetadata[*resourcespb.VaultAccessPolicyArgs, *VaultAccessPolicy, *resourcespb.VaultAccessPolicyResource]{
	CreateFunc:        CreateVaultAccessPolicy,
	UpdateFunc:        UpdateVaultAccessPolicy,
	ReadFromStateFunc: VaultAccessPolicyFromState,
	ExportFunc: func(r *VaultAccessPolicy, _ *resources.Resources) (*resourcespb.VaultAccessPolicyArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewVaultAccessPolicy,
	AbbreviatedName: "kv",
}

type VaultAccessPolicy struct {
	resources.ChildResourceWithId[*Vault, *resourcespb.VaultAccessPolicyArgs]
	Vault *Vault
}

func (r *VaultAccessPolicy) GetMetadata() resources.ResourceMetadataInterface {
	return &vaultAccessPolicyMetadata
}

func CreateVaultAccessPolicy(resourceId string, args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) (*VaultAccessPolicy, error) {
	return NewVaultAccessPolicy(resourceId, args, others)
}

func UpdateVaultAccessPolicy(resource *VaultAccessPolicy, vn *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) error {
	_, err := NewVaultAccessPolicy(resource.ResourceId, vn, others)
	return err
}

func VaultAccessPolicyFromState(resource *VaultAccessPolicy, _ *output.TfState) (*resourcespb.VaultAccessPolicyResource, error) {
	return &resourcespb.VaultAccessPolicyResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  resource.ResourceId,
			NeedsUpdate: false,
		},
		VaultId:  resource.Args.VaultId,
		Identity: resource.Args.Identity,
		Access:   resource.Args.Access,
	}, nil
}

func NewVaultAccessPolicy(resourceId string, args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) (*VaultAccessPolicy, error) {
	vap := &VaultAccessPolicy{
		ChildResourceWithId: resources.ChildResourceWithId[*Vault, *resourcespb.VaultAccessPolicyArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	v, err := resources.Get[*Vault](resourceId, others, args.VaultId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{vap.NewValidationError(err, "vault_id")})
	}
	vap.Parent = v
	vap.Vault = v

	// This dependency is not modelled directly by the id, so we have to do it this way
	for _, res := range others.ResourceMap {
		if vm, ok := res.(*VirtualMachine); ok && vm.GetAwsIdentity() == args.Identity {
			others.AddDependency(resourceId, vm.GetResourceId())
		}
	}

	return vap, nil
}

func (r *VaultAccessPolicy) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		//return []output.TfBlock{
		//	vault_secret.AwsSsmParameter{
		//		AwsResource: &common.AwsResource{
		//			TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
		//		},
		//		Name:  fmt.Sprintf("/%s/%s", r.Vault.Name, r.Name),
		//		ResourceGroup:  "SecureString",
		//		Value: r.Value,
		//	},
		//}
		return nil, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		vaultId, err := r.Vault.GetVaultId()
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			AzureClientConfig{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}},
			vault_access_policy.AzureKeyVaultAccessPolicy{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				},
				KeyVaultId: vaultId,
				AzureKeyVaultAccessPolicyInline: &vault.AzureKeyVaultAccessPolicyInline{
					TenantId: fmt.Sprintf(
						"data.azurerm_client_config.%s.tenant_id", r.ResourceId,
					),
					ObjectId: "\"" + r.Args.Identity + "\"",
					// fixme
					AzureKeyVaultPermissions: r.GetAccessPolicyRules(),
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

// fix return
func (r *VaultAccessPolicy) GetAccessPolicyRules() *vault.AzureKeyVaultPermissions {
	switch r.GetCloud() {
	case common.AWS:
		switch r.Args.Access {
		case resourcespb.VaultAccess_READ:
			return &vault.AzureKeyVaultPermissions{}
		case resourcespb.VaultAccess_WRITE:
			return &vault.AzureKeyVaultPermissions{}
		case resourcespb.VaultAccess_OWNER:
			return &vault.AzureKeyVaultPermissions{}
		}
	case common.AZURE:
		switch r.Args.Access {
		case resourcespb.VaultAccess_READ:
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"List", "Get"},
			}
		case resourcespb.VaultAccess_WRITE:
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"Set", "Delete"},
			}
		case resourcespb.VaultAccess_OWNER:
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"},
			}
		}

	}
	return nil
}

func (r *VaultAccessPolicy) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if r.Args.Access == resourcespb.VaultAccess_UNKNOWN {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown vault access"), "access"))
	}
	return errs
}

func (r *VaultAccessPolicy) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return vault_secret.AwsResourceName, nil
	case common.AZURE:
		return vault_secret.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *VaultAccessPolicy) GetCloudSpecificLocation() string {
	return r.Vault.GetCloudSpecificLocation()
}
