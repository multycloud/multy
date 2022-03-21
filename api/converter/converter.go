package converter

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resources"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
)

type ResourceConverters[Arg proto.Message, OutT proto.Message] interface {
	Convert(resourceId string, request []Arg) OutT
	NewArg() Arg
	Nil() OutT
}

type MultyResourceConverter interface {
	ConvertToMultyResource(resourceId string, arg proto.Message, resources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error)
	NewArg() proto.Message
}

type VnConverter struct {
}

func (v VnConverter) NewArg() proto.Message {
	return &resources.CloudSpecificVirtualNetworkArgs{}
}

func (v VnConverter) ConvertToMultyResource(resourceId string, m proto.Message, _ map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificVirtualNetworkArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	vn := types.VirtualNetwork{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name:      arg.Name,
		CidrBlock: arg.CidrBlock,
	}
	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &vn,
		ImplicitlyCreated: false,
	}, nil
}

type SubnetConverter struct {
}

func (v SubnetConverter) NewArg() proto.Message {
	return &resources.CloudSpecificSubnetArgs{}
}

func (v SubnetConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificSubnetArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	subnet := types.Subnet{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name:             arg.Name,
		CidrBlock:        arg.CidrBlock,
		AvailabilityZone: int(arg.AvailabilityZone),
	}

	if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
		// Connect to vn in the same cloud
		subnet.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &subnet,
		ImplicitlyCreated: false,
	}, nil
}

type NetworkInterfaceConverter struct {
}

func (v NetworkInterfaceConverter) NewArg() proto.Message {
	return &resources.CloudSpecificNetworkInterfaceArgs{}
}

func (v NetworkInterfaceConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificNetworkInterfaceArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	ni := types.NetworkInterface{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
	}

	if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SubnetId, c)]; ok {
		// Connect to subnet in the same cloud
		ni.SubnetId = subnet.Resource.(*types.Subnet)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", arg.SubnetId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &ni,
		ImplicitlyCreated: false,
	}, nil
}

type RouteTableConverter struct {
}

func (v RouteTableConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v RouteTableConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificRouteTableArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	rt := types.RouteTable{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
		Routes: util.MapSliceValues(arg.Routes, func(route *resources.Route) types.RouteTableRoute {
			return types.RouteTableRoute{
				CidrBlock:   route.CidrBlock,
				Destination: strings.ToLower(route.Destination.String()),
			}
		}),
	}

	if arg.VirtualNetworkId != "" {
		if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
			// Connect to vn in the same cloud
			rt.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &rt,
		ImplicitlyCreated: false,
	}, nil
}

type RouteTableAssociationConverter struct {
}

func (v RouteTableAssociationConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableAssociationArgs{}
}

func (v RouteTableAssociationConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificRouteTableAssociationArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	rta := types.RouteTableAssociation{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
	}

	if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(arg.SubnetId, c)]; ok {
		// Connect to subnet in the same cloud
		rta.SubnetId = subnet.Resource.(*types.Subnet)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", arg.SubnetId, c)
	}

	if rt, ok := otherResources[common_resources.GetResourceIdForCloud(arg.RouteTableId, c)]; ok {
		// Connect to subnet in the same cloud
		rta.RouteTableId = rt.Resource.(*types.RouteTable)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("route table with id %s not found in %s", arg.RouteTableId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &rta,
		ImplicitlyCreated: false,
	}, nil
}

type NetworkSecurityGroupConverter struct {
}

func (v NetworkSecurityGroupConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v NetworkSecurityGroupConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificNetworkSecurityGroupArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	nsg := types.NetworkSecurityGroup{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
		Rules: util.MapSliceValues(arg.Rules, func(rule *resources.NetworkSecurityRule) types.RuleType {
			return types.RuleType{
				Protocol:  rule.Protocol,
				Priority:  int(rule.Priority),
				FromPort:  convertPort(rule.PortRange.From),
				ToPort:    convertPort(rule.PortRange.To),
				CidrBlock: rule.CidrBlock,
				Direction: convertRuleDirection(rule.Direction),
			}
		}),
	}

	if vn, ok := otherResources[common_resources.GetResourceIdForCloud(arg.VirtualNetworkId, c)]; ok {
		// Connect to vn in the same cloud
		nsg.VirtualNetwork = vn.Resource.(*types.VirtualNetwork)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("virtual network with id %s not found in %s", arg.VirtualNetworkId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &nsg,
		ImplicitlyCreated: false,
	}, nil
}

func convertRuleDirection(direction resources.Direction) string {
	switch direction {
	case resources.Direction_BOTH_DIRECTIONS:
		return "both"
	case resources.Direction_INGRESS:
		return "ingress"
	case resources.Direction_EGRESS:
		return "egress"
	}

	return "unknown"
}

func convertPort(port int32) string {
	if port == 0 {
		return "*"
	}

	return strconv.Itoa(int(port))
}
