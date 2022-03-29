package object_storage_object

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type ObjectStorageObjectService struct {
	Service services.Service[*resources.ObjectStorageObjectArgs, *resources.ObjectStorageObjectResource]
}

func (s ObjectStorageObjectService) Convert(resourceId string, args *resources.ObjectStorageObjectArgs, state *output.TfState) (*resources.ObjectStorageObjectResource, error) {
	return &resources.ObjectStorageObjectResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Acl:              args.Acl,
		ObjectStorageId:  args.ObjectStorageId,
		Content:          args.Content,
		ContentType:      args.ContentType,
		Source:           args.Source,
	}, nil
}

func NewObjectStorageObjectService(database *db.Database) ObjectStorageObjectService {
	nsg := ObjectStorageObjectService{
		Service: services.Service[*resources.ObjectStorageObjectArgs, *resources.ObjectStorageObjectResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
