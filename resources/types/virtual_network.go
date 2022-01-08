package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/network_security_group"
	"multy-go/resources/output/route_table"
	"multy-go/resources/output/virtual_network"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"

	"github.com/zclconf/go-cty/cty"
)

/*
Virtual network traffic is only internal
AWS: Default security group defaults to allow all traffic to  mirror Azure
	IGW created by default (to be changed)
AZ: Route table created to restrict traffic on vnet

*/

type VirtualNetwork struct {
	*resources.CommonResourceParams `hcl:",block"`
	Name                            string `hcl:"name"`
	CidrBlock                       string `hcl:"cidr_block"`
}

type VirtualNetworkOutput struct {
	*resources.CommonResourceOutputs
	CidrBlock string
}

func (vn *VirtualNetwork) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	if cloud == common.AWS {
		vpc := virtual_network.AwsVpc{
			AwsResource: common.AwsResource{
				ResourceName: virtual_network.AwsResourceName,
				ResourceId:   vn.GetTfResourceId(cloud),
				Tags:         map[string]string{"Name": vn.Name},
			},
			CidrBlock:          vn.CidrBlock,
			EnableDnsHostnames: true,
		}
		// TODO make conditional on route_table_association with Internet Destination
		igw := virtual_network.AwsInternetGateway{
			AwsResource: common.AwsResource{
				ResourceName: virtual_network.AwsInternetGatewayName,
				ResourceId:   vn.GetTfResourceId(cloud),
				Tags:         map[string]string{"Name": vn.Name},
			},
			VpcId: vn.GetVirtualNetworkId(common.AWS),
		}
		allowAllSgRule := []network_security_group.AwsSecurityGroupRule{{
			Protocol:   "-1",
			FromPort:   0,
			ToPort:     0,
			CidrBlocks: "0.0.0.0/0",
			Self:       true,
		}}
		sg := network_security_group.AwsDefaultSecurityGroup{
			AwsResource: common.AwsResource{
				ResourceName: network_security_group.AwsDefaultSecurityGroupResourceName,
				ResourceId:   vn.GetTfResourceId(cloud),
				Tags:         map[string]string{"Name": vn.Name},
			},
			VpcId:   vn.GetVirtualNetworkId(common.AWS),
			Ingress: allowAllSgRule,
			Egress:  allowAllSgRule,
		}
		return []interface{}{
			vpc,
			igw,
			sg,
		}
	} else if cloud == common.AZURE {
		return []interface{}{virtual_network.AzureVnet{
			AzResource: common.AzResource{
				ResourceName:      virtual_network.AzureResourceName,
				ResourceId:        vn.GetTfResourceId(cloud),
				Name:              vn.Name,
				ResourceGroupName: rg.GetResourceGroupName(vn.ResourceGroupId, cloud),
				Location:          ctx.GetLocationFromCommonParams(vn.CommonResourceParams, cloud),
			},
			AddressSpace: []string{vn.CidrBlock},
		}, route_table.AzureRouteTable{
			AzResource: common.AzResource{
				ResourceName:      route_table.AzureResourceName,
				ResourceId:        vn.GetTfResourceId(cloud),
				ResourceGroupName: rg.GetResourceGroupName(vn.ResourceGroupId, cloud),
				Name:              vn.Name,
				Location:          ctx.GetLocationFromCommonParams(vn.CommonResourceParams, cloud),
			},
			Routes: []route_table.AzureRouteTableRoute{{
				Name:          "local",
				AddressPrefix: "0.0.0.0/0",
				NextHopType:   "VnetLocal",
			}},
		}}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (vn *VirtualNetwork) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	if cloud == common.AWS {
		return map[string]cty.Value{
			"test": cty.StringVal(fmt.Sprintf("%s.id", vn.GetTfResourceId(cloud))),
		}
	} else if cloud == common.AZURE {
		return map[string]cty.Value{
			"test": cty.StringVal(fmt.Sprintf("%s.id", vn.GetTfResourceId(cloud))),
		}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (vn *VirtualNetwork) GetVirtualNetworkId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: "aws_vpc", common.AZURE: "azurerm_virtual_network"}
	return fmt.Sprintf("%s.%s.id", types[cloud], vn.GetTfResourceId(cloud))
}

func (vn *VirtualNetwork) GetVirtualNetworkName(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: "aws_vpc", common.AZURE: "azurerm_virtual_network"}
	return fmt.Sprintf("%s.%s.name", types[cloud], vn.GetTfResourceId(cloud))
}

func (vn *VirtualNetwork) GetDefaultNetworkAclId(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return fmt.Sprintf("aws_vpc.%s.default_network_acl_id", vn.GetTfResourceId(common.AWS))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (vn *VirtualNetwork) GetAssociatedRouteTableId(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("%s.%s.id", route_table.AzureResourceName, vn.GetTfResourceId(common.AZURE))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (vn *VirtualNetwork) GetAssociatedInternetGateway(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return fmt.Sprintf("%s.%s.id", virtual_network.AwsInternetGatewayName, vn.GetTfResourceId(common.AWS))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

// TODO validate commonparams
func (vn *VirtualNetwork) Validate(ctx resources.MultyContext) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.CidrBlock valid CIDR { return false }
	if len(vn.CidrBlock) == 0 { // max len?
		vn.LogFatal(vn.ResourceId, "cidr_block", "cidr_block length is invalid")
	}
	return
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
