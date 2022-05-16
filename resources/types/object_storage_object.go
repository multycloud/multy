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
	"github.com/multycloud/multy/resources/output/terraform"
	"github.com/multycloud/multy/validate"
)

// AWS: aws_s3_object
// Azure: azurerm_storage_blob

var objectStorageObjectMetadata = resources.ResourceMetadata[*resourcespb.ObjectStorageObjectArgs, *ObjectStorageObject, *resourcespb.ObjectStorageObjectResource]{
	CreateFunc:        CreateObjectStorageObject,
	UpdateFunc:        UpdateObjectStorageObject,
	ReadFromStateFunc: ObjectStorageObjectFromState,
	ExportFunc: func(r *ObjectStorageObject, _ *resources.Resources) (*resourcespb.ObjectStorageObjectArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewObjectStorageObject,
	AbbreviatedName: "st",
}

type ObjectStorageObject struct {
	resources.ChildResourceWithId[*ObjectStorage, *resourcespb.ObjectStorageObjectArgs]

	ObjectStorage *ObjectStorage `mhcl:"ref=object_storage"`
}

func (r *ObjectStorageObject) GetMetadata() resources.ResourceMetadataInterface {
	return &objectStorageObjectMetadata
}

func CreateObjectStorageObject(resourceId string, args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) (*ObjectStorageObject, error) {
	return NewObjectStorageObject(resourceId, args, others)
}

func UpdateObjectStorageObject(resource *ObjectStorageObject, vn *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func ObjectStorageObjectFromState(resource *ObjectStorageObject, _ *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
	return &resourcespb.ObjectStorageObjectResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  resource.ResourceId,
			NeedsUpdate: false,
		},
		Name:            resource.Args.Name,
		Acl:             resource.Args.Acl,
		ObjectStorageId: resource.Args.ObjectStorageId,
		ContentBase64:   resource.Args.ContentBase64,
		ContentType:     resource.Args.ContentType,
		Source:          resource.Args.Source,
	}, nil
}

func NewObjectStorageObject(resourceId string, args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) (*ObjectStorageObject, error) {
	o := &ObjectStorageObject{
		ChildResourceWithId: resources.ChildResourceWithId[*ObjectStorage, *resourcespb.ObjectStorageObjectArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	obj, err := resources.Get[*ObjectStorage](resourceId, others, args.ObjectStorageId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{o.NewValidationError(err, "object_storage_id")})
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
	out.ContentBase64 = r.Args.ContentBase64
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
			Bucket:        r.ObjectStorage.GetResourceName(),
			Key:           r.Args.Name,
			Acl:           acl,
			ContentBase64: r.Args.ContentBase64,
			ContentType:   r.Args.ContentType,
			Source:        r.Args.Source,
		}}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		var containerName string
		if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
			containerName = r.ObjectStorage.GetAssociatedPublicContainerResourceName()
		} else {
			containerName = r.ObjectStorage.GetAssociatedPrivateContainerResourceName()
		}
		contentFile := terraform.NewLocalFile(r.ResourceId, r.Args.ContentBase64)
		return []output.TfBlock{
			contentFile,
			object_storage_object.AzureStorageAccountBlob{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					Name:              r.Args.Name,
				},
				StorageAccountName:   r.ObjectStorage.GetResourceName(),
				StorageContainerName: containerName,
				Type:                 "Block",
				Source:               contentFile.GetFilename(),
				ContentType:          r.Args.ContentType,
			}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *ObjectStorageObject) GetS3Key() string {
	return fmt.Sprintf("%s.%s.key", "aws_s3_object", r.ResourceId)
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
	if len(r.Args.ContentBase64) == 0 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("content_base64 must be set"), ""))
	}
	return errs
}

func (r *ObjectStorageObject) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return "aws_s3_object", nil
	case commonpb.CloudProvider_AZURE:
		return "azurerm_storage_blob", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *ObjectStorageObject) GetCloudSpecificLocation() string {
	return r.ObjectStorage.GetCloudSpecificLocation()
}
