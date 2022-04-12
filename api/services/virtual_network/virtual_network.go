package virtual_network

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VnService struct {
	Service services.Service[*resourcespb.VirtualNetworkArgs, *resourcespb.VirtualNetworkResource]
}

func (s VnService) Convert(resourceId string, args *resourcespb.VirtualNetworkArgs, state *output.TfState) (*resourcespb.VirtualNetworkResource, error) {
	return &resourcespb.VirtualNetworkResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		CidrBlock:        args.CidrBlock,
	}, nil
}

func NewVnService(database *db.Database) VnService {
	vn := VnService{
		Service: services.Service[*resourcespb.VirtualNetworkArgs, *resourcespb.VirtualNetworkResource]{
			Db: database,
		},
	}
	vn.Service.Converters = &vn
	return vn
}
