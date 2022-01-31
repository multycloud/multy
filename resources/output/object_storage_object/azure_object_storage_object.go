package object_storage_object

import (
	"fmt"
	"multy-go/resources/common"
)

// azurerm_storage_blob
type AzureStorageAccountBlob struct {
	common.AzResource    `hcl:",squash"`
	StorageAccountName   string `hcl:"storage_account_name,expr"`
	StorageContainerName string `hcl:"storage_container_name,expr"`
	Type                 string `hcl:"type"`
	SourceContent        string `hcl:"source_content"  hcle:"omitempty"`
	ContentType          string `hcl:"content_type" hcle:"omitempty"`
	Source               string `hcl:"content_type" hcle:"omitempty"`
}

// azurerm_storage_container
type AzureStorageContainer struct {
	common.AzResource   `hcl:",squash"`
	StorageAccountName  string `hcl:"storage_account_name,expr"`
	ContainerAccessType string `hcl:"container_access_type"`
}

func (r AzureStorageContainer) GetResourceName() string {
	return fmt.Sprintf("azurerm_storage_container.%s.name", r.ResourceId)
}
