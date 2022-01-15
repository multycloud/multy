package types

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/route_table_association"
	"multy-go/resources/output/subnet"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

/*
Notes:
Azure: New subnets will be associated with a default route table to block traffic to internet
*/

type Subnet struct {
	*resources.CommonResourceParams
	Name             string `hcl:"name"`
	CidrBlock        string `hcl:"cidr_block"`
	VirtualNetworkId string `hcl:"virtual_network_id"`
	AvailabilityZone int    `hcl:"availability_zone,optional"`
}

func (s *Subnet) getMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return subnet.AwsResourceName
	case common.AZURE:
		return subnet.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}

func (s *Subnet) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{
		"id": cty.StringVal(s.getOutputId(cloud)),
	}
}

func (s *Subnet) getOutputId(cloud common.CloudProvider) string {
	return fmt.Sprintf("${%s.%s.id}", s.getMainResourceName(cloud), s.GetTfResourceId(cloud))
}

func (s *Subnet) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	var virtualNetwork *VirtualNetwork
	if vn, err := ctx.GetResource(s.VirtualNetworkId); err != nil {
		s.LogFatal(s.ResourceId, "virtual_network_id", err.Error())
	} else {
		virtualNetwork = vn.Resource.(*VirtualNetwork)
	}
	if cloud == common.AWS {
		return []interface{}{subnet.AwsSubnet{
			AwsResource: common.AwsResource{
				ResourceName: subnet.AwsResourceName,
				ResourceId:   s.GetTfResourceId(cloud),
				Tags:         map[string]string{"Name": s.Name},
			},
			CidrBlock:        s.CidrBlock,
			VpcId:            virtualNetwork.GetVirtualNetworkId(cloud),
			AvailabilityZone: common.GetAvailabilityZone(ctx.Location, s.AvailabilityZone, cloud),
		}}
	} else if cloud == common.AZURE {
		var azResources []interface{}
		azSubnet := subnet.AzureSubnet{
			AzResource: common.AzResource{
				ResourceName:      subnet.AzureResourceName,
				ResourceId:        s.GetTfResourceId(cloud),
				Name:              s.Name,
				ResourceGroupName: rg.GetResourceGroupName(s.ResourceGroupId, cloud),
			},
			AddressPrefixes:    []string{s.CidrBlock},
			VirtualNetworkName: virtualNetwork.GetVirtualNetworkName(cloud),
		}
		azSubnet.ServiceEndpoints = getServiceEndpointSubnetReferences(ctx, s.getOutputId(cloud))
		azResources = append(azResources, azSubnet)

		// there must be a better way to do this
		if !checkSubnetRouteTableAssociated(ctx, s.getOutputId(cloud)) {
			rt := virtualNetwork.GetAssociatedRouteTableId(cloud)
			rtAssociation := route_table_association.AzureRouteTableAssociation{
				AzResource: common.AzResource{
					ResourceName: route_table_association.AzureResourceName,
					ResourceId:   s.GetTfResourceId(cloud),
				},
				SubnetId:     s.getOutputId(cloud),
				RouteTableId: rt,
			}
			azResources = append(azResources, rtAssociation)
		}

		return azResources
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (s *Subnet) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: "aws_subnet", common.AZURE: "azurerm_subnet"}
	return fmt.Sprintf("%s.%s.id", types[cloud], s.GetTfResourceId(cloud))
}

func getServiceEndpointSubnetReferences(ctx resources.MultyContext, id string) []string {
	const (
		DATABASE = "Microsoft.Sql"
	)

	var serviceEndpoints []string
	for _, resource := range ctx.Resources {
		switch resource.Resource.(type) {
		case *Database:
			r := resource.Resource.(*Database)
			if common.StringInSlice(id, r.SubnetIds) {
				serviceEndpoints = append(serviceEndpoints, DATABASE)
			}
		}
	}
	return serviceEndpoints
}

func checkSubnetRouteTableAssociated(ctx resources.MultyContext, sId string) bool {
	for _, resource := range ctx.Resources {
		switch resource.Resource.(type) {
		case *RouteTableAssociation:
			r := resource.Resource.(*RouteTableAssociation)
			if sId == r.SubnetId {
				return true
			}
		}
	}
	return false
}

func (s *Subnet) Validate(ctx resources.MultyContext) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.CidrBlock valid CIDR { return false }
	//if vn.AvailbilityZone valid { return false }
	if len(s.CidrBlock) == 0 { // max len?
		s.LogFatal(s.ResourceId, "cidr_block", fmt.Sprintf("%s cidr_block length is invalid", s.ResourceId))
	}

	return
}
