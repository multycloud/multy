package object_storage

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type ObjectStorageService struct {
	Service services.Service[*resourcespb.ObjectStorageArgs, *resourcespb.ObjectStorageResource]
}

func (s ObjectStorageService) Convert(resourceId string, args *resourcespb.ObjectStorageArgs, state *output.TfState) (*resourcespb.ObjectStorageResource, error) {
	return &resourcespb.ObjectStorageResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Versioning:       args.Versioning,
	}, nil
}

func NewObjectStorageService(database *db.Database) ObjectStorageService {
	nsg := ObjectStorageService{
		Service: services.Service[*resourcespb.ObjectStorageArgs, *resourcespb.ObjectStorageResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "object_storage",
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
