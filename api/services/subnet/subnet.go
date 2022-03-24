package subnet

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type SubnetService struct {
	Service services.Service[*resources.CloudSpecificSubnetArgs, *resources.SubnetResource]
}

func (s SubnetService) Convert(resourceId string, args []*resources.CloudSpecificSubnetArgs, state *output.TfState) (*resources.SubnetResource, error) {
	var result []*resources.CloudSpecificSubnetResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificSubnetResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
			AvailabilityZone: r.AvailabilityZone,
			VirtualNetworkId: r.VirtualNetworkId,
		})
	}

	return &resources.SubnetResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func NewSubnetService(database *db.Database) SubnetService {
	vn := SubnetService{
		Service: services.Service[*resources.CloudSpecificSubnetArgs, *resources.SubnetResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.Service.Converters = &vn
	return vn
}
