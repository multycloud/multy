package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/output/subnet"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

/*
Notes:
Azure: New subnets will be associated with a default route table to block traffic to internet
*/

type Subnet struct {
	*resources.CommonResourceParams
	Name             string          `hcl:"name"`
	CidrBlock        string          `hcl:"cidr_block"`
	VirtualNetwork   *VirtualNetwork `mhcl:"ref=virtual_network"`
	AvailabilityZone int             `hcl:"availability_zone,optional"`
}

func (s *Subnet) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		location := s.VirtualNetwork.Location
		if location == "" {
			location = ctx.Location
		}
		awsSubnet := subnet.AwsSubnet{
			AwsResource:      common.NewAwsResource(s.GetTfResourceId(cloud), s.Name),
			CidrBlock:        s.CidrBlock,
			VpcId:            s.VirtualNetwork.GetVirtualNetworkId(cloud),
			AvailabilityZone: common.GetAvailabilityZone(location, s.AvailabilityZone, cloud),
		}
		// This flag needs to be set so that eks nodes can connect to the kubernetes cluster
		// https://aws.amazon.com/blogs/containers/upcoming-changes-to-ip-assignment-for-eks-managed-node-groups/
		// How to tell if this subnet is private?
		for _, resource := range resources.GetAllResources[*KubernetesServiceNodePool](ctx) {
			for _, sn := range resource.SubnetIds {
				if sn.ResourceId == s.ResourceId {
					awsSubnet.MapPublicIpOnLaunch = true
				}
			}
		}
		return []output.TfBlock{awsSubnet}, nil
	} else if cloud == common.AZURE {
		var azResources []output.TfBlock
		azSubnet := subnet.AzureSubnet{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: s.GetTfResourceId(cloud)},
				Name:              s.Name,
				ResourceGroupName: rg.GetResourceGroupName(s.VirtualNetwork.ResourceGroupId, cloud),
			},
			AddressPrefixes:    []string{s.CidrBlock},
			VirtualNetworkName: s.VirtualNetwork.GetVirtualNetworkName(cloud),
		}
		azSubnet.ServiceEndpoints = getServiceEndpointSubnetReferences(ctx, s.ResourceId)
		azResources = append(azResources, azSubnet)

		// there must be a better way to do this
		if !checkSubnetRouteTableAssociated(ctx, s.ResourceId) {
			rt, err := s.VirtualNetwork.GetAssociatedRouteTableId(cloud)
			if err != nil {
				return nil, err
			}
			subnetId, err := resources.GetMainOutputId(s, cloud)
			if err != nil {
				return nil, err
			}
			rtAssociation := route_table_association.AzureRouteTableAssociation{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: s.GetTfResourceId(cloud)},
				},
				SubnetId:     subnetId,
				RouteTableId: rt,
			}
			azResources = append(azResources, rtAssociation)
		}

		return azResources, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
	return nil, nil
}

func (s *Subnet) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: "aws_subnet", common.AZURE: "azurerm_subnet"}
	return fmt.Sprintf("%s.%s.id", types[cloud], s.GetTfResourceId(cloud))
}

func getServiceEndpointSubnetReferences(ctx resources.MultyContext, id string) []string {
	const (
		DATABASE = "Microsoft.Sql"
	)

	serviceEndpoints := map[string]bool{}
	for _, resource := range resources.GetAllResources[*Database](ctx) {
		for _, sn := range resource.SubnetIds {
			if sn.ResourceId == id {
				serviceEndpoints[DATABASE] = true
			}
		}
	}
	return util.Keys(serviceEndpoints)
}

func checkSubnetRouteTableAssociated(ctx resources.MultyContext, sId string) bool {
	for _, resource := range resources.GetAllResources[*RouteTableAssociation](ctx) {
		if sId == resource.SubnetId.ResourceId {
			return true
		}
	}
	return false
}

func (s *Subnet) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.CidrBlock valid CIDR { return false }
	//if vn.AvailbilityZone valid { return false }
	if len(s.CidrBlock) == 0 { // max len?
		s.NewError("cidr_block", fmt.Sprintf("%s cidr_block length is invalid", s.ResourceId))
	}

	return errs
}

func (s *Subnet) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return subnet.AwsResourceName, nil
	case common.AZURE:
		return subnet.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}

func (s *Subnet) GetLocation(cloud common.CloudProvider, ctx resources.MultyContext) string {
	return s.VirtualNetwork.GetLocation(cloud, ctx)
}
