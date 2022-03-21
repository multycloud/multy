package vault

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type VaultSecretService struct {
	Service services.Service[*resources.CloudSpecificVaultSecretArgs, *resources.VaultSecretResource]
}

func (s VaultSecretService) Convert(resourceId string, args []*resources.CloudSpecificVaultSecretArgs) *resources.VaultSecretResource {
	var result []*resources.CloudSpecificVaultSecretResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificVaultSecretResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			Value:            r.Value,
			VaultId:          r.VaultId,
		})
	}

	return &resources.VaultSecretResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s VaultSecretService) NewArg() *resources.CloudSpecificVaultSecretArgs {
	return &resources.CloudSpecificVaultSecretArgs{}
}

func (s VaultSecretService) Nil() *resources.VaultSecretResource {
	return nil
}

func NewVaultSecretService(database *db.Database) VaultSecretService {
	ni := VaultSecretService{
		Service: services.Service[*resources.CloudSpecificVaultSecretArgs, *resources.VaultSecretResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
