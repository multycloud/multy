package object_storage_object

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
)

type GoogleStorageBucketObject struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_storage_bucket_object"`
	Bucket              string `hcl:"bucket,expr" hcl:"Bucket" json:"bucket"`
	Content             string `hcl:"content,expr" json:"content"`
	ContentType         string `hcl:"content_type"  hcle:"omitempty" json:"content_type"`
}

type GoogleStorageObjectAccessControl struct {
	*output.TerraformResource `hcl:",squash"  default:"name=google_storage_object_access_control"`
	Object                    string `hcl:"object,expr" json:"object"`
	Bucket                    string `hcl:"bucket,expr" json:"bucket"`
	Role                      string `hcl:"role" json:"role"`
	Entity                    string `hcl:"entity" json:"entity"`
}
