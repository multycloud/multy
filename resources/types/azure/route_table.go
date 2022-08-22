package azure_resources

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

type AzureRouteTable struct {
	*types.RouteTable
}

func InitRouteTable(vn *types.RouteTable) resources.ResourceTranslator[*resourcespb.RouteTableResource] {
	return AzureRouteTable{vn}
}

func (r AzureRouteTable) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.RouteTableResource, error) {
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
	out.AzureOutputs = &resourcespb.RouteTableAzureOutputs{}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[route_table.AzureRouteTable](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.VirtualNetworkId = r.Args.VirtualNetworkId
		var routes []*resourcespb.Route
		for _, r := range stateResource.Routes {
			route := &resourcespb.Route{
				CidrBlock:   r.AddressPrefix,
				Destination: resourcespb.RouteDestination_INTERNET,
			}
			if r.NextHopType != "Internet" {
				route.Destination = resourcespb.RouteDestination_UNKNOWN_DESTINATION
			}

			routes = append(routes, route)
		}
		out.Routes = routes
		out.AzureOutputs.RouteTableId = stateResource.ResourceId
	} else {
		statuses["azure_route_table"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureRouteTable) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rt := route_table.AzureRouteTable{
		AzResource: common.NewAzResource(
			r.ResourceId, r.Args.Name, GetResourceGroupName(r.VirtualNetwork.Args.GetCommonParameters().ResourceGroupId),
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
					NextHopType:   "Internet",
				},
			)
		}
	}
	rt.Routes = routes
	return []output.TfBlock{rt}, nil
}

func (r AzureRouteTable) GetMainResourceName() (string, error) {
	return route_table.AzureResourceName, nil
}
