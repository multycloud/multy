package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
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
	VirtualNetwork *VirtualNetwork   `mhcl:"ref=virtual_network"`
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

func (r *RouteTable) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		vpcId, err := resources.GetMainOutputId(r.VirtualNetwork, cloud)
		if err != nil {
			return nil, err
		}
		rt := route_table.AwsRouteTable{
			AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
			VpcId:       vpcId,
		}
		gtw, err := r.VirtualNetwork.GetAssociatedInternetGateway(cloud)
		if err != nil {
			return nil, err
		}

		var routes []route_table.AwsRouteTableRoute
		for _, route := range r.Routes {
			if strings.EqualFold(route.Destination, INTERNET) {
				routes = append(
					routes, route_table.AwsRouteTableRoute{
						CidrBlock: route.CidrBlock,
						GatewayId: gtw,
					},
				)
			}
		}
		rt.Routes = routes

		return []output.TfBlock{rt}, nil
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
		return []output.TfBlock{rt}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *RouteTable) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: route_table.AwsResourceName, common.AZURE: route_table.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func (r *RouteTable) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, r.CommonResourceParams.Validate(ctx, cloud)...)
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

func (r *RouteTable) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return route_table.AwsResourceName, nil
	case common.AZURE:
		return route_table.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
