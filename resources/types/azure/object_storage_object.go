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
	"github.com/multycloud/multy/resources/output/terraform"
	"github.com/multycloud/multy/resources/types"
)

type AzureObjectStorageObject struct {
	*types.ObjectStorageObject
}

func InitObjectStorageObject(vn *types.ObjectStorageObject) resources.ResourceTranslator[*resourcespb.ObjectStorageObjectResource] {
	return AzureObjectStorageObject{vn}
}

func (r AzureObjectStorageObject) FromState(state *output.TfState) (*resourcespb.ObjectStorageObjectResource, error) {
	out := new(resourcespb.ObjectStorageObjectResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}

	id, err := resources.GetMainOutputRef(AzureObjectStorage{r.Parent})
	if err != nil {
		return nil, err
	}

	out.Name = r.Args.Name
	out.ContentBase64 = r.Args.ContentBase64
	out.ContentType = r.Args.ContentType
	out.ObjectStorageId = r.Args.ObjectStorageId
	out.Acl = r.Args.Acl
	out.Source = r.Args.Source

	if !flags.DryRun {
		stateResource, err := output.GetParsed[object_storage.AzureStorageAccount](state, id)
		if err != nil {
			return nil, err
		}
		out.Url = fmt.Sprintf("https://%s.blob.core.windows.net/public/%s", stateResource.AzResource.Name, r.Args.Name)

	} else {
		out.Url = "dryrun"
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
	contentFile := terraform.NewLocalFile(r.ResourceId, r.Args.ContentBase64)
	return []output.TfBlock{
		contentFile,
		object_storage_object.AzureStorageAccountBlob{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				Name:              r.Args.Name,
			},
			StorageAccountName:   fmt.Sprintf("azurerm_storage_account.%s.name", r.ObjectStorage.ResourceId),
			StorageContainerName: containerName,
			Type:                 "Block",
			Source:               contentFile.GetFilename(),
			ContentType:          r.Args.ContentType,
		}}, nil
}

func (r AzureObjectStorageObject) GetMainResourceName() (string, error) {
	return "azurerm_storage_blob", nil
}
