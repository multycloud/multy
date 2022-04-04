package network_security_group

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type NetworkSecurityGroupService struct {
	Service services.Service[*resourcespb.NetworkSecurityGroupArgs, *resourcespb.NetworkSecurityGroupResource]
}

func (s NetworkSecurityGroupService) Convert(resourceId string, args *resourcespb.NetworkSecurityGroupArgs, state *output.TfState) (*resourcespb.NetworkSecurityGroupResource, error) {
	return &resourcespb.NetworkSecurityGroupResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		VirtualNetworkId: args.VirtualNetworkId,
		Rules:            args.Rules,
	}, nil
}

func NewNetworkSecurityGroupService(database *db.Database) NetworkSecurityGroupService {
	nsg := NetworkSecurityGroupService{
		Service: services.Service[*resourcespb.NetworkSecurityGroupArgs, *resourcespb.NetworkSecurityGroupResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
