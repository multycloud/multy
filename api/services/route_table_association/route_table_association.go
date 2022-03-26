package route_table_association

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type RouteTableAssociationService struct {
	Service services.Service[*resources.CloudSpecificRouteTableAssociationArgs, *resources.RouteTableAssociationResource]
}

func (s RouteTableAssociationService) Convert(resourceId string, args []*resources.CloudSpecificRouteTableAssociationArgs, state *output.TfState) (*resources.RouteTableAssociationResource, error) {
	var result []*resources.CloudSpecificRouteTableAssociationResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificRouteTableAssociationResource{
			CommonParameters: util.ConvertCommonChildParams(r.CommonParameters),
			SubnetId:         r.SubnetId,
			RouteTableId:     r.RouteTableId,
		})
	}

	return &resources.RouteTableAssociationResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func NewRouteTableAssociationService(database *db.Database) RouteTableAssociationService {
	rta := RouteTableAssociationService{
		Service: services.Service[*resources.CloudSpecificRouteTableAssociationArgs, *resources.RouteTableAssociationResource]{
			Db:         database,
			Converters: nil,
		},
	}
	rta.Service.Converters = &rta
	return rta
}
