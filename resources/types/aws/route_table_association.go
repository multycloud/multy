package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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

func (r AwsRouteTableAssociation) FromState(state *output.TfState) (*resourcespb.RouteTableAssociationResource, error) {
	return &resourcespb.RouteTableAssociationResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		SubnetId:     r.Args.SubnetId,
		RouteTableId: r.Args.RouteTableId,
	}, nil
}

func (r AwsRouteTableAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rtId, err := resources.GetMainOutputId(AwsRouteTable{r.RouteTable})
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(AwsSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		route_table_association.AwsRouteTableAssociation{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.Subnet.ResourceId},
			},
			RouteTableId: rtId,
			SubnetId:     subnetId,
		},
	}, nil
}

func (r AwsRouteTableAssociation) GetMainResourceName() (string, error) {
	return virtual_network.AwsResourceName, nil
}
