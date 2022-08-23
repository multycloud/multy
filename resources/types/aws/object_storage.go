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

func (r AwsObjectStorage) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.ObjectStorageResource, error) {
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

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage.AwsS3Bucket](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.Name = stateResource.Bucket
		out.AwsOutputs = &resourcespb.ObjectStorageAwsOutputs{S3BucketArn: stateResource.Arn}
		output.AddToStatuses(statuses, "aws_s3_bucket", output.MaybeGetPlannedChageById[object_storage.AwsS3Bucket](plan, r.ResourceId))
	} else {
		statuses["aws_s3_bucket"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage.AwsS3BucketVersioning](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.Versioning = len(stateResource.VersioningConfiguration) > 0 && stateResource.VersioningConfiguration[0].Status == "Enabled"
		output.AddToStatuses(statuses, "aws_s3_bucket_versioning", output.MaybeGetPlannedChageById[object_storage.AwsS3BucketVersioning](plan, r.ResourceId))
	} else {
		out.Versioning = false
		if r.Args.Versioning {
			statuses["aws_s3_bucket_versioning"] = commonpb.ResourceStatus_NEEDS_CREATE
		}
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
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
			VersioningConfiguration: []object_storage.VersioningConfiguration{{Status: "Enabled"}},
		})
	}
	return awsResources, nil
}

func (r AwsObjectStorage) GetMainResourceName() (string, error) {
	return "aws_s3_bucket", nil
}
