package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault"
	"github.com/multycloud/multy/resources/output/vault_access_policy"
	"github.com/multycloud/multy/resources/types"
)

type AzureVaultAccessPolicy struct {
	*types.VaultAccessPolicy
}

func InitVaultAccessPolicy(vn *types.VaultAccessPolicy) resources.ResourceTranslator[*resourcespb.VaultAccessPolicyResource] {
	return AzureVaultAccessPolicy{vn}
}

func (r AzureVaultAccessPolicy) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VaultAccessPolicyResource, error) {
	out := &resourcespb.VaultAccessPolicyResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		VaultId:  r.Args.VaultId,
		Identity: r.Args.Identity,
		Access:   r.Args.Access,
	}

	if flags.DryRun {
		return out, nil
	}

	// TODO: parse access
	statuses := map[string]commonpb.ResourceStatus_Status{}
	if stateResource, exists, err := output.MaybeGetParsedById[vault_access_policy.AzureKeyVaultAccessPolicy](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.AzureOutputs = &resourcespb.VaultAccessPolicyAzureOutputs{KeyVaultAccessPolicyId: stateResource.ResourceId}
		out.Identity = stateResource.ObjectId
	} else {
		statuses["azure_key_vault_access_policy"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureVaultAccessPolicy) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		AzureClientConfig{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}},
		vault_access_policy.AzureKeyVaultAccessPolicy{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			KeyVaultId: fmt.Sprintf("%s.%s.id", vault.AzureResourceName, r.Vault.ResourceId),
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

func (r AzureVaultAccessPolicy) GetAccessPolicyRules() *vault.AzureKeyVaultPermissions {
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
	default:
		return nil
	}
}

func (r AzureVaultAccessPolicy) GetMainResourceName() (string, error) {
	return output.GetResourceName(vault_access_policy.AzureKeyVaultAccessPolicy{}), nil
}
