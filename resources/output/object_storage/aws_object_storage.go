package object_storage

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

// aws_s3_bucket
type AwsS3Bucket struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_bucket"`
	Bucket              string `hcl:"bucket" json:"bucket"`

	Arn string `json:"arn" hcle:"omitempty"`
}

type AwsS3BucketVersioning struct {
	*common.AwsResource     `hcl:",squash" default:"name=aws_s3_bucket_versioning" json:"*_common_._aws_resource"`
	BucketId                string                    `hcl:"bucket,expr" json:"bucket_id"`
	VersioningConfiguration []VersioningConfiguration `hcl:"versioning_configuration,blocks" json:"versioning_configuration"`
}

type VersioningConfiguration struct {
	Status string `hcl:"status" json:"status"`
}

func (vpc *AwsS3Bucket) GetBucketId() string {
	return fmt.Sprintf("aws_s3_bucket.%s.id", vpc.ResourceId)
}
