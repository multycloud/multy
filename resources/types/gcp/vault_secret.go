package gcp_resources

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

type GcpVaultSecret struct {
	*types.VaultSecret
}

func InitVaultSecret(vn *types.VaultSecret) resources.ResourceTranslator[*resourcespb.VaultSecretResource] {
	return GcpVaultSecret{vn}
}

func (r GcpVaultSecret) FromState(state *output.TfState) (*resourcespb.VaultSecretResource, error) {
	return &resourcespb.VaultSecretResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:        r.Args.Name,
		Value:       r.Args.Value,
		VaultId:     r.Args.VaultId,
		GcpOverride: r.Args.GcpOverride,
	}, nil
}

func (r GcpVaultSecret) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	replication := vault_secret.GoogleSecretManagerSecretReplication{Automatic: true}
	if !r.Args.GetGcpOverride().GetGlobalReplication() {
		replication = vault_secret.NewManagedSecretReplication(r.Vault.GetCloudSpecificLocation())
	}
	secret := vault_secret.GoogleSecretManagerSecret{
		GcpResource: common.NewGcpResource(r.ResourceId, "", r.Vault.Args.GetGcpOverride().GetProject()),
		SecretId:    r.Args.Name,
		Replication: replication,
	}

	version := vault_secret.GoogleSecretManagerSecretVersion{
		GcpResource: common.NewGcpResourceWithNoProject(r.ResourceId, ""),
		SecretId:    fmt.Sprintf("%s.%s.id", output.GetResourceName(secret), r.ResourceId),
		SecretData:  r.Args.Value,
	}

	return []output.TfBlock{&secret, &version}, nil
}

func (r GcpVaultSecret) GetMainResourceName() (string, error) {
	return output.GetResourceName(vault_secret.GoogleSecretManagerSecret{}), nil
}
