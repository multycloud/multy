package azure_resources

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
)

type AzureVirtualNetwork struct {
	*types.VirtualNetwork
}

func InitVirtualNetwork(vn *types.VirtualNetwork) resources.ResourceTranslator[*resourcespb.VirtualNetworkResource] {
	return AzureVirtualNetwork{vn}
}

func (r AzureVirtualNetwork) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VirtualNetworkResource, error) {
	if flags.DryRun {
		return &resourcespb.VirtualNetworkResource{
			CommonParameters: &commonpb.CommonResourceParameters{
				ResourceId:      r.ResourceId,
				ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
				Location:        r.Args.CommonParameters.Location,
				CloudProvider:   r.Args.CommonParameters.CloudProvider,
				NeedsUpdate:     false,
			},
			Name:        r.Args.Name,
			CidrBlock:   r.Args.CidrBlock,
			GcpOverride: r.Args.GcpOverride,
		}, nil
	}
	out := new(resourcespb.VirtualNetworkResource)
	out.CommonParameters = &commonpb.CommonResourceParameters{
		ResourceId:      r.ResourceId,
		ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
		Location:        r.Args.CommonParameters.Location,
		CloudProvider:   r.GetCloud(),
		NeedsUpdate:     false,
	}
	out.GcpOverride = r.Args.GcpOverride
	out.AzureOutputs = &resourcespb.VirtualNetworkAzureOutputs{}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[virtual_network.AzureVnet](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.CidrBlock = stateResource.AddressSpace[0]

		out.AzureOutputs.VirtualNetworkId = stateResource.AzResource.ResourceId
		output.AddToStatuses(statuses, "azure_virtual_network", output.MaybeGetPlannedChageById[virtual_network.AzureVnet](plan, r.ResourceId))
	} else {
		statuses["azure_virtual_network"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[route_table.AzureRouteTable](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AzureOutputs.LocalRouteTableId = stateResource.ResourceId
		output.AddToStatuses(statuses, "azure_local_route_table", output.MaybeGetPlannedChageById[route_table.AzureRouteTable](plan, r.ResourceId))
	} else {
		statuses["azure_local_route_table"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureVirtualNetwork) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{virtual_network.AzureVnet{
		AzResource: common.NewAzResource(
			r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
			r.GetCloudSpecificLocation(),
		),
		AddressSpace: []string{r.Args.CidrBlock},
	}, route_table.AzureRouteTable{
		AzResource: common.NewAzResource(
			r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
			r.GetCloudSpecificLocation(),
		),
		Routes: []route_table.AzureRouteTableRoute{{
			Name:          "local",
			AddressPrefix: "0.0.0.0/0",
			NextHopType:   "VnetLocal",
		}},
	}}, nil
}

func (r AzureVirtualNetwork) GetMainResourceName() (string, error) {
	return virtual_network.AzureResourceName, nil
}

func (r AzureVirtualNetwork) GetAssociatedRouteTableId() string {
	return fmt.Sprintf("%s.%s.id", route_table.AzureResourceName, r.ResourceId)
}
