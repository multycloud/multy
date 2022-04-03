package network_interface

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type NetworkInterfaceService struct {
	Service services.Service[*resourcespb.NetworkInterfaceArgs, *resourcespb.NetworkInterfaceResource]
}

func (s NetworkInterfaceService) Convert(resourceId string, args *resourcespb.NetworkInterfaceArgs, state *output.TfState) (*resourcespb.NetworkInterfaceResource, error) {
	return &resourcespb.NetworkInterfaceResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		SubnetId:         args.SubnetId,
	}, nil
}

func NewNetworkInterfaceService(database *db.Database) NetworkInterfaceService {
	ni := NetworkInterfaceService{
		Service: services.Service[*resourcespb.NetworkInterfaceArgs, *resourcespb.NetworkInterfaceResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
