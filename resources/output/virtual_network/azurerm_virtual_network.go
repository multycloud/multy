package virtual_network

import (
	"github.com/multycloud/multy/resources/common"
)

const AzureResourceName = "azurerm_virtual_network"

type AzureVnet struct {
	*common.AzResource `hcl:",squash"  default:"name=azurerm_virtual_network"`
	AddressSpace       []string `hcl:"address_space" json:"address_space"`
}
