package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/object_storage"
	"multy-go/resources/output/object_storage_object"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

/*
AWS: aws_s3_bucket
Azure: azurerm_storage_account
*/

/*
resource "aws_s3_bucket" "b" {
  bucket = "my-tf-test-bucket"
  acl    = "private"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}
resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
*/

type ObjectStorage struct {
	*resources.CommonResourceParams
	Name         string     `hcl:"name"`
	Acl          []AclRules `hcl:"acl,optional"`
	Versioning   bool       `hcl:"versioning,optional"`
	RandomSuffix bool       `hcl:"random_suffix,optional"` // name must be unique
}

type AclRules struct{}

func (r *ObjectStorage) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	name := r.Name
	if r.RandomSuffix {
		name += fmt.Sprintf("-%s", common.RandomString(6))
	}
	if cloud == common.AWS {
		return []interface{}{object_storage.AwsS3Bucket{
			AwsResource: common.AwsResource{
				ResourceName: "aws_s3_bucket",
				ResourceId:   r.GetTfResourceId(cloud),
			}, Bucket: name}}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)

		storageAccount := object_storage.AzureStorageAccount{
			AzResource: common.AzResource{
				ResourceName:      "azurerm_storage_account",
				ResourceId:        r.GetTfResourceId(cloud),
				ResourceGroupName: rgName,
				Name:              common.RemoveSpecialChars(name),
				Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			},
			AccountTier:            "Standard",
			AccountReplicationType: "GZRS",
			AllowBlobPublicAccess:  true,
		}

		return []interface{}{
			storageAccount,
			object_storage_object.AzureStorageContainer{
				AzResource: common.AzResource{
					ResourceName: "azurerm_storage_container",
					ResourceId:   fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "public"),
					Name:         "public",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "blob",
			}, object_storage_object.AzureStorageContainer{
				AzResource: common.AzResource{
					ResourceName: "azurerm_storage_container",
					ResourceId:   fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "private"),
					Name:         "private",
				},
				StorageAccountName:  storageAccount.GetResourceName(),
				ContainerAccessType: "private",
			}}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *ObjectStorage) GetAssociatedPublicContainerResourceName(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_public.name", r.GetTfResourceId(common.AZURE))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *ObjectStorage) GetAssociatedPrivateContainerResourceName(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_container.%s_private.name", r.GetTfResourceId(common.AZURE))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *ObjectStorage) GetResourceName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return fmt.Sprintf("aws_s3_bucket.%s.id", r.GetTfResourceId(common.AWS))
	} else if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_storage_account.%s.name", r.GetTfResourceId(common.AZURE))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *ObjectStorage) Validate(ctx resources.MultyContext) {
	return
}
