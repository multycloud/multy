package object_storage_object

import (
	"github.com/multycloud/multy/resources/common"
)

type AwsS3BucketObject struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_object"`
	Bucket              string `hcl:"bucket,expr" json:"bucket"`
	Key                 string `hcl:"key" json:"key"`
	Acl                 string `hcl:"acl" hcle:"omitempty" json:"acl"`
	ContentBase64       string `hcl:"content_base64"  hcle:"omitempty" json:"content_base64"`
	ContentType         string `hcl:"content_type" hcle:"omitempty" json:"content_type"`
	Source              string `hcl:"source" hcle:"omitempty" json:"source"`
}
