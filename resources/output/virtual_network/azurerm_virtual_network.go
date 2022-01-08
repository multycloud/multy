package virtual_network

import (
	"multy-go/resources/common"
)

const AzureResourceName = "azurerm_virtual_network"

type AzureVnet struct {
	common.AzResource `hcl:",squash"`
	AddressSpace      []string `hcl:"address_space"`
}
