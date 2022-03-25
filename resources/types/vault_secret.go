package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/validate"
)

type VaultSecret struct {
	*resources.CommonResourceParams
	Name  string `hcl:"name"`
	Value string `hcl:"value"`
	Vault *Vault `mhcl:"ref=vault"`
}

func (r *VaultSecret) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		return []output.TfBlock{
			vault_secret.AwsSsmParameter{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				Name:  fmt.Sprintf("/%s/%s", r.Vault.Name, r.Name),
				Type:  "SecureString",
				Value: r.Value,
			},
		}, nil
	} else if cloud == common.AZURE {
		vaultId, err := r.Vault.GetVaultId(cloud)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			vault_secret.AzureKeyVaultSecret{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
					Name:              r.Name,
				},
				KeyVaultId: vaultId,
				Value:      r.Value,
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *VaultSecret) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *VaultSecret) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return vault_secret.AwsResourceName, nil
	case common.AZURE:
		return vault_secret.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
