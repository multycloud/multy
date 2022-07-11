package object_storage_object

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
)

type GoogleStorageBucketObject struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_storage_bucket_object"`
	Bucket              string `hcl:"bucket,expr" hcl:"Bucket"`
	Content             string `hcl:"content,expr"`
	ContentType         string `hcl:"content_type"  hcle:"omitempty"`
}

type GoogleStorageObjectAccessControl struct {
	*output.TerraformResource `hcl:",squash"  default:"name=google_storage_object_access_control"`
	Object                    string `hcl:"object,expr"`
	Bucket                    string `hcl:"bucket,expr"`
	Role                      string `hcl:"role"`
	Entity                    string `hcl:"entity"`
}
