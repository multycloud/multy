package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"golang.org/x/exp/slices"
)

type GcpRouteTable struct {
	*types.RouteTable
}

func InitRouteTable(rg *types.RouteTable) resources.ResourceTranslator[*resourcespb.RouteTableResource] {
	return GcpRouteTable{rg}
}

func (r GcpRouteTable) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	associations := resources.GetAllResourcesWithRef[*types.RouteTableAssociation](ctx,
		func(rt *types.RouteTableAssociation) *types.RouteTable {
			return rt.RouteTable
		},
		r.RouteTable,
	)

	var routes []output.TfBlock
	for i, routeId := range r.getRouteIds() {
		routeName := fmt.Sprintf("%s-%d", r.Args.Name, i)
		route := r.Args.Routes[i]
		outputRoute := &route_table.GoogleComputeRoute{
			GcpResource: common.NewGcpResource(routeId, routeName, r.VirtualNetwork.Args.GetGcpOverride().GetProject()),
			DestRange:   route.CidrBlock,
			Network:     fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
			Priority:    1000,
			Tags:        r.getTags(associations),
		}
		if route.Destination == resourcespb.RouteDestination_INTERNET {
			outputRoute.NextHopGateway = "default-internet-gateway"
		}
		routes = append(routes, outputRoute)
	}

	return routes, nil
}

func (r GcpRouteTable) getTags(rtas []*types.RouteTableAssociation) []string {
	if len(rtas) == 0 {
		return []string{"no-subnet-attached"}
	}
	var out []string
	for _, rta := range rtas {
		out = append(out, GcpSubnet{rta.Subnet}.getNetworkTag())
	}
	slices.Sort(out)
	return out
}

func (r GcpRouteTable) getRouteIds() (out []string) {
	for i := range r.Args.Routes {
		out = append(out, fmt.Sprintf("%s-%d", r.ResourceId, i))
	}
	return
}

func (r GcpRouteTable) FromState(state *output.TfState) (*resourcespb.RouteTableResource, error) {
	out := &resourcespb.RouteTableResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:             r.Args.Name,
		VirtualNetworkId: r.Args.VirtualNetworkId,
		Routes:           r.Args.Routes,
	}

	if flags.DryRun {
		return out, nil
	}

	out.GcpOutputs = &resourcespb.RouteTableGcpOutputs{}

	for _, routeId := range r.getRouteIds() {
		stateResource, err := output.GetParsedById[route_table.GoogleComputeRoute](state, routeId)
		if err != nil {
			return nil, err
		}
		out.GcpOutputs.ComputeRouteId = append(out.GcpOutputs.ComputeRouteId, stateResource.SelfLink)
	}

	return out, nil
}

func (r GcpRouteTable) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main id exists for route table gcp")
}
