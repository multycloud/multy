package virtual_network

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VnService struct {
	Service services.Service[*resources.VirtualNetworkArgs, *resources.VirtualNetworkResource]
}

func (s VnService) Convert(resourceId string, args *resources.VirtualNetworkArgs, state *output.TfState) (*resources.VirtualNetworkResource, error) {
	return &resources.VirtualNetworkResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		CidrBlock:        args.CidrBlock,
	}, nil
}

func NewVnService(database *db.Database) VnService {
	vn := VnService{
		Service: services.Service[*resources.VirtualNetworkArgs, *resources.VirtualNetworkResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.Service.Converters = &vn
	return vn
}
