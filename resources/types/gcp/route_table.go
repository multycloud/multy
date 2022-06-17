package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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
	for i, route := range r.Args.Routes {
		routeName := fmt.Sprintf("%s-%d", r.Args.Name, i)
		routeId := fmt.Sprintf("%s-%d", r.ResourceId, i)
		outputRoute := &route_table.GoogleComputeRoute{
			GcpResource: common.NewGcpResource(routeId, routeName),
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

func (r GcpRouteTable) FromState(_ *output.TfState) (*resourcespb.RouteTableResource, error) {
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

func (r GcpRouteTable) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main id exists for route table gcp")
}
