package object_storage_object

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage_object"
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
		Content:          args.Content,
		ContentType:      args.ContentType,
		Source:           args.Source,
	}, nil
}

func NewObjectStorageObjectService(database *db.Database) ObjectStorageObjectService {
	nsg := ObjectStorageObjectService{
		Service: services.Service[*resourcespb.ObjectStorageObjectArgs, *resourcespb.ObjectStorageObjectResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}

func getUrl(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(object_storage_object.AwsS3BucketObject{}, resourceId)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", values["bucket"].(string), values["name"].(string)), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(object_storage_object.AzureStorageAccountBlob{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["url"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}
