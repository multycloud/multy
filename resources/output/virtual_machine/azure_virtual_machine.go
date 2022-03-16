package virtual_machine

import (
	"multy-go/resources/common"
)

const AzureResourceName = "azurerm_linux_virtual_machine"

type AzureVirtualMachine struct {
	*common.AzResource            `hcl:",squash" default:"name=azurerm_linux_virtual_machine"`
	Location                      string                    `hcl:"location"`
	Size                          string                    `hcl:"size"`
	NetworkInterfaceIds           []string                  `hcl:"network_interface_ids"`
	CustomData                    string                    `hcl:"custom_data,expr" hcle:"omitempty"`
	OsDisk                        AzureOsDisk               `hcl:"os_disk"`
	AdminUsername                 string                    `hcl:"admin_username"`
	AdminPassword                 string                    `hcl:"admin_password,expr" hcle:"omitempty"`
	AdminSshKey                   AzureAdminSshKey          `hcl:"admin_ssh_key" hcle:"omitempty"`
	SourceImageReference          AzureSourceImageReference `hcl:"source_image_reference"`
	DisablePasswordAuthentication bool                      `hcl:"disable_password_authentication"`
	Identity                      AzureIdentity             `hcl:"identity"`
}

type AzureIdentity struct {
	Type string `hcl:"type"`
}

type AzureAdminSshKey struct {
	Username  string `hcl:"username"`
	PublicKey string `hcl:"public_key,expr"`
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
