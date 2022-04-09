package route_table

import "github.com/multycloud/multy/resources/common"

const AzureResourceName = "azurerm_route_table"

type AzureRouteTable struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_route_table"`
	Routes             []AzureRouteTableRoute `hcl:"route,blocks"`
}

type AzureRouteTableRoute struct {
	Name          string `hcl:"name"`
	AddressPrefix string `hcl:"address_prefix"`
	NextHopType   string `hcl:"next_hop_type"` // VirtualNetworkGateway , VnetLocal , Internet , VirtualAppliance and None
}

type AzureRouteTableAssociation struct {
	*common.AzResource `hcl:",squash"`
	RouteTableId       string `hcl:"route_table_id,expr"`
	SubnetId           string `hcl:"subnet_id,expr"`
}
