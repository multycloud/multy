package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
)

type AzureSubnet struct {
	*types.Subnet
}

func InitSubnet(r *types.Subnet) resources.ResourceTranslator[*resourcespb.SubnetResource] {
	return AzureSubnet{r}
}

func (r AzureSubnet) FromState(state *output.TfState) (*resourcespb.SubnetResource, error) {
	if flags.DryRun {
		return &resourcespb.SubnetResource{
			CommonParameters: &commonpb.CommonChildResourceParameters{
				ResourceId:  r.ResourceId,
				NeedsUpdate: false,
			},
			Name:             r.Args.Name,
			CidrBlock:        r.Args.CidrBlock,
			VirtualNetworkId: r.Args.VirtualNetworkId,
		}, nil
	}
	out := new(resourcespb.SubnetResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}
	out.VirtualNetworkId = r.Args.GetVirtualNetworkId()
	out.AzureOutputs = &resourcespb.SubnetAzureOutputs{}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[subnet.AzureSubnet](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.CidrBlock = stateResource.AddressPrefixes[0]
		out.AzureOutputs.SubnetId = stateResource.ResourceId
	} else {
		statuses["azure_subnet"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	// TODO: check if subnet is associated with route table defined in the VN

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureSubnet) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var azResources []output.TfBlock
	azSubnet := subnet.AzureSubnet{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			Name:              r.Args.Name,
			ResourceGroupName: GetResourceGroupName(r.VirtualNetwork.Args.GetCommonParameters().GetResourceGroupId()),
		},
		AddressPrefixes:    []string{r.Args.CidrBlock},
		VirtualNetworkName: fmt.Sprintf("%s.%s.name", virtual_network.AzureResourceName, r.VirtualNetwork.ResourceId),
	}
	azSubnet.ServiceEndpoints = getServiceEndpointSubnetReferences(ctx, r.Subnet)
	azResources = append(azResources, azSubnet)

	// there must be a better way to do this
	if !checkSubnetRouteTableAssociated(ctx, r.Subnet) {
		subnetId, err := resources.GetMainOutputId(r)
		if err != nil {
			return nil, err
		}
		rtAssociation := route_table_association.AzureRouteTableAssociation{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			SubnetId:     subnetId,
			RouteTableId: AzureVirtualNetwork{r.VirtualNetwork}.GetAssociatedRouteTableId(),
		}
		azResources = append(azResources, rtAssociation)
	}

	return azResources, nil
}

func (r AzureSubnet) GetMainResourceName() (string, error) {
	return subnet.AzureResourceName, nil
}

func getServiceEndpointSubnetReferences(ctx resources.MultyContext, r *types.Subnet) []string {
	const (
		DATABASE = "Microsoft.Sql"
	)

	serviceEndpoints := map[string]bool{}
	if len(resources.GetAllResourcesWithRef(ctx, func(db *types.Database) *types.Subnet { return db.Subnet }, r)) > 0 {
		serviceEndpoints[DATABASE] = true
	}
	return util.SortedKeys(serviceEndpoints)
}

func checkSubnetRouteTableAssociated(ctx resources.MultyContext, r *types.Subnet) bool {
	return len(resources.GetAllResourcesWithRef(ctx, func(rt *types.RouteTableAssociation) *types.Subnet { return rt.Subnet }, r)) > 0
}
