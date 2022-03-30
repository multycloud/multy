package route_table

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type RouteTableService struct {
	Service services.Service[*resources.RouteTableArgs, *resources.RouteTableResource]
}

func (s RouteTableService) Convert(resourceId string, args *resources.RouteTableArgs, state *output.TfState) (*resources.RouteTableResource, error) {
	return &resources.RouteTableResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		VirtualNetworkId: args.VirtualNetworkId,
		Routes:           args.Routes,
	}, nil
}

func NewRouteTableService(database *db.Database) RouteTableService {
	rt := RouteTableService{
		Service: services.Service[*resources.RouteTableArgs, *resources.RouteTableResource]{
			Db:         database,
			Converters: nil,
		},
	}
	rt.Service.Converters = &rt
	return rt
}
