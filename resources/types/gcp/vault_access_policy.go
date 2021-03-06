package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault_access_policy"
	"github.com/multycloud/multy/resources/types"
	"strings"
)

type GcpVaultAccessPolicy struct {
	*types.VaultAccessPolicy
}

func InitVaultAccessPolicy(vn *types.VaultAccessPolicy) resources.ResourceTranslator[*resourcespb.VaultAccessPolicyResource] {
	return GcpVaultAccessPolicy{vn}
}

func (r GcpVaultAccessPolicy) FromState(state *output.TfState) (*resourcespb.VaultAccessPolicyResource, error) {
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

	out.GcpOutputs = &resourcespb.VaultAccessPolicyGcpOutputs{}

	var ids []string

	prefix := output.GetResourceName(vault_access_policy.GoogleSecretManagerSecretIamMember{}) + "."
	for _, resource := range state.Values.RootModule.Resources {
		if strings.HasPrefix(resource.Address,
			fmt.Sprintf("%s%s-", prefix, r.ResourceId)) {
			ids = append(ids, strings.TrimPrefix(resource.Address, prefix))
		}
	}

	for _, resourceId := range ids {
		stateResource, err := output.GetParsedById[vault_access_policy.GoogleSecretManagerSecretIamMember](state, resourceId)
		if err != nil {
			return nil, err
		}
		out.GcpOutputs.SecretManagerSecretIamMembershipId = append(out.GcpOutputs.SecretManagerSecretIamMembershipId, stateResource.ResourceId)
	}

	return out, nil
}

func (r GcpVaultAccessPolicy) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	role, err := getGcpIamRole(r.Args.Access)
	if err != nil {
		return nil, err
	}

	secrets := resources.GetAllResourcesWithRef(ctx, func(t *types.VaultSecret) *types.Vault {
		return t.Vault
	}, r.Vault)

	var result []output.TfBlock
	for _, secret := range secrets {
		result = append(result,
			&vault_access_policy.GoogleSecretManagerSecretIamMember{
				GcpResource: common.NewGcpResource(
					fmt.Sprintf("%s-%s", r.ResourceId, secret.ResourceId), "", r.Vault.Args.GetGcpOverride().GetProject()),
				SecretId: GcpVaultSecret{secret}.getSecretId(),
				Role:     role,
				Member:   fmt.Sprintf("serviceAccount:%s", r.Args.Identity),
			})
	}

	return result, nil
}

func getGcpIamRole(acl resourcespb.VaultAccess_Enum) (string, error) {
	switch acl {
	case resourcespb.VaultAccess_READ:
		return vault_access_policy.SecretAccessorRole, nil
	case resourcespb.VaultAccess_WRITE:
		return vault_access_policy.SecretWriterRole, nil
	case resourcespb.VaultAccess_OWNER:
		return vault_access_policy.SecretOwnerRole, nil
	default:
		return "", fmt.Errorf("unknown vault access: %s", acl.String())
	}
}

func (r GcpVaultAccessPolicy) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main resource in vault_access_policy for gcp")
}
