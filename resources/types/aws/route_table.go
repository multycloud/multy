package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table"
	"github.com/multycloud/multy/resources/types"
)

type AwsRouteTable struct {
	*types.RouteTable
}

func InitRouteTable(vn *types.RouteTable) resources.ResourceTranslator[*resourcespb.RouteTableResource] {
	return AwsRouteTable{vn}
}

func (r AwsRouteTable) FromState(state *output.TfState) (*resourcespb.RouteTableResource, error) {
	if flags.DryRun {
		return &resourcespb.RouteTableResource{
			CommonParameters: &commonpb.CommonChildResourceParameters{
				ResourceId:  r.ResourceId,
				NeedsUpdate: false,
			},
			Name:             r.Args.Name,
			VirtualNetworkId: r.Args.VirtualNetworkId,
			Routes:           r.Args.Routes,
		}, nil
	}
	out := new(resourcespb.RouteTableResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}
	out.AwsOutputs = &resourcespb.RouteTableAwsOutputs{}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[route_table.AwsRouteTable](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.AwsResource.Tags["Name"]
		out.VirtualNetworkId = r.Args.VirtualNetworkId
		var routes []*resourcespb.Route
		for _, r := range stateResource.Routes {
			route := &resourcespb.Route{
				CidrBlock:   r.CidrBlock,
				Destination: resourcespb.RouteDestination_INTERNET,
			}
			if r.GatewayId == "" {
				route.Destination = resourcespb.RouteDestination_UNKNOWN_DESTINATION
			}

			routes = append(routes, route)
		}
		out.Routes = routes
		out.AwsOutputs.RouteTableId = stateResource.ResourceId
	} else {
		statuses["aws_route_table"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AwsRouteTable) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	vpcId, err := resources.GetMainOutputId(AwsVirtualNetwork{r.VirtualNetwork})
	if err != nil {
		return nil, err
	}
	rt := route_table.AwsRouteTable{
		AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
		VpcId:       vpcId,
	}

	var routes []route_table.AwsRouteTableRoute
	for _, route := range r.Args.Routes {

		if route.Destination == resourcespb.RouteDestination_INTERNET {
			routes = append(
				routes, route_table.AwsRouteTableRoute{
					CidrBlock: route.CidrBlock,
					GatewayId: AwsVirtualNetwork{r.VirtualNetwork}.GetAssociatedInternetGateway(),
				},
			)
		}
	}
	rt.Routes = routes

	return []output.TfBlock{rt}, nil
}

func (r AwsRouteTable) GetMainResourceName() (string, error) {
	return route_table.AwsResourceName, nil
}
