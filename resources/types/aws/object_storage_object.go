package aws_resources

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

type AwsObjectStorageObject struct {
	*types.ObjectStorageObject
}

func InitObjectStorageObject(vn *types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource] {
	return AwsObjectStorageObject{vn}
}

func (r AwsObjectStorageObject) FromState(state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
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
		stateResource, err := output.GetParsed[object_storage.AwsS3Bucket](state, id)
		if err != nil {
			return nil, err
		}
		out.Url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", stateResource.Bucket, r.Args.Name)
		out.AwsOutputs = &resourcespb.ObjectStorageObjectAwsOutputs{S3BucketObjectId: stateResource.ResourceId}
	} else {
		out.Url = "dryrun"
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

func (r AwsObjectStorageObject) GetMainResourceName() (string, error) {
	return "aws_s3_object", nil
}
