package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
)

type AwsVault struct {
	*types.Vault
}

func InitVault(vn *types.Vault) resources.ResourceTranslator[*resourcespb.VaultResource] {
	return AwsVault{vn}
}

type awsClientConfig struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=awsrm_client_config"`
}

func (r AwsVault) FromState(state *output.TfState) (*resourcespb.VaultResource, error) {
	return &resourcespb.VaultResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name: r.Args.Name,
	}, nil
}

func (r AwsVault) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return nil, nil
}

func (r AwsVault) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("vaults doesn't output any resources in AWS")
}
