package route_table

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type RouteTableService struct {
	Service services.Service[*resourcespb.RouteTableArgs, *resourcespb.RouteTableResource]
}

func (s RouteTableService) Convert(resourceId string, args *resourcespb.RouteTableArgs, state *output.TfState) (*resourcespb.RouteTableResource, error) {
	return &resourcespb.RouteTableResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		VirtualNetworkId: args.VirtualNetworkId,
		Routes:           args.Routes,
	}, nil
}

func NewRouteTableService(database *db.Database) RouteTableService {
	rt := RouteTableService{
		Service: services.Service[*resourcespb.RouteTableArgs, *resourcespb.RouteTableResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "route_table",
		},
	}
	rt.Service.Converters = &rt
	return rt
}
