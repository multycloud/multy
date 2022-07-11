package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/resources/types"
)

type GcpObjectStorageObject struct {
	*types.ObjectStorageObject
}

func InitObjectStorageObject(vn *types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource] {
	return GcpObjectStorageObject{vn}
}

func (r GcpObjectStorageObject) FromState(state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
	out := new(resourcespb.ObjectStorageObjectResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}

	id, err := resources.GetMainOutputRef(r)
	if err != nil {
		return nil, err
	}

	out.Name = r.Args.Name
	out.ContentBase64 = r.Args.ContentBase64
	out.ContentType = r.Args.ContentType
	out.ObjectStorageId = r.Args.ObjectStorageId
	out.Acl = r.Args.Acl
	out.Source = r.Args.Source

	if !flags.DryRun {
		stateResource, err := output.GetParsed[object_storage_object.GoogleStorageBucketObject](state, id)
		if err != nil {
			return nil, err
		}
		out.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", stateResource.Bucket, r.Args.Name)
	} else {
		out.Url = "dryrun"
	}

	return out, nil
}

func (r GcpObjectStorageObject) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var result []output.TfBlock
	bucketName := fmt.Sprintf("%s.%s.name", output.GetResourceName(object_storage.GoogleStorageBucket{}), r.ObjectStorage.ResourceId)
	result = append(result, &object_storage_object.GoogleStorageBucketObject{
		GcpResource: common.NewGcpResourceWithNoProject(r.ResourceId, r.Args.Name),
		Bucket:      bucketName,
		Content:     fmt.Sprintf("base64decode(\"%s\")", r.Args.ContentBase64),
		ContentType: r.Args.ContentType,
	})

	if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
		objectName := fmt.Sprintf("%s.%s.output_name", output.GetResourceName(object_storage_object.GoogleStorageBucketObject{}), r.ResourceId)
		result = append(result, &object_storage_object.GoogleStorageObjectAccessControl{
			TerraformResource: &output.TerraformResource{
				ResourceId: r.ResourceId,
			},
			Object: objectName,
			Bucket: bucketName,
			Role:   "READER",
			Entity: "allUsers",
		})
	}
	return result, nil

}

func (r GcpObjectStorageObject) GetMainResourceName() (string, error) {
	return output.GetResourceName(object_storage_object.GoogleStorageBucketObject{}), nil
}
