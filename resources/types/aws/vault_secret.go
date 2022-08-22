package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/resources/types"
	"strings"
)

type AwsVaultSecret struct {
	*types.VaultSecret
}

func InitVaultSecret(vn *types.VaultSecret) resources.ResourceTranslator[*resourcespb.VaultSecretResource] {
	return AwsVaultSecret{vn}
}

func (r AwsVaultSecret) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VaultSecretResource, error) {
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
	if stateResource, exists, err := output.MaybeGetParsedById[vault_secret.AwsSsmParameter](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = strings.TrimPrefix(stateResource.Name, fmt.Sprintf("/%s/", r.Vault.Args.Name))
		out.Value = stateResource.Value
		out.AwsOutputs = &resourcespb.VaultSecretAwsOutputs{SsmParameterArn: stateResource.Arn}
		output.AddToStatuses(statuses, "aws_ssm_parameter", output.MaybeGetPlannedChageById[vault_secret.AwsSsmParameter](plan, r.ResourceId))
	} else {
		statuses["aws_ssm_parameter"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil

}

func (r AwsVaultSecret) Translate(resources.MultyContext) ([]output.TfBlock, error) {
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
}

func (r AwsVaultSecret) GetMainResourceName() (string, error) {
	return vault_secret.AwsResourceName, nil
}
