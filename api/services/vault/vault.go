package vault

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultService struct {
	Service services.Service[*resources.VaultArgs, *resources.VaultResource]
}

func (s VaultService) Convert(resourceId string, args *resources.VaultArgs, state *output.TfState) (*resources.VaultResource, error) {
	return &resources.VaultResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
	}, nil
}

func NewVaultService(database *db.Database) VaultService {
	ni := VaultService{
		Service: services.Service[*resources.VaultArgs, *resources.VaultResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
