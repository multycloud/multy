package aws_resources

import (
	"fmt"
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
