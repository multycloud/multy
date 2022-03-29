package public_ip

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type PublicIpService struct {
	Service services.Service[*resources.PublicIpArgs, *resources.PublicIpResource]
}

func (s PublicIpService) Convert(resourceId string, args *resources.PublicIpArgs, state *output.TfState) (*resources.PublicIpResource, error) {
	return &resources.PublicIpResource{
		CommonParameters:   util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:               args.Name,
		NetworkInterfaceId: args.NetworkInterfaceId,
	}, nil
}

func NewPublicIpService(database *db.Database) PublicIpService {
	nsg := PublicIpService{
		Service: services.Service[*resources.PublicIpArgs, *resources.PublicIpResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
