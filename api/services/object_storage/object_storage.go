package object_storage

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type ObjectStorageService struct {
	Service services.Service[*resources.CloudSpecificObjectStorageArgs, *resources.ObjectStorageResource]
}

func (s ObjectStorageService) Convert(resourceId string, args []*resources.CloudSpecificObjectStorageArgs) *resources.ObjectStorageResource {
	var result []*resources.CloudSpecificObjectStorageResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificObjectStorageResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
		})
	}

	return &resources.ObjectStorageResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s ObjectStorageService) NewArg() *resources.CloudSpecificObjectStorageArgs {
	return &resources.CloudSpecificObjectStorageArgs{}
}

func (s ObjectStorageService) Nil() *resources.ObjectStorageResource {
	return nil
}

func NewObjectStorageServiceService(database *db.Database) ObjectStorageService {
	nsg := ObjectStorageService{
		Service: services.Service[*resources.CloudSpecificObjectStorageArgs, *resources.ObjectStorageResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
