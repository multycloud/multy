package virtual_network

import (
	"multy-go/resources/common"
)

const AzureResourceName = "azurerm_virtual_network"

type AzureVnet struct {
	*common.AzResource `hcl:",squash"  default:"name=azurerm_virtual_network"`
	AddressSpace       []string `hcl:"address_space"`
}
