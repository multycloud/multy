package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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

func (r AzureVaultSecret) FromState(state *output.TfState) (*resourcespb.VaultSecretResource, error) {
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
