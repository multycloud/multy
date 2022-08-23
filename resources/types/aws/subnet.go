package aws_resources

import (
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"net"
	"strings"
)

type AwsSubnet struct {
	*types.Subnet
}

func InitSubnet(r *types.Subnet) resources.ResourceTranslator[*resourcespb.SubnetResource] {
	return AwsSubnet{r}
}

func (r AwsSubnet) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.SubnetResource, error) {
	if flags.DryRun {
		return &resourcespb.SubnetResource{
			CommonParameters: &commonpb.CommonChildResourceParameters{
				ResourceId:  r.ResourceId,
				NeedsUpdate: false,
			},
			Name:             r.Args.Name,
			CidrBlock:        r.Args.CidrBlock,
			VirtualNetworkId: r.Args.VirtualNetworkId,
		}, nil
	}
	out := new(resourcespb.SubnetResource)
	out.CommonParameters = &commonpb.CommonChildResourceParameters{
		ResourceId:  r.ResourceId,
		NeedsUpdate: false,
	}
	out.VirtualNetworkId = r.Args.GetVirtualNetworkId()
	out.Name = r.Args.Name
	out.CidrBlock = r.Args.CidrBlock
	out.AwsOutputs = &resourcespb.SubnetAwsOutputs{
		SubnetIdByAvailabilityZone: map[string]string{},
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}
	if azArray, ok := common.AVAILABILITY_ZONES[r.VirtualNetwork.GetLocation()][r.GetCloud()]; ok {
		for i, zone := range azArray {
			subnetI := fmt.Sprintf("aws_subnet_%d", i)
			stateResource, exists, err := output.MaybeGetParsedById[subnet.AwsSubnet](state, r.getSubnetResourceId(i))
			if err != nil {
				return out, err
			}
			if !exists {
				statuses[subnetI] = commonpb.ResourceStatus_NEEDS_CREATE
				continue
			}
			output.AddToStatuses(statuses, subnetI, output.MaybeGetPlannedChageById[subnet.AwsSubnet](plan, r.getSubnetResourceId(i)))

			if stateResource.AvailabilityZone != zone {
				continue
			}
			out.AwsOutputs.SubnetIdByAvailabilityZone[stateResource.AvailabilityZone] = stateResource.ResourceId
			name := stateResource.Tags["Name"]
			if name != r.getSubnetResourceName(i) {
				suffix := fmt.Sprintf("-%d", i+1)
				if strings.HasSuffix(name, suffix) {
					out.Name = strings.TrimSuffix(name, suffix)
				} else if _, ok := statuses[subnetI]; !ok {
					statuses[subnetI] = commonpb.ResourceStatus_NEEDS_UPDATE
				}
			}
		}
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AwsSubnet) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var out []output.TfBlock
	azArray := common.AVAILABILITY_ZONES[r.VirtualNetwork.GetLocation()][r.GetCloud()]
	subnetCidrBlocks, err := splitSubnet(r.Args.CidrBlock, len(azArray))
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage(
			fmt.Sprintf("Unable to split subnet %s into %d availability zones: %s", r.Args.CidrBlock, len(azArray), err.Error()),
			err)
	}

	for i, az := range azArray {
		awsSubnet := subnet.AwsSubnet{
			AwsResource:      common.NewAwsResource(r.getSubnetResourceId(i), r.getSubnetResourceName(i)),
			CidrBlock:        subnetCidrBlocks[i],
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

		out = append(out, awsSubnet)
	}

	return out, nil
}

func (r AwsSubnet) getSubnetResourceId(i int) string {
	return fmt.Sprintf("%s-%d", r.ResourceId, i+1)
}

func (r AwsSubnet) getSubnetResourceName(i int) string {
	return fmt.Sprintf("%s-%d", r.Args.Name, i+1)
}

func (r AwsSubnet) GetSubnetIds() []string {
	var out []string
	azArray := common.AVAILABILITY_ZONES[r.VirtualNetwork.GetLocation()][r.GetCloud()]
	for i := range azArray {
		out = append(out, fmt.Sprintf("%s.%s.id", output.GetResourceName(subnet.AwsSubnet{}), r.getSubnetResourceId(i)))
	}
	return out
}

func (r AwsSubnet) GetSubnetId(zone int32) (string, error) {
	if zone == 0 {
		return "", fmt.Errorf("zone 0 is invalid")
	}
	_, err := common.GetAvailabilityZone(r.VirtualNetwork.GetLocation(), int(zone), r.GetCloud())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s.id", output.GetResourceName(subnet.AwsSubnet{}), r.getSubnetResourceId(int(zone-1))), nil
}

func splitSubnet(cidrBlock string, numSubnets int) ([]string, error) {
	if numSubnets == 0 {
		return nil, nil
	}
	if numSubnets == 1 {
		return []string{cidrBlock}, nil
	}
	_, ipnet, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return nil, err
	}

	subnet1, err := cidr.Subnet(ipnet, 1, 0)
	if err != nil {
		return nil, err
	}
	subnet2, err := cidr.Subnet(ipnet, 1, 1)
	if err != nil {
		return nil, err
	}

	subnet1NumSubnets := numSubnets / 2
	subnet2NumSubnets := numSubnets - subnet1NumSubnets
	var out []string
	subnet1Subnets, err := splitSubnet(subnet1.String(), subnet1NumSubnets)
	if err != nil {
		return nil, err
	}
	subnet2Subnets, err := splitSubnet(subnet2.String(), subnet2NumSubnets)
	if err != nil {
		return nil, err
	}

	out = append(out, subnet1Subnets...)
	out = append(out, subnet2Subnets...)
	return out, nil
}

func (r AwsSubnet) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("there's many subnets")
}
