package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type AwsVirtualNetwork struct {
	*types.VirtualNetwork
}

func InitVirtualNetwork(vn *types.VirtualNetwork) resources.ResourceTranslator[*resourcespb.VirtualNetworkResource] {
	return AwsVirtualNetwork{vn}
}

func (r AwsVirtualNetwork) FromState(state *output.TfState) (*resourcespb.VirtualNetworkResource, error) {
	if flags.DryRun {
		return &resourcespb.VirtualNetworkResource{
			CommonParameters: &commonpb.CommonResourceParameters{
				ResourceId:      r.GetResourceId(),
				ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
				Location:        r.Args.CommonParameters.Location,
				CloudProvider:   r.Args.CommonParameters.CloudProvider,
				NeedsUpdate:     false,
			},
			Name:      r.Args.Name,
			CidrBlock: r.Args.CidrBlock,
		}, nil
	}
	out := new(resourcespb.VirtualNetworkResource)
	out.CommonParameters = &commonpb.CommonResourceParameters{
		ResourceId:      r.GetResourceId(),
		ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
		Location:        r.Args.CommonParameters.Location,
		CloudProvider:   r.GetCloud(),
		NeedsUpdate:     false,
	}

	id, err := resources.GetMainOutputRef(r)
	if err != nil {
		return nil, err
	}

	stateResource, err := output.GetParsed[virtual_network.AwsVpc](state, id)
	if err != nil {
		return nil, err
	}
	out.Name = stateResource.AwsResource.Tags["Name"]
	out.CidrBlock = stateResource.CidrBlock

	return out, nil
}

func (r AwsVirtualNetwork) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	vpc := virtual_network.AwsVpc{
		AwsResource:        common.NewAwsResource(r.GetResourceId(), r.Args.Name),
		CidrBlock:          r.Args.CidrBlock,
		EnableDnsHostnames: true,
	}
	// TODO make conditional on route_table_association with Internet Destination
	igw := virtual_network.AwsInternetGateway{
		AwsResource: common.NewAwsResource(r.GetResourceId(), r.Args.Name),
		VpcId:       fmt.Sprintf("%s.%s.id", virtual_network.AwsResourceName, r.ResourceId),
	}
	allowAllSgRule := []network_security_group.AwsSecurityGroupRule{{
		Protocol: "-1",
		FromPort: 0,
		ToPort:   0,
		Self:     true,
	}}
	sg := network_security_group.AwsDefaultSecurityGroup{
		AwsResource: common.NewAwsResource(r.GetResourceId(), r.Args.Name),
		VpcId:       fmt.Sprintf("%s.%s.id", virtual_network.AwsResourceName, r.ResourceId),
		Ingress:     allowAllSgRule,
		Egress:      allowAllSgRule,
	}
	return []output.TfBlock{
		vpc,
		igw,
		sg,
	}, nil
}

func (r AwsVirtualNetwork) GetMainResourceName() (string, error) {
	return virtual_network.AwsResourceName, nil
}

func (r AwsVirtualNetwork) GetAssociatedInternetGateway() string {
	return fmt.Sprintf("%s.%s.id", virtual_network.AwsInternetGatewayName, r.ResourceId)
}
