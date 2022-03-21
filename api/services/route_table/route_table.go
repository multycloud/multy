package route_table

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type RouteTableService struct {
	Service services.Service[*resources.CloudSpecificRouteTableArgs, *resources.RouteTableResource]
}

func (s RouteTableService) Convert(resourceId string, args []*resources.CloudSpecificRouteTableArgs) *resources.RouteTableResource {
	var result []*resources.CloudSpecificRouteTableResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificRouteTableResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			VirtualNetworkId: r.VirtualNetworkId,
			Routes:           r.Routes,
		})
	}

	return &resources.RouteTableResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s RouteTableService) NewArg() *resources.CloudSpecificRouteTableArgs {
	return &resources.CloudSpecificRouteTableArgs{}
}

func (s RouteTableService) Nil() *resources.RouteTableResource {
	return nil
}

func NewRouteTableServiceService(database *db.Database) RouteTableService {
	rt := RouteTableService{
		Service: services.Service[*resources.CloudSpecificRouteTableArgs, *resources.RouteTableResource]{
			Db:         database,
			Converters: nil,
		},
	}
	rt.Service.Converters = &rt
	return rt
}
