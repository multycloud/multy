package network_interface

import (
	"github.com/multycloud/multy/resources/common"
)

const AzureResourceName = "azurerm_network_interface"

type AzureNetworkInterface struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_network_interface"`
	IpConfigurations   []AzureIpConfiguration `hcl:"ip_configuration,blocks"`
}

type AzureIpConfiguration struct {
	Name                       string `hcl:"name"`
	PrivateIpAddressAllocation string `hcl:"private_ip_address_allocation"`
	SubnetId                   string `hcl:"subnet_id,expr" hcle:"omitempty"`
	PublicIpAddressId          string `hcl:"public_ip_address_id,expr" hcle:"omitempty"`
	Primary                    bool   `hcl:"primary" hcle:"omitempty"`
}
