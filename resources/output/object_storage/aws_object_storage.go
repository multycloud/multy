package object_storage

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

// aws_s3_bucket
type AwsS3Bucket struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_bucket"`
	Bucket              string `hcl:"bucket"`
}

type AwsS3BucketVersioning struct {
	*common.AwsResource     `hcl:",squash" default:"name=aws_s3_bucket_versioning"`
	BucketId                string                  `hcl:"bucket,expr"`
	VersioningConfiguration VersioningConfiguration `hcl:"versioning_configuration"`
}

type VersioningConfiguration struct {
	Status string `hcl:"status"`
}

func (vpc *AwsS3Bucket) GetBucketId() string {
	return fmt.Sprintf("aws_s3_bucket.%s.id", vpc.ResourceId)
}
