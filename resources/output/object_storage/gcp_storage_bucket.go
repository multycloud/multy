package object_storage

import "github.com/multycloud/multy/resources/common"

type GoogleStorageBucket struct {
	*common.GcpResource      `hcl:",squash"  default:"name=google_storage_bucket"`
	UniformBucketLevelAccess bool                            `hcl:"uniform_bucket_level_access"`
	Versioning               []GoogleStorageBucketVersioning `hcl:"versioning,blocks" hcle:"omitempty"`
	Location                 string                          `hcl:"location"`
}

type GoogleStorageBucketVersioning struct {
	Enabled bool `hcl:"enabled"`
}
