package vault

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultService struct {
	Service services.Service[*resourcespb.VaultArgs, *resourcespb.VaultResource]
}

func (s VaultService) Convert(resourceId string, args *resourcespb.VaultArgs, state *output.TfState) (*resourcespb.VaultResource, error) {
	return &resourcespb.VaultResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
	}, nil
}

func NewVaultService(database *db.Database) VaultService {
	ni := VaultService{
		Service: services.Service[*resourcespb.VaultArgs, *resourcespb.VaultResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
