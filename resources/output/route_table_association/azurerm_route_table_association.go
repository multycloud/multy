package route_table_association

import "multy/resources/common"

const AzureResourceName = "azurerm_subnet_route_table_association"

type AzureRouteTableAssociation struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_subnet_route_table_association"`
	SubnetId           string `hcl:"subnet_id"`
	RouteTableId       string `hcl:"route_table_id"`
}
