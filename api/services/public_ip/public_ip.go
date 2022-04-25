package public_ip

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type PublicIpService struct {
	Service services.Service[*resourcespb.PublicIpArgs, *resourcespb.PublicIpResource]
}

func (s PublicIpService) Convert(resourceId string, args *resourcespb.PublicIpArgs, state *output.TfState) (*resourcespb.PublicIpResource, error) {
	return &resourcespb.PublicIpResource{
		CommonParameters:   util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:               args.Name,
		NetworkInterfaceId: args.NetworkInterfaceId,
	}, nil
}

func NewPublicIpService(database *db.Database) PublicIpService {
	nsg := PublicIpService{
		Service: services.Service[*resourcespb.PublicIpArgs, *resourcespb.PublicIpResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "public_ip",
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
