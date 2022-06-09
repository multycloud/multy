package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/vault_secret"
	"github.com/multycloud/multy/resources/types"
)

type AwsVaultSecret struct {
	*types.VaultSecret
}

func InitVaultSecret(vn *types.VaultSecret) resources.ResourceTranslator[*resourcespb.VaultSecretResource] {
	return AwsVaultSecret{vn}
}

func (r AwsVaultSecret) FromState(state *output.TfState) (*resourcespb.VaultSecretResource, error) {
	return &resourcespb.VaultSecretResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:    r.Args.Name,
		Value:   r.Args.Value,
		VaultId: r.Args.VaultId,
	}, nil
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
