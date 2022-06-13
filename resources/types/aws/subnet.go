package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type AwsSubnet struct {
	*types.Subnet
}

func InitSubnet(r *types.Subnet) resources.ResourceTranslator[*resourcespb.SubnetResource] {
	return AwsSubnet{r}
}

func (r AwsSubnet) FromState(state *output.TfState) (*resourcespb.SubnetResource, error) {
	if flags.DryRun {
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
	out := new(resourcespb.SubnetResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}
	out.AvailabilityZone = r.Args.AvailabilityZone
	out.VirtualNetworkId = r.Args.GetVirtualNetworkId()

	stateResource, err := output.GetParsedById[subnet.AwsSubnet](state, r.ResourceId)
	if err != nil {
		return nil, err
	}
	out.Name = stateResource.AwsResource.Tags["Name"]
	out.CidrBlock = stateResource.CidrBlock
	return out, nil
}

func (r AwsSubnet) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	location := r.VirtualNetwork.GetLocation()
	az, err := common.GetAvailabilityZone(location, int(r.Args.AvailabilityZone), r.GetCloud())
	if err != nil {
		return nil, err
	}
	awsSubnet := subnet.AwsSubnet{
		AwsResource:      common.NewAwsResource(r.ResourceId, r.Args.Name),
		CidrBlock:        r.Args.CidrBlock,
		VpcId:            fmt.Sprintf("%s.%s.id", virtual_network.AwsResourceName, r.VirtualNetwork.ResourceId),
		AvailabilityZone: az,
	}
	// This flag needs to be set so that eks nodes can connect to the kubernetes cluster
	// https://aws.amazon.com/blogs/containers/upcoming-changes-to-ip-assignment-for-eks-managed-node-groups/
	// How to tell if this subnet is private?
	if len(resources.GetAllResourcesWithRef(ctx, func(k *types.KubernetesNodePool) *types.Subnet { return k.Subnet }, r.Subnet)) > 0 {
		awsSubnet.MapPublicIpOnLaunch = true
	}
	if len(resources.GetAllResourcesWithRef(ctx, func(k *types.KubernetesCluster) *types.Subnet { return k.DefaultNodePool.Subnet }, r.Subnet)) > 0 {
		awsSubnet.MapPublicIpOnLaunch = true
	}
	return []output.TfBlock{awsSubnet}, nil
}

func (r AwsSubnet) GetMainResourceName() (string, error) {
	return subnet.AwsResourceName, nil
}
