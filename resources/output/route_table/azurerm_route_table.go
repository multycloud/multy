package route_table

import "github.com/multycloud/multy/resources/common"

const AzureResourceName = "azurerm_route_table"

type AzureRouteTable struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_route_table"`
	Routes             []AzureRouteTableRoute `hcl:"route,blocks" json:"route,omitempty"`
}

type AzureRouteTableRoute struct {
	Name          string `hcl:"name" json:"name"`
	AddressPrefix string `hcl:"address_prefix" json:"address_prefix"`
	NextHopType   string `hcl:"next_hop_type" json:"next_hop_type"` // VirtualNetworkGateway , VnetLocal , Internet , VirtualAppliance and None
}
