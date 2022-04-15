package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

// AWS: aws_s3_bucket_object
// Azure: azurerm_storage_blob

var SUPPORTED_CONTENT_TYPES = []string{"text/html", "application/zip"}

type ObjectStorageObject struct {
	resources.ChildResourceWithId[*ObjectStorage, *resourcespb.ObjectStorageObjectArgs]

	ObjectStorage *ObjectStorage `mhcl:"ref=object_storage"`
}

func NewObjectStorageObject(resourceId string, args *resourcespb.ObjectStorageObjectArgs, others resources.Resources) (*ObjectStorageObject, error) {
	o := &ObjectStorageObject{
		ChildResourceWithId: resources.ChildResourceWithId[*ObjectStorage, *resourcespb.ObjectStorageObjectArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	obj, err := resources.Get[*ObjectStorage](resourceId, others, args.ObjectStorageId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{o.NewValidationError(err.Error(), "object_storage_id")})
	}
	o.Parent = obj
	o.ObjectStorage = obj
	return o, nil
}

func (r *ObjectStorageObject) FromState(state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
	out := new(resourcespb.ObjectStorageObjectResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}

	id, err := resources.GetMainOutputRef(r.Parent)
	if err != nil {
		return nil, err
	}

	out.Name = r.Args.Name
	out.Content = r.Args.Content
	out.ContentType = r.Args.ContentType
	out.ObjectStorageId = r.Args.ObjectStorageId
	out.Acl = r.Args.Acl
	out.Source = r.Args.Source

	switch r.GetCloud() {
	case common.AWS:
		stateResource, err := output.GetParsed[object_storage.AwsS3Bucket](state, id)
		if err != nil {
			return nil, err
		}
		out.Url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", stateResource.Bucket, r.Args.Name)
	case common.AZURE:
		stateResource, err := output.GetParsed[object_storage.AzureStorageAccount](state, id)
		if err != nil {
			return nil, err
		}
		out.Url = fmt.Sprintf("https://%s.blob.core.windows.net/public/%s", stateResource.AzResource.Name, r.Args.Name)
	}

	return out, nil
}

func (r *ObjectStorageObject) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var acl string
	if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
		acl = "public-read"
	} else {
		acl = "private"
	}
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{object_storage_object.AwsS3BucketObject{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			Bucket:      r.ObjectStorage.GetResourceName(),
			Key:         r.Args.Name,
			Acl:         acl,
			Content:     r.Args.Content,
			ContentType: r.Args.ContentType,
			Source:      r.Args.Source,
		}}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		var containerName string
		if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
			containerName = r.ObjectStorage.GetAssociatedPublicContainerResourceName()
		} else {
			containerName = r.ObjectStorage.GetAssociatedPrivateContainerResourceName()
		}
		return []output.TfBlock{
			object_storage_object.AzureStorageAccountBlob{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					Name:              r.Args.Name,
				},
				StorageAccountName:   r.ObjectStorage.GetResourceName(),
				StorageContainerName: containerName,
				Type:                 "Block",
				SourceContent:        r.Args.Content,
				ContentType:          r.Args.ContentType,
				Source:               r.Args.Source,
			}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *ObjectStorageObject) GetS3Key() string {
	return fmt.Sprintf("%s.%s.key", "aws_s3_bucket_object", r.ResourceId)
}

func (r *ObjectStorageObject) GetAzureBlobName() string {
	return fmt.Sprintf("%s.%s.name", "azurerm_storage_blob", r.ResourceId)
}

func (r *ObjectStorageObject) GetAzureBlobUrl() string {
	return fmt.Sprintf("%s.%s.url", "azurerm_storage_blob", r.ResourceId)
}

func (r *ObjectStorageObject) IsPrivate() bool {
	return r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PRIVATE
}

func (r *ObjectStorageObject) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if len(r.Args.Content) == 0 {
		errs = append(errs, r.NewValidationError("content must be set", ""))
	}
	if !util.Contains(SUPPORTED_CONTENT_TYPES, r.Args.ContentType) {
		errs = append(errs, r.NewValidationError(fmt.Sprintf("%s not a valid content_type", r.Args.ContentType), "content_type"))
	}
	return errs
}

func (r *ObjectStorageObject) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return "aws_s3_bucket_object", nil
	case commonpb.CloudProvider_AZURE:
		return "azurerm_storage_blob", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *ObjectStorageObject) GetCloudSpecificLocation() string {
	return r.ObjectStorage.GetCloudSpecificLocation()
}
