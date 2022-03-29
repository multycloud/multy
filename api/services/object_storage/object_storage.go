package object_storage

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type ObjectStorageService struct {
	Service services.Service[*resources.ObjectStorageArgs, *resources.ObjectStorageResource]
}

func (s ObjectStorageService) Convert(resourceId string, args *resources.ObjectStorageArgs, state *output.TfState) (*resources.ObjectStorageResource, error) {
	return &resources.ObjectStorageResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Versioning:       args.Versioning,
	}, nil
}

func NewObjectStorageService(database *db.Database) ObjectStorageService {
	nsg := ObjectStorageService{
		Service: services.Service[*resources.ObjectStorageArgs, *resources.ObjectStorageResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
