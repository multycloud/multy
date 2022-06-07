package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
)

type AwsVaultAccessPolicy struct {
	*types.VaultAccessPolicy
}

func InitVaultAccessPolicy(vn *types.VaultAccessPolicy) resources.ResourceTranslator[*resourcespb.VaultAccessPolicyResource] {
	return AwsVaultAccessPolicy{vn}
}

func (r AwsVaultAccessPolicy) FromState(state *output.TfState) (*resourcespb.VaultAccessPolicyResource, error) {
	return &resourcespb.VaultAccessPolicyResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		VaultId:  r.Args.VaultId,
		Identity: r.Args.Identity,
		Access:   r.Args.Access,
	}, nil
}

func (r AwsVaultAccessPolicy) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return nil, nil
}

func (r AwsVaultAccessPolicy) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("vault access policy doesn't output any resources in AWS")
}
