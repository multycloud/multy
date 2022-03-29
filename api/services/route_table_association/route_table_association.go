package route_table_association

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type RouteTableAssociationService struct {
	Service services.Service[*resources.RouteTableAssociationArgs, *resources.RouteTableAssociationResource]
}

func (s RouteTableAssociationService) Convert(resourceId string, args *resources.RouteTableAssociationArgs, state *output.TfState) (*resources.RouteTableAssociationResource, error) {
	return &resources.RouteTableAssociationResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		SubnetId:         args.SubnetId,
		RouteTableId:     args.RouteTableId,
	}, nil
}

func NewRouteTableAssociationService(database *db.Database) RouteTableAssociationService {
	rta := RouteTableAssociationService{
		Service: services.Service[*resources.RouteTableAssociationArgs, *resources.RouteTableAssociationResource]{
			Db:         database,
			Converters: nil,
		},
	}
	rta.Service.Converters = &rta
	return rta
}
