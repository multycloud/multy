package object_storage

import "multy-go/resources/common"

// aws_s3_bucket
type AwsS3Bucket struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_bucket"`
	Bucket              string `hcl:"bucket"`
}
