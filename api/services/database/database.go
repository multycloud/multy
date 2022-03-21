package database

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type DatabaseService struct {
	Service services.Service[*resources.CloudSpecificDatabaseArgs, *resources.DatabaseResource]
}

func (s DatabaseService) Convert(resourceId string, args []*resources.CloudSpecificDatabaseArgs) *resources.DatabaseResource {
	var result []*resources.CloudSpecificDatabaseResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificDatabaseResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			Engine:           r.Engine,
			EngineVersion:    r.EngineVersion,
			StorageMb:        r.StorageMb,
			Size:             r.Size,
			Username:         r.Username,
			Password:         r.Password,
			SubnetIds:        r.SubnetIds,
		})
	}

	return &resources.DatabaseResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s DatabaseService) NewArg() *resources.CloudSpecificDatabaseArgs {
	return &resources.CloudSpecificDatabaseArgs{}
}

func (s DatabaseService) Nil() *resources.DatabaseResource {
	return nil
}

func NewDatabaseServiceService(database *db.Database) DatabaseService {
	ni := DatabaseService{
		Service: services.Service[*resources.CloudSpecificDatabaseArgs, *resources.DatabaseResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
