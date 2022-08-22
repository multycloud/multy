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

func (r GcpVaultAccessPolicy) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VaultAccessPolicyResource, error) {
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
	statuses := map[string]commonpb.ResourceStatus_Status{}

	var ids []string

	resourceName := output.GetResourceName(vault_access_policy.GoogleSecretManagerSecretIamMember{})
	for _, resource := range state.Resources {
		if resource.Type != resourceName {
			continue
		}
		if strings.HasPrefix(resource.Name, fmt.Sprintf("%s-", r.ResourceId)) {
			ids = append(ids, resource.Name)
		}
	}

	for _, resourceId := range ids {
		if stateResource, exists, err := output.MaybeGetParsedById[vault_access_policy.GoogleSecretManagerSecretIamMember](state, resourceId); exists {
			if err != nil {
				return nil, err
			}
			out.GcpOutputs.SecretManagerSecretIamMembershipId = append(out.GcpOutputs.SecretManagerSecretIamMembershipId, stateResource.ResourceId)
			out.Access = getVaultAccess(stateResource.Role)
			out.Identity = getIdentity(stateResource.Member)
			if out.Access != r.Args.Access || out.Identity != r.Args.Identity {
				statuses[fmt.Sprintf("gcp_secret_manager_secret_iam_member_%s", resourceId)] = commonpb.ResourceStatus_NEEDS_UPDATE
			}
		} else {
			statuses[fmt.Sprintf("gcp_secret_manager_secret_iam_member_%s", resourceId)] = commonpb.ResourceStatus_NEEDS_CREATE
		}
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
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

func getVaultAccess(acl string) resourcespb.VaultAccess_Enum {
	switch acl {
	case vault_access_policy.SecretAccessorRole:
		return resourcespb.VaultAccess_READ
	case vault_access_policy.SecretWriterRole:
		return resourcespb.VaultAccess_WRITE
	case vault_access_policy.SecretOwnerRole:
		return resourcespb.VaultAccess_OWNER
	default:
		return resourcespb.VaultAccess_UNKNOWN
	}
}

func getIdentity(member string) string {
	if strings.HasPrefix(member, "serviceAccount:") {
		return strings.TrimPrefix(member, "serviceAccount:")
	}
	return "unknown"
}

func (r GcpVaultAccessPolicy) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main resource in vault_access_policy for gcp")
}
