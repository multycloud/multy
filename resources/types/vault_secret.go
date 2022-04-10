package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/validate"
)

type VaultSecret struct {
	resources.ChildResourceWithId[*Vault, *resourcespb.VaultSecretArgs]

	Vault *Vault
}

func NewVaultSecret(resourceId string, args *resourcespb.VaultSecretArgs, others resources.Resources) (*VaultSecret, error) {
	vap := &VaultSecret{
		ChildResourceWithId: resources.ChildResourceWithId[*Vault, *resourcespb.VaultSecretArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	v, err := resources.Get[*Vault](resourceId, others, args.VaultId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{vap.NewValidationError(err.Error(), "vault_id")})
	}
	vap.Parent = v
	vap.Vault = v
	return vap, nil
}

func (r *VaultSecret) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{
			vault_secret.AwsSsmParameter{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				},
				Name:  fmt.Sprintf("/%s/%s", r.Vault.Args.Name, r.Args.Name),
				Type:  "SecureString",
				Value: r.Args.Value,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		vaultId, err := r.Vault.GetVaultId()
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			vault_secret.AzureKeyVaultSecret{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					Name:              r.Args.Name,
				},
				KeyVaultId: vaultId,
				Value:      r.Args.Value,
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *VaultSecret) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	return errs
}

func (r *VaultSecret) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return vault_secret.AwsResourceName, nil
	case common.AZURE:
		return vault_secret.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *VaultSecret) GetCloudSpecificLocation() string {
	return r.Vault.GetCloudSpecificLocation()
}
