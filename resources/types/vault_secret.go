package types

import (
	"fmt"
	"multy/resources"
	"multy/resources/common"
	"multy/resources/output"
	"multy/resources/output/vault_secret"
	"multy/validate"
)

type VaultSecret struct {
	*resources.CommonResourceParams
	Name  string `hcl:"name"`
	Value string `hcl:"value"`
	Vault *Vault `mhcl:"ref=vault"`
}

func (r *VaultSecret) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
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
		}
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			vault_secret.AzureKeyVaultSecret{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
					Name:              r.Name,
				},
				KeyVaultId: r.Vault.GetVaultId(cloud),
				Value:      r.Value,
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *VaultSecret) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *VaultSecret) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return vault_secret.AwsResourceName
	case common.AZURE:
		return vault_secret.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
