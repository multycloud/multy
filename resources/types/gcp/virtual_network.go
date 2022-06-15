package gcp_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type GcpVirtualNetwork struct {
	*types.VirtualNetwork
}

func InitVirtualNetwork(vn *types.VirtualNetwork) resources.ResourceTranslator[*resourcespb.VirtualNetworkResource] {
	return GcpVirtualNetwork{vn}
}

func (r GcpVirtualNetwork) FromState(_ *output.TfState) (*resourcespb.VirtualNetworkResource, error) {
	return &resourcespb.VirtualNetworkResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.GetCloud(),
		},
		Name:      r.Args.Name,
		CidrBlock: r.Args.CidrBlock,
	}, nil
}

func (r GcpVirtualNetwork) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{&virtual_network.GoogleComputeNetwork{
		GcpResource:                 common.NewGcpResource(r.ResourceId, r.Args.Name),
		RoutingMode:                 "REGIONAL",
		Description:                 "Managed by Multy",
		AutoCreateSubnetworks:       false,
		DeleteDefaultRoutesOnCreate: true,
	}}, nil
}

func (r GcpVirtualNetwork) GetMainResourceName() (string, error) {
	return virtual_network.GcpResourceName, nil
}
