package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/vault_secret"
	"multy-go/validate"
)

type VaultSecret struct {
	*resources.CommonResourceParams
	Name  string `hcl:"name"`
	Value string `hcl:"value"`
	Vault *Vault `mhcl:"ref=vault"`
}

func (r *VaultSecret) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	if cloud == common.AWS {
		return []interface{}{
			vault_secret.AwsSsmParameter{
				AwsResource: common.AwsResource{
					ResourceName: vault_secret.AwsResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
				},
				Name:  fmt.Sprintf("/%s/%s", r.Vault.Name, r.Name),
				Type:  "SecureString",
				Value: r.Value,
			},
		}
	} else if cloud == common.AZURE {
		return []interface{}{
			vault_secret.AzureKeyVaultSecret{
				AzResource: common.AzResource{
					ResourceName: vault_secret.AzureResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
					Name:         r.Name,
				},
				KeyVaultId: r.Vault.GetVaultId(cloud),
				Value:      r.Value,
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *VaultSecret) Validate(ctx resources.MultyContext) {
	return
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
