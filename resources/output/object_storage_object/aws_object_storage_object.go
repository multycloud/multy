package object_storage_object

import "multy-go/resources/common"

// aws_s3_bucket_object
type AwsS3BucketObject struct {
	common.AwsResource `hcl:",squash"`
	Bucket             string `hcl:"bucket,expr"`
	Key                string `hcl:"key"`
	Acl                string `hcl:"acl"`
	Content            string `hcl:"content"`
	ContentType        string `hcl:"content_type"`
}
