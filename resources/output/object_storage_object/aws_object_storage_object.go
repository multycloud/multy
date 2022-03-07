package object_storage_object

import (
	"multy-go/resources/common"
)

// aws_s3_bucket_object
type AwsS3BucketObject struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_bucket_object"`
	Bucket              string `hcl:"bucket,expr"`
	Key                 string `hcl:"key"`
	Acl                 string `hcl:"acl" hcle:"omitempty"`
	Content             string `hcl:"content"  hcle:"omitempty"`
	ContentType         string `hcl:"content_type" hcle:"omitempty"`
	Source              string `hcl:"source" hcle:"omitempty"`
}
