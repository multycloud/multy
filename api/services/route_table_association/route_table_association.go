package route_table_association

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type RouteTableAssociationService struct {
	Service services.Service[*resourcespb.RouteTableAssociationArgs, *resourcespb.RouteTableAssociationResource]
}

func (s RouteTableAssociationService) Convert(resourceId string, args *resourcespb.RouteTableAssociationArgs, state *output.TfState) (*resourcespb.RouteTableAssociationResource, error) {
	return &resourcespb.RouteTableAssociationResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		SubnetId:         args.SubnetId,
		RouteTableId:     args.RouteTableId,
	}, nil
}

func NewRouteTableAssociationService(database *db.Database) RouteTableAssociationService {
	rta := RouteTableAssociationService{
		Service: services.Service[*resourcespb.RouteTableAssociationArgs, *resourcespb.RouteTableAssociationResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "route_table_associatio",
		},
	}
	rta.Service.Converters = &rta
	return rta
}
