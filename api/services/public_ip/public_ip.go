package public_ip

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type PublicIpService struct {
	Service services.Service[*resources.CloudSpecificPublicIpArgs, *resources.PublicIpResource]
}

func (s PublicIpService) Convert(resourceId string, args []*resources.CloudSpecificPublicIpArgs, state *output.TfState) (*resources.PublicIpResource, error) {
	var result []*resources.CloudSpecificPublicIpResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificPublicIpResource{
			CommonParameters:   util.ConvertCommonParams(r.CommonParameters),
			Name:               r.Name,
			NetworkInterfaceId: r.NetworkInterfaceId,
		})
	}

	return &resources.PublicIpResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func (s PublicIpService) NewArg() *resources.CloudSpecificPublicIpArgs {
	return &resources.CloudSpecificPublicIpArgs{}
}

func (s PublicIpService) Nil() *resources.PublicIpResource {
	return nil
}

func NewPublicIpService(database *db.Database) PublicIpService {
	nsg := PublicIpService{
		Service: services.Service[*resources.CloudSpecificPublicIpArgs, *resources.PublicIpResource]{
			Db:         database,
			Converters: nil,
		},
	}
	nsg.Service.Converters = &nsg
	return nsg
}
