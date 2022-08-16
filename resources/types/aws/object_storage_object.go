package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/resources/types"
)

type AwsObjectStorageObject struct {
	*types.ObjectStorageObject
}

func InitObjectStorageObject(vn *types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource] {
	return AwsObjectStorageObject{vn}
}

func (r AwsObjectStorageObject) FromState(state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
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

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage_object.AwsS3BucketObject](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Key
		out.Acl = parseObjectAcl(stateResource.Acl)
		out.ContentBase64 = stateResource.ContentBase64
		out.ContentType = stateResource.ContentType
		out.Source = stateResource.Source
		if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
			out.Url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", stateResource.Bucket, stateResource.Key)
		}
		out.AwsOutputs = &resourcespb.ObjectStorageObjectAwsOutputs{S3BucketObjectId: stateResource.ResourceId}
	} else {
		statuses["aws_s3_bucket_object"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AwsObjectStorageObject) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var acl string
	if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
		acl = "public-read"
	} else {
		acl = "private"
	}

	bucketId, err := resources.GetMainOutputId(AwsObjectStorage{r.ObjectStorage})
	if err != nil {
		return nil, err
	}

	return []output.TfBlock{object_storage_object.AwsS3BucketObject{
		AwsResource: &common.AwsResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
		},
		Bucket:        bucketId,
		Key:           r.Args.Name,
		Acl:           acl,
		ContentBase64: r.Args.ContentBase64,
		ContentType:   r.Args.ContentType,
		Source:        r.Args.Source,
	}}, nil
}

func parseObjectAcl(acl string) resourcespb.ObjectStorageObjectAcl {
	switch acl {
	case "public-read":
		return resourcespb.ObjectStorageObjectAcl_PUBLIC_READ
	default:
		return resourcespb.ObjectStorageObjectAcl_PRIVATE
	}
}

func (r AwsObjectStorageObject) GetMainResourceName() (string, error) {
	return "aws_s3_object", nil
}
