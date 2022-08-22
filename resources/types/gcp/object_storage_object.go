package gcp_resources

import (
	"encoding/base64"
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

func (r GcpObjectStorageObject) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.ObjectStorageObjectResource, error) {
	out := &resourcespb.ObjectStorageObjectResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:            r.Args.Name,
		Acl:             r.Args.Acl,
		ObjectStorageId: r.Args.ObjectStorageId,
		ContentBase64:   r.Args.ContentBase64,
		ContentType:     r.Args.ContentType,
		Source:          r.Args.Source,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage_object.GoogleStorageBucketObject](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.ContentType = stateResource.ContentType
		out.ContentBase64 = base64.StdEncoding.EncodeToString([]byte(stateResource.Content))
		out.Url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", stateResource.Bucket, r.Args.Name)
		out.GcpOutputs = &resourcespb.ObjectStorageObjectGcpOutputs{StorageBucketObjectId: stateResource.SelfLink}
		if aclResource, exists, err := output.MaybeGetParsedById[object_storage_object.GoogleStorageObjectAccessControl](state, r.ResourceId); exists {
			if err != nil {
				return nil, err
			}
			out.GcpOutputs.StorageObjectAccessControl = aclResource.ResourceId
			out.Acl = resourcespb.ObjectStorageObjectAcl_PUBLIC_READ
		} else {
			out.Acl = resourcespb.ObjectStorageObjectAcl_PRIVATE
		}
	} else {
		statuses["azure_storage_account_blob"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
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
