package azure_resources

import (
	"encoding/base64"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	"github.com/multycloud/multy/resources/types"
)

type AzureObjectStorageObject struct {
	*types.ObjectStorageObject
}

func InitObjectStorageObject(vn *types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource] {
	return AzureObjectStorageObject{vn}
}

func (r AzureObjectStorageObject) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.ObjectStorageObjectResource, error) {
	out := &resourcespb.ObjectStorageObjectResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:            r.Args.Name,
		Acl:             r.Args.Acl,
		ObjectStorageId: r.Args.ObjectStorageId,
		ContentBase64:   r.Args.ContentBase64,
		ContentType:     r.Args.ContentType,
		Source:          r.Args.Source,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage_object.AzureStorageAccountBlob](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.ContentType = stateResource.ContentType
		out.Source = stateResource.Source
		out.ContentBase64 = base64.StdEncoding.EncodeToString([]byte(stateResource.SourceContent))
		if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
			out.Url = fmt.Sprintf("https://%s.blob.core.windows.net/public/%s", stateResource.StorageAccountName, r.Args.Name)
		}
		out.AzureOutputs = &resourcespb.ObjectStorageObjectAzureOutputs{StorageBlobId: stateResource.ResourceId}
		output.AddToStatuses(statuses, "azure_storage_account_blob", output.MaybeGetPlannedChageById[object_storage_object.AzureStorageAccountBlob](plan, r.ResourceId))
	} else {
		statuses["azure_storage_account_blob"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureObjectStorageObject) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var containerName string
	if r.Args.Acl == resourcespb.ObjectStorageObjectAcl_PUBLIC_READ {
		containerName = fmt.Sprintf("azurerm_storage_container.%s_public.name", r.ObjectStorage.ResourceId)
	} else {
		containerName = fmt.Sprintf("azurerm_storage_container.%s_private.name", r.ObjectStorage.ResourceId)
	}
	return []output.TfBlock{
		object_storage_object.AzureStorageAccountBlob{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				Name:              r.Args.Name,
			},
			StorageAccountName:   fmt.Sprintf("azurerm_storage_account.%s.name", r.ObjectStorage.ResourceId),
			StorageContainerName: containerName,
			Type:                 "Block",
			ContentType:          r.Args.ContentType,
			SourceContent:        fmt.Sprintf("base64decode(\"%s\")", r.Args.ContentBase64),
		}}, nil
}

func (r AzureObjectStorageObject) GetMainResourceName() (string, error) {
	return "azurerm_storage_blob", nil
}
