package converter

import (
	"github.com/multycloud/multy/api/proto/resources"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/protobuf/proto"
	"strings"
)

type ResourceConverters[Arg proto.Message, OutT proto.Message] interface {
	Convert(resourceId string, request []Arg) OutT
	NewArg() Arg
	Nil() OutT
}

type MultyResourceConverter interface {
	ConvertToMultyResource(resourceId string, arg proto.Message, resources map[string]common_resources.CloudSpecificResource) common_resources.CloudSpecificResource
	NewArg() proto.Message
}

type VnConverter struct {
}

func (v VnConverter) NewArg() proto.Message {
	return &resources.CloudSpecificVirtualNetworkArgs{}
}

func (v VnConverter) ConvertToMultyResource(resourceId string, m proto.Message, _ map[string]common_resources.CloudSpecificResource) common_resources.CloudSpecificResource {
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
	}
}

type SubnetConverter struct {
}

func (v SubnetConverter) NewArg() proto.Message {
	return &resources.CloudSpecificSubnetArgs{}
}

func (v SubnetConverter) ConvertToMultyResource(resourceId string, m proto.Message, otherResources map[string]common_resources.CloudSpecificResource) common_resources.CloudSpecificResource {
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

	// Connect to vn in the same cloud
	subnet.VirtualNetwork = otherResources[common_resources.GetCloudSpecificResourceId(&subnet, c)].Resource.(*types.VirtualNetwork)

	return common_resources.CloudSpecificResource{
		Cloud:             c,
		Resource:          &subnet,
		ImplicitlyCreated: false,
	}
}
