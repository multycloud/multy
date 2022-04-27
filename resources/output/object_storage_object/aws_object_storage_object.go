package object_storage_object

import (
	"github.com/multycloud/multy/resources/common"
)

type AwsS3BucketObject struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_object"`
	Bucket              string `hcl:"bucket,expr"`
	Key                 string `hcl:"key"`
	Acl                 string `hcl:"acl" hcle:"omitempty"`
	ContentBase64       string `hcl:"content_base64"  hcle:"omitempty"`
	ContentType         string `hcl:"content_type" hcle:"omitempty"`
	Source              string `hcl:"source" hcle:"omitempty"`
}
