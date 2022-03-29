package network_interface

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type NetworkInterfaceService struct {
	Service services.Service[*resources.NetworkInterfaceArgs, *resources.NetworkInterfaceResource]
}

func (s NetworkInterfaceService) Convert(resourceId string, args *resources.NetworkInterfaceArgs, state *output.TfState) (*resources.NetworkInterfaceResource, error) {
	return &resources.NetworkInterfaceResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		SubnetId:         args.SubnetId,
	}, nil
}

func NewNetworkInterfaceService(database *db.Database) NetworkInterfaceService {
	ni := NetworkInterfaceService{
		Service: services.Service[*resources.NetworkInterfaceArgs, *resources.NetworkInterfaceResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
