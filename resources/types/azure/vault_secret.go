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
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/resources/types"
)

type AzureVaultSecret struct {
	*types.VaultSecret
}

func InitVaultSecret(vn *types.VaultSecret) resources.ResourceTranslator[*resourcespb.VaultSecretResource] {
	return AzureVaultSecret{vn}
}

func (r AzureVaultSecret) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VaultSecretResource, error) {
	out := &resourcespb.VaultSecretResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:        r.Args.Name,
		Value:       r.Args.Value,
		VaultId:     r.Args.VaultId,
		GcpOverride: r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[vault_secret.AzureKeyVaultSecret](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.Value = stateResource.Value
		out.AzureOutputs = &resourcespb.VaultSecretAzureOutputs{KeyVaultSecretId: stateResource.ResourceId}
		output.AddToStatuses(statuses, "azure_key_vault_secret", output.MaybeGetPlannedChageById[vault_secret.AzureKeyVaultSecret](plan, r.ResourceId))
	} else {
		statuses["azure_key_vault_secret"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureVaultSecret) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		vault_secret.AzureKeyVaultSecret{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				Name:              r.Args.Name,
			},
			KeyVaultId: fmt.Sprintf("%s.%s.id", vault.AzureResourceName, r.Vault.ResourceId),
			Value:      r.Args.Value,
		},
	}, nil
}

func (r AzureVaultSecret) GetMainResourceName() (string, error) {
	return vault_secret.AzureResourceName, nil
}
