package route_table_association

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type RouteTableAssociationService struct {
	Service services.Service[*resources.CloudSpecificRouteTableAssociationArgs, *resources.RouteTableAssociationResource]
}

func (s RouteTableAssociationService) Convert(resourceId string, args []*resources.CloudSpecificRouteTableAssociationArgs) *resources.RouteTableAssociationResource {
	var result []*resources.CloudSpecificRouteTableAssociationResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificRouteTableAssociationResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			SubnetId:         r.SubnetId,
			RouteTableId:     r.RouteTableId,
		})
	}

	return &resources.RouteTableAssociationResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s RouteTableAssociationService) NewArg() *resources.CloudSpecificRouteTableAssociationArgs {
	return &resources.CloudSpecificRouteTableAssociationArgs{}
}

func (s RouteTableAssociationService) Nil() *resources.RouteTableAssociationResource {
	return nil
}

func NewRouteTableAssociationServiceService(database *db.Database) RouteTableAssociationService {
	rta := RouteTableAssociationService{
		Service: services.Service[*resources.CloudSpecificRouteTableAssociationArgs, *resources.RouteTableAssociationResource]{
			Db:         database,
			Converters: nil,
		},
	}
	rta.Service.Converters = &rta
	return rta
}
