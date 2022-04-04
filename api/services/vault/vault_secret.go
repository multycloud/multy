package vault

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultSecretService struct {
	Service services.Service[*resourcespb.VaultSecretArgs, *resourcespb.VaultSecretResource]
}

func (s VaultSecretService) Convert(resourceId string, args *resourcespb.VaultSecretArgs, state *output.TfState) (*resourcespb.VaultSecretResource, error) {
	return &resourcespb.VaultSecretResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Value:            args.Value,
		VaultId:          args.VaultId,
	}, nil
}

func NewVaultSecretService(database *db.Database) VaultSecretService {
	ni := VaultSecretService{
		Service: services.Service[*resourcespb.VaultSecretArgs, *resourcespb.VaultSecretResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
