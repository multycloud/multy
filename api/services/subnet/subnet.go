package subnet

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type SubnetService struct {
	Service services.Service[*resources.SubnetArgs, *resources.SubnetResource]
}

func (s SubnetService) Convert(resourceId string, args *resources.SubnetArgs, state *output.TfState) (*resources.SubnetResource, error) {
	return &resources.SubnetResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		CidrBlock:        args.CidrBlock,
		AvailabilityZone: args.AvailabilityZone,
		VirtualNetworkId: args.VirtualNetworkId,
	}, nil
}

func NewSubnetService(database *db.Database) SubnetService {
	vn := SubnetService{
		Service: services.Service[*resources.SubnetArgs, *resources.SubnetResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.Service.Converters = &vn
	return vn
}
