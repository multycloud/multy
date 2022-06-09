package subnet

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
)

const AzureResourceName = "azurerm_subnet"

type AzureSubnet struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_subnet"`
	AddressPrefixes    []string `hcl:"address_prefixes" json:"address_prefixes"`
	VirtualNetworkName string   `hcl:"virtual_network_name,expr" json:"virtual_network_name"`
	ServiceEndpoints   []string `hcl:"service_endpoints" hcle:"omitempty" json:"service_endpoints"`
}

func (subnet *AzureSubnet) GetId() string {
	return fmt.Sprintf("azurerm_subnet.%s.id", subnet.ResourceId)
}
