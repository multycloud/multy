package object_storage_object

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type ObjectStorageObjectService struct {
	Service services.Service[*resourcespb.ObjectStorageObjectArgs, *resourcespb.ObjectStorageObjectResource]
}

func (s ObjectStorageObjectService) Convert(resourceId string, args *resourcespb.ObjectStorageObjectArgs, state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
	return &resourcespb.ObjectStorageObjectResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Acl:              args.Acl,
		ObjectStorageId:  args.ObjectStorageId,
		ContentBase64:    args.ContentBase64,
		ContentType:      args.ContentType,
		Source:           args.Source,
	}, nil
}

func NewObjectStorageObjectService(database *db.Database) ObjectStorageObjectService {
	nsg := ObjectStorageObjectService{
		Service: services.Service[*resourcespb.ObjectStorageObjectArgs, *resourcespb.ObjectStorageObjectResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "object_storage_object",
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
