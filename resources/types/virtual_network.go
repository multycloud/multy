package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/route_table"
	"github.com/multycloud/multy/resources/output/virtual_network"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

/*
Virtual network traffic is only internal
AWS: Default security group defaults to allow all traffic to  mirror Azure
	IGW created by default (to be changed)
AZ: Route table created to restrict traffic on vnet

*/

type VirtualNetwork struct {
	resources.ResourceWithId[*resourcespb.VirtualNetworkArgs]
}

func NewVirtualNetwork(resourceId string, vn *resourcespb.VirtualNetworkArgs, _ resources.Resources) (*VirtualNetwork, error) {
	return &VirtualNetwork{
		ResourceWithId: resources.ResourceWithId[*resourcespb.VirtualNetworkArgs]{
			ResourceId: resourceId,
			Args:       vn,
		},
	}, nil
}

func (r *VirtualNetwork) FromState(state *output.TfState) (*resourcespb.VirtualNetworkResource, error) {
	out := new(resourcespb.VirtualNetworkResource)

	id, err := resources.GetMainOutputRef(r)
	if err != nil {
		return nil, err
	}

	switch r.GetCloud() {
	case common.AWS:
		stateResource, err := output.GetParsed[virtual_network.AwsVpc](state, id)
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.AwsResource.Tags["Name"]
		out.CidrBlock = stateResource.CidrBlock
		out.CommonParameters = &commonpb.CommonResourceParameters{
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.GetCloud(),
			NeedsUpdate:     false,
		}
	case common.AZURE:
		stateResource, err := output.GetParsed[virtual_network.AzureVnet](state, id)
		if err != nil {
			return nil, err
		}
		out.Name = stateResource.Name
		out.CidrBlock = stateResource.AddressSpace[0]
		out.CommonParameters = &commonpb.CommonResourceParameters{
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.GetCloud(),
			NeedsUpdate:     false,
		}
	}

	return out, nil
}

func (r *VirtualNetwork) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		vpc := virtual_network.AwsVpc{
			AwsResource:        common.NewAwsResource(r.ResourceId, r.Args.Name),
			CidrBlock:          r.Args.CidrBlock,
			EnableDnsHostnames: true,
		}
		// TODO make conditional on route_table_association with Internet Destination
		igw := virtual_network.AwsInternetGateway{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			VpcId:       r.GetVirtualNetworkId(),
		}
		allowAllSgRule := []network_security_group.AwsSecurityGroupRule{{
			Protocol:   "-1",
			FromPort:   0,
			ToPort:     0,
			CidrBlocks: []string{"0.0.0.0/0"},
			Self:       true,
		}}
		sg := network_security_group.AwsDefaultSecurityGroup{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			VpcId:       r.GetVirtualNetworkId(),
			Ingress:     allowAllSgRule,
			Egress:      allowAllSgRule,
		}
		return []output.TfBlock{
			vpc,
			igw,
			sg,
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return []output.TfBlock{virtual_network.AzureVnet{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
				r.GetCloudSpecificLocation(),
			),
			AddressSpace: []string{r.Args.CidrBlock},
		}, route_table.AzureRouteTable{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
				r.GetCloudSpecificLocation(),
			),
			Routes: []route_table.AzureRouteTableRoute{{
				Name:          "local",
				AddressPrefix: "0.0.0.0/0",
				NextHopType:   "VnetLocal",
			}},
		}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *VirtualNetwork) GetVirtualNetworkId() string {
	t, _ := r.GetMainResourceName()
	return fmt.Sprintf("%s.%s.id", t, r.ResourceId)
}

func (r *VirtualNetwork) GetVirtualNetworkName() string {
	t, _ := r.GetMainResourceName()
	return fmt.Sprintf("%s.%s.name", t, r.ResourceId)
}

func (r *VirtualNetwork) GetAssociatedRouteTableId() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return fmt.Sprintf("${%s.%s.id}", route_table.AzureResourceName, r.ResourceId), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *VirtualNetwork) GetAssociatedInternetGateway() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return fmt.Sprintf("%s.%s.id", virtual_network.AwsInternetGatewayName, r.ResourceId), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

//// TODO validate commonparams
func (r *VirtualNetwork) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	//if r.Name contains not letters,numbers,_,- { return false }
	//if r.Name length? { return false }
	//if r.CidrBlock valid CIDR { return false }
	if len(r.Args.CidrBlock) == 0 { // max len?
		errs = append(errs, validate.ValidationError{
			ErrorMessage: "cidr_block length is invalid",
			ResourceId:   r.ResourceId,
			FieldName:    "cidr_block",
		})
	}
	return errs
}

func (r *VirtualNetwork) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return virtual_network.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return virtual_network.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

/*
Virtual Network is a private address space when resources can be placed.

By default, resources inside `virtual_network` cannot access the internet. To enable internet access look at`route_table`

Mapping:
AWS:
`aws_vpc` - VPC
`aws_internet_gateway` - Internet Gateway attached to VPC. Default route table does not route traffic to outside VPC

Azure:
`azurerm_virtual_network` - Vnet
`azurerm_route_table` - Default route table to block internet access

Inputs:

cidr_block `[]string` (required): Address range of virtual network.
name `string` (required): Name of virtual network. Read about multy naming [here](xxx)
location `string` (optional): Region to deploy resource into. See available locations [here](xxx)
clouds `[]string` (optional): Clouds where services will be deployed to. Learn about multy clouds [here](xxx)
*/
