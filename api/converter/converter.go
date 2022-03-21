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

type DatabaseConverter struct {
}

func (v DatabaseConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v DatabaseConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificDatabaseArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	db := types.Database{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name:          arg.Name,
		Engine:        arg.Engine.String(),
		EngineVersion: arg.EngineVersion,
		Storage:       int(arg.StorageMb / 1024),
		Size:          arg.Size.String(),
		DbUsername:    arg.Username,
		DbPassword:    arg.Password,
	}

	for _, subnetId := range arg.SubnetIds {

		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			db.SubnetIds = append(db.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &db,
		ImplicitlyCreated: false,
	}, nil
}

type ObjectStorageConverter struct {
}

func (v ObjectStorageConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v ObjectStorageConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificObjectStorageArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	db := types.ObjectStorage{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &db,
		ImplicitlyCreated: false,
	}, nil
}

type ObjectStorageObjectConverter struct {
}

func (v ObjectStorageObjectConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v ObjectStorageObjectConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificObjectStorageObjectArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	obj := types.ObjectStorageObject{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name:        arg.Name,
		Content:     arg.Content,
		ContentType: arg.ContentType,
		Acl:         arg.Acl.String(),
		Source:      arg.Source,
	}

	if o, ok := otherResources[common_resources.GetResourceIdForCloud(arg.ObjectStorageId, c)]; ok {
		// Connect to vn in the same cloud
		obj.ObjectStorage = o.Resource.(*types.ObjectStorage)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("object storage with id %s not found in %s", arg.ObjectStorageId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &obj,
		ImplicitlyCreated: false,
	}, nil
}

type PublicIpConverter struct {
}

func (v PublicIpConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v PublicIpConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificPublicIpArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	obj := types.PublicIp{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
	}

	if ni, ok := otherResources[common_resources.GetResourceIdForCloud(arg.NetworkInterfaceId, c)]; ok {
		// Connect to vn in the same cloud
		obj.NetworkInterfaceId = ni.Resource.(*types.NetworkInterface)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("network interface with id %s not found in %s", arg.NetworkInterfaceId, c)
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &obj,
		ImplicitlyCreated: false,
	}, nil
}

type KubernetesClusterConverter struct {
}

func (v KubernetesClusterConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (v KubernetesClusterConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificKubernetesClusterArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	kc := types.KubernetesService{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name: arg.Name,
	}

	for _, subnetId := range arg.SubnetIds {

		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			kc.SubnetIds = append(kc.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &kc,
		ImplicitlyCreated: false,
	}, nil
}

type KubernetesNodePoolConverter struct {
}

func (v KubernetesNodePoolConverter) NewArg() proto.Message {
	return &resources.CloudSpecificRouteTableArgs{}
}

func zeroToNil(a int32) *int {
	var result *int
	if a != 0 {
		n := int(a)
		result = &n
	}
	return result
}

func (v KubernetesNodePoolConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) (common_resources.CloudSpecificResource, error) {
	arg := m.(*resources.CloudSpecificKubernetesNodePoolArgs)
	c := cloud_providers.CloudProvider(strings.ToLower(arg.CommonParameters.CloudProvider.String()))
	knp := types.KubernetesServiceNodePool{
		CommonResourceParams: &common_resources.CommonResourceParams{
			ResourceId:      resourceId,
			ResourceGroupId: arg.CommonParameters.ResourceGroupId,
			Location:        strings.ToLower(arg.CommonParameters.Location.String()),
			Clouds:          []string{string(c)},
		},
		Name:              arg.Name,
		IsDefaultPool:     arg.IsDefaultPool,
		StartingNodeCount: zeroToNil(arg.StartingNodeCount),
		MaxNodeCount:      int(arg.MaxNodeCount),
		MinNodeCount:      int(arg.MinNodeCount),
		Labels:            arg.Labels,
		VmSize:            arg.VmSize.String(),
		DiskSizeGiB:       int(arg.DiskSizeGb),
	}

	if kc, ok := otherResources[common_resources.GetResourceIdForCloud(arg.ClusterId, c)]; ok {
		// Connect to vn in the same cloud
		knp.ClusterId = kc.Resource.(*types.KubernetesService)
	} else {
		return common_resources.CloudSpecificResource{}, fmt.Errorf("cluster with id %s not found in %s", arg.ClusterId, c)
	}

	for _, subnetId := range arg.SubnetIds {
		if subnet, ok := otherResources[common_resources.GetResourceIdForCloud(subnetId, c)]; ok {
			// Connect to vn in the same cloud
			knp.SubnetIds = append(knp.SubnetIds, subnet.Resource.(*types.Subnet))
		} else {
			return common_resources.CloudSpecificResource{}, fmt.Errorf("subnet with id %s not found in %s", subnetId, c)
		}
	}

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &knp,
		ImplicitlyCreated: false,
	}, nil
}
