package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/resources/types"
)

type AzureObjectStorage struct {
	*types.ObjectStorage
}

func InitObjectStorage(vn *types.ObjectStorage) resources.ResourceTranslator[*resourcespb.ObjectStorageResource] {
	return AzureObjectStorage{vn}
}

func (r AzureObjectStorage) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.ObjectStorageResource, error) {
	out := &resourcespb.ObjectStorageResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Versioning:  r.Args.Versioning,
		GcpOverride: r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}
	out.AzureOutputs = &resourcespb.ObjectStorageAzureOutputs{}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage.AzureStorageAccount](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AzureOutputs.StorageAccountId = stateResource.ResourceId
		out.Name = stateResource.Name
		out.Versioning = len(stateResource.BlobProperties) > 0 && stateResource.BlobProperties[0].VersioningEnabled
		output.AddToStatuses(statuses, "azure_storage_account", output.MaybeGetPlannedChageById[object_storage.AzureStorageAccount](plan, r.ResourceId))
	} else {
		statuses["azure_storage_account"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if privContainer, exists, err := output.MaybeGetParsedById[object_storage_object.AzureStorageContainer](state, r.getPrivateContainerId()); exists {
		if err != nil {
			return nil, err
		}
		out.AzureOutputs.PrivateStorageContainerId = privContainer.ResourceId
		output.AddToStatuses(statuses, "azure_private_storage_container", output.MaybeGetPlannedChageById[object_storage_object.AzureStorageContainer](plan, r.getPrivateContainerId()))
	} else {
		statuses["azure_private_storage_container"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if publicContainer, exists, err := output.MaybeGetParsedById[object_storage_object.AzureStorageContainer](state, r.getPublicContainerId()); exists {
		if err != nil {
			return nil, err
		}
		out.AzureOutputs.PublicStorageContainerId = publicContainer.ResourceId
		output.AddToStatuses(statuses, "azure_public_storage_container", output.MaybeGetPlannedChageById[object_storage_object.AzureStorageContainer](plan, r.getPublicContainerId()))
	} else {
		statuses["azure_public_storage_container"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureObjectStorage) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)

	storageAccount := object_storage.AzureStorageAccount{
		AzResource: common.NewAzResource(
			r.ResourceId, r.Args.Name, rgName,
			r.GetCloudSpecificLocation(),
		),
		AccountTier:                "Standard",
		AccountReplicationType:     "GZRS",
		AllowNestedItemsToBePublic: true,
		BlobProperties: []object_storage.BlobProperties{{
			VersioningEnabled: r.Args.Versioning,
		}},
	}

	return []output.TfBlock{
		storageAccount,
		object_storage_object.AzureStorageContainer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{
					ResourceId: r.getPublicContainerId(),
				},
				Name: "public",
			},
			StorageAccountName:  storageAccount.GetResourceName(),
			ContainerAccessType: "blob",
		}, object_storage_object.AzureStorageContainer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{
					ResourceId: r.getPrivateContainerId(),
				},
				Name: "private",
			},
			StorageAccountName:  storageAccount.GetResourceName(),
			ContainerAccessType: "private",
		}}, nil
}

func (r AzureObjectStorage) getPrivateContainerId() string {
	return fmt.Sprintf("%s_%s", r.ResourceId, "private")
}

func (r AzureObjectStorage) getPublicContainerId() string {
	return fmt.Sprintf("%s_%s", r.ResourceId, "public")
}

func (r AzureObjectStorage) GetMainResourceName() (string, error) {
	return object_storage.AzureResourceName, nil
}
