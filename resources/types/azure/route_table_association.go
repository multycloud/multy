package azure_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/types"
)

type AzureRouteTableAssociation struct {
	*types.RouteTableAssociation
}

func InitRouteTableAssociation(vn *types.RouteTableAssociation) resources.ResourceTranslator[*resourcespb.RouteTableAssociationResource] {
	return AzureRouteTableAssociation{vn}
}

func (r AzureRouteTableAssociation) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.RouteTableAssociationResource, error) {
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
	if _, exists, _ := output.MaybeGetParsedById[route_table_association.AzureRouteTableAssociation](state, r.getRtaId()); !exists {
		statuses["azure_route_table_association"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureRouteTableAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rtId, err := resources.GetMainOutputId(AzureRouteTable{r.RouteTable})
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(AzureSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		route_table_association.AzureRouteTableAssociation{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.getRtaId()},
			},
			RouteTableId: rtId,
			SubnetId:     subnetId,
		},
	}, nil
}

func (r AzureRouteTableAssociation) getRtaId() string {
	// Here we use the subnet id so that it is the same as the one that is created by default in the subnet.
	// This ensures that if a RTA is created after the default RTA is created, they will have the same ID and
	// terraform will either update it in place or destroy it before creating it.
	return r.Subnet.ResourceId
}
func (r AzureRouteTableAssociation) GetMainResourceName() (string, error) {
	return route_table_association.AzureResourceName, nil
}
