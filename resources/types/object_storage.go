package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/validate"
)

var objectStorageMetadata = resources.ResourceMetadata[*resourcespb.ObjectStorageArgs, *ObjectStorage, *resourcespb.ObjectStorageResource]{
	CreateFunc:        CreateObjectStorage,
	UpdateFunc:        UpdateObjectStorage,
	ReadFromStateFunc: ObjectStorageFromState,
	ExportFunc: func(r *ObjectStorage, _ *resources.Resources) (*resourcespb.ObjectStorageArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewObjectStorage,
	AbbreviatedName: "st",
}

type ObjectStorage struct {
	resources.ResourceWithId[*resourcespb.ObjectStorageArgs]
}

func (r *ObjectStorage) GetMetadata() resources.ResourceMetadataInterface {
	return &objectStorageMetadata
}

func NewObjectStorage(resourceId string, db *resourcespb.ObjectStorageArgs, _ *resources.Resources) (*ObjectStorage, error) {
	return &ObjectStorage{
		ResourceWithId: resources.ResourceWithId[*resourcespb.ObjectStorageArgs]{
			ResourceId: resourceId,
			Args:       db,
		},
	}, nil
}

func CreateObjectStorage(resourceId string, args *resourcespb.ObjectStorageArgs, others *resources.Resources) (*ObjectStorage, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		rgId, err := NewRg("st", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return nil, err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewObjectStorage(resourceId, args, others)
}

func UpdateObjectStorage(resource *ObjectStorage, vn *resourcespb.ObjectStorageArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

type AclRules struct{}

func ObjectStorageFromState(resource *ObjectStorage, _ *output.TfState) (*resourcespb.ObjectStorageResource, error) {
	return &resourcespb.ObjectStorageResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:       resource.Args.Name,
		Versioning: resource.Args.Versioning,
	}, nil
}

func (r *ObjectStorage) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var awsResources []output.TfBlock
		s3Bucket := object_storage.AwsS3Bucket{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			Bucket: r.Args.Name,
		}
		awsResources = append(awsResources, s3Bucket)

		if r.Args.Versioning {
			awsResources = append(awsResources, object_storage.AwsS3BucketVersioning{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				},
				BucketId:                s3Bucket.GetBucketId(),
				VersioningConfiguration: object_storage.VersioningConfiguration{Status: "Enabled"},
			})
		}
		return awsResources, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)

		storageAccount := object_storage.AzureStorageAccount{
			AzResource: common.NewAzResource(
				r.ResourceId, common.RemoveSpecialChars(r.Args.Name), rgName,
				r.GetCloudSpecificLocation(),
			),
			AccountTier:                "Standard",
			AccountReplicationType:     "GZRS",
			AllowNestedItemsToBePublic: true,
			BlobProperties: object_storage.BlobProperties{
				VersioningEnabled: r.Args.Versioning,
			},
		}

		return []output.TfBlock{
			storageAccount,
			object_storage_object.AzureStorageContainer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{
						ResourceId: fmt.Sprintf("%s_%s", r.ResourceId, "public"),
					},
					Name: "public",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "blob",
			}, object_storage_object.AzureStorageContainer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{
						ResourceId: fmt.Sprintf("%s_%s", r.ResourceId, "private"),
					},
					Name: "private",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "private",
			}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *ObjectStorage) GetAssociatedPublicContainerResourceName() string {
	if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_public.name", r.ResourceId)
	}
	return ""
}

func (r *ObjectStorage) GetAssociatedPrivateContainerResourceName() string {
	if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_private.name", r.ResourceId)
	}
	return ""
}

func (r *ObjectStorage) GetResourceName() string {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return fmt.Sprintf("aws_s3_bucket.%s.id", r.ResourceId)
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return fmt.Sprintf("azurerm_storage_account.%s.name", r.ResourceId)
	}
	return ""
}

func (r *ObjectStorage) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}

func (r *ObjectStorage) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return "aws_s3_bucket", nil
	case commonpb.CloudProvider_AZURE:
		return "azurerm_storage_account", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
