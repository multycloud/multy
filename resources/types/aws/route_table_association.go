package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type AwsRouteTableAssociation struct {
	*types.RouteTableAssociation
}

func InitRouteTableAssociation(vn *types.RouteTableAssociation) resources.ResourceTranslator[*resourcespb.RouteTableAssociationResource] {
	return AwsRouteTableAssociation{vn}
}

func (r AwsRouteTableAssociation) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.RouteTableAssociationResource, error) {
	out := &resourcespb.RouteTableAssociationResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		SubnetId:     r.Args.SubnetId,
		RouteTableId: r.Args.RouteTableId,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}
	out.AwsOutputs = &resourcespb.RouteTableAssociationAwsOutputs{RouteTableAssociationIdByAvailabilityZone: map[string]string{}}

	azArray := common.AVAILABILITY_ZONES[r.Subnet.VirtualNetwork.GetLocation()][r.GetCloud()]
	for i, zone := range azArray {
		resourceId := fmt.Sprintf("%s-%d", r.ResourceId, i+1)
		if stateResource, exists, err := output.MaybeGetParsedById[route_table_association.AwsRouteTableAssociation](state, resourceId); exists {
			if err != nil {
				return nil, err
			}
			out.AwsOutputs.RouteTableAssociationIdByAvailabilityZone[zone] = stateResource.ResourceId
		} else {
			statuses[fmt.Sprintf("aws_route_table_association_%d", i)] = commonpb.ResourceStatus_NEEDS_CREATE
		}
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AwsRouteTableAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rtId, err := resources.GetMainOutputId(AwsRouteTable{r.RouteTable})
	if err != nil {
		return nil, err
	}
	var result []output.TfBlock
	for i, subnetId := range (AwsSubnet{r.Subnet}.GetSubnetIds()) {
		resourceId := fmt.Sprintf("%s-%d", r.ResourceId, i+1)
		result = append(result, route_table_association.AwsRouteTableAssociation{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: resourceId},
			},
			RouteTableId: rtId,
			SubnetId:     subnetId,
		})
	}
	return result, nil
}

func (r AwsRouteTableAssociation) GetMainResourceName() (string, error) {
	return virtual_network.AwsResourceName, nil
}
