package route_table_association

import "multy-go/resources/common"

const AzureResourceName = "azurerm_subnet_route_table_association"

type AzureRouteTableAssociation struct {
	common.AzResource `hcl:",squash"`
	SubnetId          string `hcl:"subnet_id"`
	RouteTableId      string `hcl:"route_table_id,expr"`
}
