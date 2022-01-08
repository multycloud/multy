package virtual_machine

import (
	"multy-go/resources/common"
)

const AzureResourceName = "azurerm_linux_virtual_machine"

type AzureVirtualMachine struct {
	common.AzResource             `hcl:",squash"`
	Location                      string                    `hcl:"location"`
	Size                          string                    `hcl:"size"`
	NetworkInterfaceIds           []string                  `hcl:"network_interface_ids,expr"`
	CustomData                    string                    `hcl:"custom_data" hcle:"omitempty"`
	OsDisk                        AzureOsDisk               `hcl:"os_disk"`
	AdminUsername                 string                    `hcl:"admin_username"`
	AdminPassword                 string                    `hcl:"admin_password"`
	SourceImageReference          AzureSourceImageReference `hcl:"source_image_reference"`
	DisablePasswordAuthentication bool                      `hcl:"disable_password_authentication"`
}

type AzureOsDisk struct {
	Caching            string `hcl:"caching"`
	StorageAccountType string `hcl:"storage_account_type"`
}

type AzureSourceImageReference struct {
	Publisher string `hcl:"publisher"`
	Offer     string `hcl:"offer"`
	Sku       string `hcl:"sku"`
	Version   string `hcl:"version"`
}
