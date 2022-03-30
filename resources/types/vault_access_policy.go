package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	"github.com/multycloud/multy/resources/output/vault_access_policy"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/validate"
)

type VaultAccessPolicy struct {
	*resources.CommonResourceParams
	Vault    *Vault `mhcl:"ref=vault"`
	Identity string `hcl:"identity"`
	Access   string `hcl:"access"`
}

const (
	READ  = "read"
	WRITE = "write"
	OWNER = "owner"
)

func (r *VaultAccessPolicy) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		//return []output.TfBlock{
		//	vault_secret.AwsSsmParameter{
		//		AwsResource: &common.AwsResource{
		//			TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
		//		},
		//		Name:  fmt.Sprintf("/%s/%s", r.Vault.Name, r.Name),
		//		Type:  "SecureString",
		//		Value: r.Value,
		//	},
		//}
		return nil, nil
	} else if cloud == common.AZURE {
		vaultId, err := r.Vault.GetVaultId(cloud)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			AzureClientConfig{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.GetTfResourceId(cloud)}},
			vault_access_policy.AzureKeyVaultAccessPolicy{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				KeyVaultId: vaultId,
				AzureKeyVaultAccessPolicyInline: &vault.AzureKeyVaultAccessPolicyInline{
					TenantId: fmt.Sprintf(
						"data.azurerm_client_config.%s.tenant_id", r.GetTfResourceId(cloud),
					),
					ObjectId: "\"" + r.Identity + "\"",
					// fixme
					AzureKeyVaultPermissions: r.GetAccessPolicyRules(cloud),
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

// fix return
func (r *VaultAccessPolicy) GetAccessPolicyRules(cloud common.CloudProvider) *vault.AzureKeyVaultPermissions {
	switch cloud {
	case common.AWS:
		switch r.Access {
		case "read":
			return &vault.AzureKeyVaultPermissions{}
		case "write":
			return &vault.AzureKeyVaultPermissions{}
		case "owner":
			return &vault.AzureKeyVaultPermissions{}
		}
	case common.AZURE:
		switch r.Access {
		case "read":
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"List", "Get"},
			}
		case "write":
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"Set", "Delete"},
			}
		case "owner":
			return &vault.AzureKeyVaultPermissions{
				CertificatePermissions: []string{},
				KeyPermissions:         []string{},
				SecretPermissions:      []string{"List", "Get", "Set", "Delete", "Recover", "Backup", "Restore", "Purge"},
			}
		}

	}
	return nil
}

func (r *VaultAccessPolicy) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, r.CommonResourceParams.Validate(ctx, cloud)...)
	if r.Access != READ && r.Access != OWNER && r.Access != WRITE {
		errs = append(errs, r.NewError("access", fmt.Sprintf("%s access is invalid", r.ResourceId)))
	}
	return errs
}

func (r *VaultAccessPolicy) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return vault_secret.AwsResourceName, nil
	case common.AZURE:
		return vault_secret.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}

func (r *VaultAccessPolicy) GetLocation(cloud common.CloudProvider, ctx resources.MultyContext) string {
	return r.Vault.GetLocation(cloud, ctx)
}
