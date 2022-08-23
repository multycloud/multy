package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
)

type GcpRouteTableAssociation struct {
	*types.RouteTableAssociation
}

func InitRouteTableAssociation(vn *types.RouteTableAssociation) resources.ResourceTranslator[*resourcespb.RouteTableAssociationResource] {
	return GcpRouteTableAssociation{vn}
}

func (r GcpRouteTableAssociation) FromState(_ *output.TfState, _ *output.TfPlan) (*resourcespb.RouteTableAssociationResource, error) {
	return &resourcespb.RouteTableAssociationResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		SubnetId:     r.Args.SubnetId,
		RouteTableId: r.Args.RouteTableId,
	}, nil
}

func (r GcpRouteTableAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return nil, nil
}

func (r GcpRouteTableAssociation) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("route table associations don't exist in gcp")
}
