package object_storage_object

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

// azurerm_storage_blob
type AzureStorageAccountBlob struct {
	*common.AzResource   `hcl:",squash" default:"name=azurerm_storage_blob"`
	StorageAccountName   string `hcl:"storage_account_name,expr" json:"storage_account_name"`
	StorageContainerName string `hcl:"storage_container_name,expr" json:"storage_container_name"`
	Type                 string `hcl:"type" json:"type"`
	SourceContent        string `hcl:"source_content,expr"  hcle:"omitempty" json:"source_content"`
	ContentType          string `hcl:"content_type" hcle:"omitempty" json:"content_type"`
	Source               string `hcl:"source,expr" hcle:"omitempty" json:"source"`
}

// azurerm_storage_container
type AzureStorageContainer struct {
	*common.AzResource  `hcl:",squash"  default:"name=azurerm_storage_container"`
	StorageAccountName  string `hcl:"storage_account_name,expr"`
	ContainerAccessType string `hcl:"container_access_type"`
}

func (r AzureStorageContainer) GetResourceName() string {
	return fmt.Sprintf("azurerm_storage_container.%s.name", r.ResourceId)
}
