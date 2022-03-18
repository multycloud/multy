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
	ConvertToMultyResource(resourceId string, arg proto.Message) common_resources.CloudSpecificResource
	NewArg() proto.Message
}

type VnConverter struct {
}

func (v VnConverter) NewArg() proto.Message {
	return &resources.CloudSpecificVirtualNetworkArgs{}
}

func (v VnConverter) ConvertToMultyResource(resourceId string, m proto.Message) common_resources.CloudSpecificResource {
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
