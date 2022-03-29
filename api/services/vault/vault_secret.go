package vault

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultSecretService struct {
	Service services.Service[*resources.VaultSecretArgs, *resources.VaultSecretResource]
}

func (s VaultSecretService) Convert(resourceId string, args *resources.VaultSecretArgs, state *output.TfState) (*resources.VaultSecretResource, error) {
	return &resources.VaultSecretResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Value:            args.Value,
		VaultId:          args.VaultId,
	}, nil
}

func NewVaultSecretService(database *db.Database) VaultSecretService {
	ni := VaultSecretService{
		Service: services.Service[*resources.VaultSecretArgs, *resources.VaultSecretResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
