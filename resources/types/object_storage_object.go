package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

// AWS: aws_s3_bucket_object
// Azure: azurerm_storage_blob

var SUPPORTED_CONTENT_TYPES = []string{"text/html", "application/zip"}

type ObjectStorageObject struct {
	*resources.CommonResourceParams
	Name          string         `hcl:"name"`
	Content       string         `hcl:"content,optional"`
	ObjectStorage *ObjectStorage `mhcl:"ref=object_storage"`
	ContentType   string         `hcl:"content_type,optional"`
	Acl           string         `hcl:"acl,optional"`
	Source        string         `hcl:"source,optional"`
}

func (r *ObjectStorageObject) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	var acl string
	if r.Acl == "public_read" {
		acl = "public-read"
	} else {
		acl = "private"
	}
	if cloud == common.AWS {
		return []output.TfBlock{object_storage_object.AwsS3BucketObject{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
			},
			Bucket:      r.ObjectStorage.GetResourceName(cloud),
			Key:         r.Name,
			Acl:         acl,
			Content:     r.Content,
			ContentType: r.ContentType,
			Source:      r.Source,
		}}, nil
	} else if cloud == common.AZURE {
		var containerName string
		if r.Acl == "public_read" {
			containerName = r.ObjectStorage.GetAssociatedPublicContainerResourceName(cloud)
		} else {
			containerName = r.ObjectStorage.GetAssociatedPrivateContainerResourceName(cloud)
		}
		return []output.TfBlock{
			object_storage_object.AzureStorageAccountBlob{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
					Name:              r.Name,
				},
				StorageAccountName:   r.ObjectStorage.GetResourceName(cloud),
				StorageContainerName: containerName,
				Type:                 "Block",
				SourceContent:        r.Content,
				ContentType:          r.ContentType,
				Source:               r.Source,
			}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *ObjectStorageObject) GetS3Key() string {
	return fmt.Sprintf("%s.%s.key", "aws_s3_bucket_object", r.GetTfResourceId(common.AWS))
}

func (r *ObjectStorageObject) GetAzureBlobName() string {
	return fmt.Sprintf("%s.%s.name", "azurerm_storage_blob", r.GetTfResourceId(common.AZURE))
}

func (r *ObjectStorageObject) GetAzureBlobUrl() string {
	return fmt.Sprintf("%s.%s.url", "azurerm_storage_blob", r.GetTfResourceId(common.AZURE))
}

func (r *ObjectStorageObject) IsPrivate() bool {
	return r.Acl == "private"
}

func (r *ObjectStorageObject) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, r.CommonResourceParams.Validate(ctx, cloud)...)
	if len(r.Content) > 0 && len(r.Source) > 0 {
		errs = append(errs, r.NewError("content", "content can't be set if source is already set"))
	}
	if len(r.Content) == 0 && len(r.Source) == 0 {
		errs = append(errs, r.NewError("", "content or source must be set"))
	}
	if len(r.Content) > 0 {
		if !util.Contains(SUPPORTED_CONTENT_TYPES, r.ContentType) {
			errs = append(errs, r.NewError("content_type", fmt.Sprintf("%s not a valid content_type", r.ContentType)))
		}
	}
	if r.Acl != "" && r.Acl != "public_read" && r.Acl != "private" {
		errs = append(errs, r.NewError("acl", fmt.Sprintf("%s not a valid acl", r.Acl)))
	}
	return errs
}

func (r *ObjectStorageObject) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return "aws_s3_bucket_object", nil
	case common.AZURE:
		return "azurerm_storage_blob", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}

func (r *ObjectStorageObject) GetLocation(cloud common.CloudProvider, ctx resources.MultyContext) string {
	return r.ObjectStorage.GetLocation(cloud, ctx)
}
