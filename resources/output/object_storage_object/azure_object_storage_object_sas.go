package object_storage_object

import (
	"fmt"
	"multy/resources/output"
)

type AzureStorageAccountBlobSas struct {
	*output.TerraformDataSource           `hcl:",squash" default:"name=azurerm_storage_account_blob_container_sas"`
	ConnectionString                      string `hcl:"connection_string,expr"`
	ContainerName                         string `hcl:"container_name,expr"`
	Start                                 string `hcl:"start"`
	Expiry                                string `hcl:"expiry"  hcle:"omitempty"`
	AzureStorageAccountBlobSasPermissions `hcl:"permissions"`
}

type AzureStorageAccountBlobSasPermissions struct {
	Read   bool `hcl:"read"`
	Write  bool `hcl:"write"`
	Add    bool `hcl:"add"`
	Create bool `hcl:"create"`
	List   bool `hcl:"list"`
	Delete bool `hcl:"delete"`
}

func (sas AzureStorageAccountBlobSas) GetSignedUrl(storageAccountName string, containerName string,
	blobName string) string {
	return fmt.Sprintf(
		"https://${%s}.blob.core.windows.net/${%s}/${%s}${data.azurerm_storage_account_blob_container_sas.%s.sas}",
		storageAccountName, containerName, blobName, sas.ResourceId,
	)
}
