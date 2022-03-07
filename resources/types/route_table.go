package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/route_table"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
	"strings"
)

/*
Notes:
AWS: Internet route to IGW
Azure: Internet route nextHop Internet
*/

type RouteTable struct {
	*resources.CommonResourceParams
	Name           string            `hcl:"name"`
	VirtualNetwork *VirtualNetwork   `mhcl:"ref=virtual_network,optional"`
	Routes         []RouteTableRoute `hcl:"routes,optional"`
}

type RouteTableRoute struct {
	CidrBlock   string `cty:"cidr_block"`
	Destination string `cty:"destination"` // allowed: Internet
}

const (
	INTERNET       = "Internet"
	VIRTUALNETWORK = "VirtualNetwork"
)

func (r *RouteTable) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		rt := route_table.AwsRouteTable{
			AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
			VpcId:       resources.GetMainOutputId(r.VirtualNetwork, cloud),
		}

		var routes []route_table.AwsRouteTableRoute
		for _, route := range r.Routes {
			if strings.EqualFold(route.Destination, INTERNET) {
				routes = append(
					routes, route_table.AwsRouteTableRoute{
						CidrBlock: route.CidrBlock,
						GatewayId: r.VirtualNetwork.GetAssociatedInternetGateway(cloud),
					},
				)
			}
		}
		rt.Routes = routes

		return []output.TfBlock{rt}
	} else if cloud == common.AZURE {
		rt := route_table.AzureRouteTable{
			AzResource: common.NewAzResource(
				r.GetTfResourceId(cloud), r.Name, rg.GetResourceGroupName(r.ResourceGroupId, cloud),
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			),
		}

		var routes []route_table.AzureRouteTableRoute
		for _, route := range r.Routes {
			if strings.EqualFold(route.Destination, INTERNET) {
				routes = append(
					routes, route_table.AzureRouteTableRoute{
						Name:          "internet",
						AddressPrefix: route.CidrBlock,
						NextHopType:   INTERNET,
					},
				)
			}
		}
		rt.Routes = routes
		return []output.TfBlock{rt}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *RouteTable) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: route_table.AwsResourceName, common.AZURE: route_table.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func (r *RouteTable) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	if len(r.Routes) > 20 {
		errs = append(errs, r.NewError("routes", fmt.Sprintf("\"%d\" exceeds routes limit is 20", len(r.Routes))))
	}
	for _, route := range r.Routes {
		if !strings.EqualFold(route.Destination, INTERNET) {
			errs = append(errs, r.NewError("route", fmt.Sprintf("\"%s\" must be Internet", route.Destination)))
		}
		//	if route.CidrBlock valid CIDR
	}
	return errs
}

func (r *RouteTable) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return route_table.AwsResourceName
	case common.AZURE:
		return route_table.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
