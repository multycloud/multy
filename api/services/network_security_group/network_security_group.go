package network_security_group

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type NetworkSecurityGroupService struct {
	Service services.Service[*resources.NetworkSecurityGroupArgs, *resources.NetworkSecurityGroupResource]
}

func (s NetworkSecurityGroupService) Convert(resourceId string, args *resources.NetworkSecurityGroupArgs, state *output.TfState) (*resources.NetworkSecurityGroupResource, error) {
	return &resources.NetworkSecurityGroupResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		VirtualNetworkId: args.VirtualNetworkId,
		Rules:            args.Rules,
	}, nil
}

func NewNetworkSecurityGroupService(database *db.Database) NetworkSecurityGroupService {
	nsg := NetworkSecurityGroupService{
		Service: services.Service[*resources.NetworkSecurityGroupArgs, *resources.NetworkSecurityGroupResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
