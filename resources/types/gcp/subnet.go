package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type GcpSubnet struct {
	*types.Subnet
}

func InitSubnet(r *types.Subnet) resources.ResourceTranslator[*resourcespb.SubnetResource] {
	return GcpSubnet{r}
}

func (r GcpSubnet) FromState(_ *output.TfState) (*resourcespb.SubnetResource, error) {
	return &resourcespb.SubnetResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:             r.Args.Name,
		CidrBlock:        r.Args.CidrBlock,
		AvailabilityZone: r.Args.AvailabilityZone,
		VirtualNetworkId: r.Args.VirtualNetworkId,
	}, nil
}

func (r GcpSubnet) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{&subnet.GoogleComputeSubnetwork{
		GcpResource:           common.NewGcpResource(r.ResourceId, r.Args.Name),
		IpCidrRange:           r.Args.CidrBlock,
		PrivateIpGoogleAccess: true,
		Network:               fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
	}}, nil
}

func (r GcpSubnet) getNetworkTag() string {
	return fmt.Sprintf("subnet-%s", r.Args.Name)
}

func (r GcpSubnet) GetMainResourceName() (string, error) {
	return output.GetResourceName(subnet.GoogleComputeSubnetwork{}), nil
}
