package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/object_storage_object"
	"multy-go/validate"
)

// AWS: aws_s3_bucket_object
// Azure: azurerm_storage_blob

type ObjectStorageObject struct {
	*resources.CommonResourceParams
	Name              string `hcl:"name"`
	Content           string `hcl:"content"`
	ObjectStorageName string `hcl:"object_storage_name"`
	ContentType       string `hcl:"content_type"`
	Acl               string `hcl:"acl,optional"`
}

func (r *ObjectStorageObject) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	var objectStorage *ObjectStorage
	if o, err := ctx.GetResource(r.ObjectStorageName); err != nil {
		r.LogFatal(r.ResourceId, "object_storage_name", err.Error())
	} else {
		objectStorage = o.Resource.(*ObjectStorage)
	}
	var acl string
	if r.Acl == "public_read" {
		acl = "public-read"
	}
	if cloud == common.AWS {
		return []interface{}{object_storage_object.AwsS3BucketObject{
			AwsResource: common.AwsResource{
				ResourceName: "aws_s3_bucket_object",
				ResourceId:   r.GetTfResourceId(cloud),
			},
			Bucket:      objectStorage.GetResourceName(cloud),
			Key:         r.Name,
			Acl:         acl,
			Content:     r.Content,
			ContentType: r.ContentType,
		}}
	} else if cloud == common.AZURE {
		var containerName string
		if r.Acl == "public_read" {
			containerName = objectStorage.GetAssociatedPublicContainerResourceName(cloud)
		} else {
			containerName = objectStorage.GetAssociatedPrivateContainerResourceName(cloud)
		}
		return []interface{}{
			object_storage_object.AzureStorageAccountBlob{
				AzResource: common.AzResource{
					ResourceName: "azurerm_storage_blob",
					ResourceId:   r.GetTfResourceId(cloud),
					Name:         r.Name,
				},
				StorageAccountName:   objectStorage.GetResourceName(cloud),
				StorageContainerName: containerName,
				Type:                 "Block",
				SourceContent:        r.Content,
				ContentType:          r.ContentType,
			}}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *ObjectStorageObject) Validate(ctx resources.MultyContext) {
	if r.ContentType != "text/html" {
		r.LogFatal(r.ResourceId, "content_type", fmt.Sprintf("%s not a valid content_type", r.ContentType))
	}
	if r.Acl != "" && r.Acl != "public_read" {
		r.LogFatal(r.ResourceId, "content_type", fmt.Sprintf("%s not a valid acl", r.Acl))
	}
	return
}
