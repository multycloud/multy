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

func (r AwsVirtualNetwork) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VirtualNetworkResource, error) {
	if flags.DryRun {
		return &resourcespb.VirtualNetworkResource{
			CommonParameters: &commonpb.CommonResourceParameters{
				ResourceId:      r.GetResourceId(),
				ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
				Location:        r.Args.CommonParameters.Location,
				CloudProvider:   r.Args.CommonParameters.CloudProvider,
				NeedsUpdate:     false,
			},
			Name:        r.Args.Name,
			CidrBlock:   r.Args.CidrBlock,
			GcpOverride: r.Args.GcpOverride,
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
	out.GcpOverride = r.Args.GcpOverride
	out.AwsOutputs = &resourcespb.VirtualNetworkAwsOutputs{}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[virtual_network.AwsVpc](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.AwsResource.Tags["Name"]
		out.CidrBlock = stateResource.CidrBlock
		out.AwsOutputs.VpcId = stateResource.AwsResource.ResourceId
		output.AddToStatuses(statuses, "aws_vpc", output.MaybeGetPlannedChageById[virtual_network.AwsVpc](plan, r.ResourceId))
	} else {
		statuses["aws_vpc"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[virtual_network.AwsInternetGateway](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.InternetGatewayId = stateResource.ResourceId
		output.AddToStatuses(statuses, "aws_internet_gateway", output.MaybeGetPlannedChageById[virtual_network.AwsInternetGateway](plan, r.ResourceId))
	} else {
		statuses["aws_internet_gateway"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.AwsDefaultSecurityGroup](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.DefaultSecurityGroupId = stateResource.ResourceId
		output.AddToStatuses(statuses, "aws_default_security_group", output.MaybeGetPlannedChageById[network_security_group.AwsDefaultSecurityGroup](plan, r.ResourceId))
	} else {
		statuses["aws_default_security_group"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}

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
