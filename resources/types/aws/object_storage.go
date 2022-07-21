package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/types"
)

type AwsObjectStorage struct {
	*types.ObjectStorage
}

func InitObjectStorage(vn *types.ObjectStorage) resources.ResourceTranslator[*resourcespb.ObjectStorageResource] {
	return AwsObjectStorage{vn}
}

func (r AwsObjectStorage) FromState(state *output.TfState) (*resourcespb.ObjectStorageResource, error) {
	out := &resourcespb.ObjectStorageResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Versioning:  r.Args.Versioning,
		GcpOverride: r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	stateResource, err := output.GetParsedById[object_storage.AwsS3Bucket](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	out.AwsOutputs = &resourcespb.ObjectStorageAwsOutputs{S3BucketArn: stateResource.Arn}

	return out, nil
}

func (r AwsObjectStorage) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var awsResources []output.TfBlock
	s3Bucket := object_storage.AwsS3Bucket{
		AwsResource: &common.AwsResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
		},
		Bucket: r.Args.Name,
	}
	awsResources = append(awsResources, s3Bucket)

	if r.Args.Versioning {
		awsResources = append(awsResources, object_storage.AwsS3BucketVersioning{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			BucketId:                s3Bucket.GetBucketId(),
			VersioningConfiguration: object_storage.VersioningConfiguration{Status: "Enabled"},
		})
	}
	return awsResources, nil
}

func (r AwsObjectStorage) GetMainResourceName() (string, error) {
	return "aws_s3_bucket", nil
}
