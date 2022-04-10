package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

/*
Notes:
AWS: Internet route to IGW
Azure: Internet route nextHop Internet
*/

type RouteTable struct {
	resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.RouteTableArgs]

	VirtualNetwork *VirtualNetwork `mhcl:"ref=virtual_network"`
}

type RouteTableRoute struct {
	CidrBlock   string `cty:"cidr_block"`
	Destination string `cty:"destination"` // allowed: Internet
}

const (
	INTERNET       = "Internet"
	VIRTUALNETWORK = "VirtualNetwork"
)

func NewRouteTable(resourceId string, args *resourcespb.RouteTableArgs, others resources.Resources) (*RouteTable, error) {
	rt := &RouteTable{
		ChildResourceWithId: resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.RouteTableArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{rt.NewValidationError(err.Error(), "virtual_network_id")})
	}
	rt.Parent = vn
	rt.VirtualNetwork = vn
	return rt, nil
}

func (r *RouteTable) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		vpcId, err := resources.GetMainOutputId(r.VirtualNetwork)
		if err != nil {
			return nil, err
		}
		rt := route_table.AwsRouteTable{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			VpcId:       vpcId,
		}
		gtw, err := r.VirtualNetwork.GetAssociatedInternetGateway()
		if err != nil {
			return nil, err
		}

		var routes []route_table.AwsRouteTableRoute
		for _, route := range r.Args.Routes {
			if route.Destination == resourcespb.RouteDestination_INTERNET {
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
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		rt := route_table.AzureRouteTable{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rg.GetResourceGroupName(r.VirtualNetwork.Args.GetCommonParameters().ResourceGroupId),
				r.GetCloudSpecificLocation(),
			),
		}

		var routes []route_table.AzureRouteTableRoute
		for _, route := range r.Args.Routes {
			if route.Destination == resourcespb.RouteDestination_INTERNET {
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
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *RouteTable) GetId(cloud commonpb.CloudProvider) string {
	types := map[commonpb.CloudProvider]string{common.AWS: route_table.AwsResourceName, common.AZURE: route_table.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.ResourceId)
}

func (r *RouteTable) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if len(r.Args.Routes) > 20 {
		errs = append(errs, r.NewValidationError(fmt.Sprintf("\"%d\" exceeds routes limit is 20", len(r.Args.Routes)), "routes"))
	}
	for _, route := range r.Args.Routes {
		if route.Destination == resourcespb.RouteDestination_UNKNOWN_DESTINATION {
			errs = append(errs, r.NewValidationError("unknown route destination", "route"))
		}
		//	if route.CidrBlock valid CIDR
	}
	return errs
}

func (r *RouteTable) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return route_table.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return route_table.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
