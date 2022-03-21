package network_security_group

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type NetworkSecurityGroupService struct {
	Service services.Service[*resources.CloudSpecificNetworkSecurityGroupArgs, *resources.NetworkSecurityGroupResource]
}

func (s NetworkSecurityGroupService) Convert(resourceId string, args []*resources.CloudSpecificNetworkSecurityGroupArgs) *resources.NetworkSecurityGroupResource {
	var result []*resources.CloudSpecificNetworkSecurityGroupResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificNetworkSecurityGroupResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			VirtualNetworkId: r.VirtualNetworkId,
			Rules:            r.Rules,
		})
	}

	return &resources.NetworkSecurityGroupResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s NetworkSecurityGroupService) NewArg() *resources.CloudSpecificNetworkSecurityGroupArgs {
	return &resources.CloudSpecificNetworkSecurityGroupArgs{}
}

func (s NetworkSecurityGroupService) Nil() *resources.NetworkSecurityGroupResource {
	return nil
}

func NewNetworkSecurityGroupServiceService(database *db.Database) NetworkSecurityGroupService {
	nsg := NetworkSecurityGroupService{
		Service: services.Service[*resources.CloudSpecificNetworkSecurityGroupArgs, *resources.NetworkSecurityGroupResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
