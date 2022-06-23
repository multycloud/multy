package virtual_machine

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"golang.org/x/exp/slices"
	"strings"
)

const AzureResourceName = "azurerm_linux_virtual_machine"

type AzureVirtualMachine struct {
	*common.AzResource            `hcl:",squash" default:"name=azurerm_linux_virtual_machine"`
	Location                      string                    `hcl:"location"`
	Size                          string                    `hcl:"size"`
	NetworkInterfaceIds           []string                  `hcl:"network_interface_ids,expr"`
	CustomData                    string                    `hcl:"custom_data" hcle:"omitempty"`
	OsDisk                        AzureOsDisk               `hcl:"os_disk"`
	AdminUsername                 string                    `hcl:"admin_username"`
	AdminPassword                 string                    `hcl:"admin_password,expr" hcle:"omitempty"`
	AdminSshKey                   AzureAdminSshKey          `hcl:"admin_ssh_key" hcle:"omitempty"`
	SourceImageReference          AzureSourceImageReference `hcl:"source_image_reference"`
	DisablePasswordAuthentication bool                      `hcl:"disable_password_authentication"`
	Identity                      AzureIdentity             `hcl:"identity"`
	Identities                    []AzureGeneratedIdentity  `json:"identity"`
	ComputerName                  string                    `hcl:"computer_name"`
}

type AzureGeneratedIdentity struct {
	// outputs
	PrincipalId string `json:"principal_id" hcle:"omitempty"`
}

type AzureIdentity struct {
	Type string `hcl:"type"`

	// outputs
	PrincipalId string `json:"principal_id" hcle:"omitempty"`
}

type AzureAdminSshKey struct {
	Username  string `hcl:"username"`
	PublicKey string `hcl:"public_key"`
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

func GetLatestAzureSourceImageReference(ref *resourcespb.ImageReference) (AzureSourceImageReference, error) {
	var offer string
	var publisher string
	version := ref.Version
	switch ref.Os {
	case resourcespb.ImageReference_UBUNTU:
		offer = "UbuntuServer"
		publisher = "Canonical"
		version = fmt.Sprintf("%s-LTS", version)
		if ref.Version == "20.04" {
			offer = "0001-com-ubuntu-server-focal"
			version = "20_04-lts"
		}
	case resourcespb.ImageReference_DEBIAN:
		offer = fmt.Sprintf("debian-%s", ref.Version)
		publisher = "Debian"
	case resourcespb.ImageReference_CENT_OS:
		offer = "CentOs"
		publisher = "OpenLogic"
		dotVersions := []string{"7.2", "7.3", "7.4", "7.5", "7.6", "7.7", "8.0"}
		if !slices.Contains(dotVersions, version) {
			version = strings.Replace(version, ".", "_", 1)
		}
	default:
		return AzureSourceImageReference{}, fmt.Errorf("unknown operating system distibution %s", ref.Os)
	}

	return AzureSourceImageReference{
		Publisher: publisher,
		Offer:     offer,
		Sku:       version,
		Version:   "latest",
	}, nil

}
