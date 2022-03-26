package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

type ObjectStorage struct {
	*resources.CommonResourceParams
	Name       string     `hcl:"name"`
	Acl        []AclRules `hcl:"acl,optional"`
	Versioning bool       `hcl:"versioning,optional"`
}

type AclRules struct{}

func (r *ObjectStorage) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		var awsResources []output.TfBlock
		s3Bucket := object_storage.AwsS3Bucket{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
			},
			Bucket: r.Name,
		}
		awsResources = append(awsResources, s3Bucket)

		if r.Versioning {
			awsResources = append(awsResources, object_storage.AwsS3BucketVersioning{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				BucketId:                s3Bucket.GetBucketId(),
				VersioningConfiguration: object_storage.VersioningConfiguration{Status: "Enabled"},
			})
		}
		return awsResources, nil
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)

		storageAccount := object_storage.AzureStorageAccount{
			AzResource: common.NewAzResource(
				r.GetTfResourceId(cloud), common.RemoveSpecialChars(r.Name), rgName,
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			),
			AccountTier:                "Standard",
			AccountReplicationType:     "GZRS",
			AllowNestedItemsToBePublic: true,
			BlobProperties: object_storage.BlobProperties{
				VersioningEnabled: r.Versioning,
			},
		}

		return []output.TfBlock{
			storageAccount,
			object_storage_object.AzureStorageContainer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{
						ResourceId: fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "public"),
					},
					Name: "public",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "blob",
			}, object_storage_object.AzureStorageContainer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{
						ResourceId: fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "private"),
					},
					Name: "private",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "private",
			}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *ObjectStorage) GetAssociatedPublicContainerResourceName(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_public.name", r.GetTfResourceId(common.AZURE))
	}
	return ""
}

func (r *ObjectStorage) GetAssociatedPrivateContainerResourceName(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_private.name", r.GetTfResourceId(common.AZURE))
	}
	return ""
}

func (r *ObjectStorage) GetResourceName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return fmt.Sprintf("aws_s3_bucket.%s.id", r.GetTfResourceId(common.AWS))
	} else if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_account.%s.name", r.GetTfResourceId(common.AZURE))
	}
	return ""
}

func (r *ObjectStorage) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *ObjectStorage) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return "aws_s3_bucket", nil
	case common.AZURE:
		return "azurerm_storage_account", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
