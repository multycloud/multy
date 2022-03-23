package virtual_network

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VnService struct {
	Service services.Service[*resources.CloudSpecificVirtualNetworkArgs, *resources.VirtualNetworkResource]
}

func (s VnService) Convert(resourceId string, args []*resources.CloudSpecificVirtualNetworkArgs, state *output.TfState) (*resources.VirtualNetworkResource, error) {
	var result []*resources.CloudSpecificVirtualNetworkResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificVirtualNetworkResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
		})
	}

	return &resources.VirtualNetworkResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func (s VnService) NewArg() *resources.CloudSpecificVirtualNetworkArgs {
	return &resources.CloudSpecificVirtualNetworkArgs{}
}

func (s VnService) Nil() *resources.VirtualNetworkResource {
	return nil
}

func NewVnService(database *db.Database) VnService {
	vn := VnService{
		Service: services.Service[*resources.CloudSpecificVirtualNetworkArgs, *resources.VirtualNetworkResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.Service.Converters = &vn
	return vn
}
