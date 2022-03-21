package object_storage_object

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type ObjectStorageObjectService struct {
	Service services.Service[*resources.CloudSpecificObjectStorageObjectArgs, *resources.ObjectStorageObjectResource]
}

func (s ObjectStorageObjectService) Convert(resourceId string, args []*resources.CloudSpecificObjectStorageObjectArgs) *resources.ObjectStorageObjectResource {
	var result []*resources.CloudSpecificObjectStorageObjectResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificObjectStorageObjectResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			Acl:              r.Acl,
			ObjectStorageId:  r.ObjectStorageId,
			Content:          r.Content,
			ContentType:      r.ContentType,
			Source:           r.Source,
		})
	}

	return &resources.ObjectStorageObjectResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s ObjectStorageObjectService) NewArg() *resources.CloudSpecificObjectStorageObjectArgs {
	return &resources.CloudSpecificObjectStorageObjectArgs{}
}

func (s ObjectStorageObjectService) Nil() *resources.ObjectStorageObjectResource {
	return nil
}

func NewObjectStorageObjectService(database *db.Database) ObjectStorageObjectService {
	nsg := ObjectStorageObjectService{
		Service: services.Service[*resources.CloudSpecificObjectStorageObjectArgs, *resources.ObjectStorageObjectResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
