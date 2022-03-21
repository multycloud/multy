package lambda

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type LambdaService struct {
	Service services.Service[*resources.CloudSpecificLambdaArgs, *resources.LambdaResource]
}

func (s LambdaService) Convert(resourceId string, args []*resources.CloudSpecificLambdaArgs) *resources.LambdaResource {
	var result []*resources.CloudSpecificLambdaResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificLambdaResource{
			CommonParameters:   util.ConvertCommonParams(r.CommonParameters),
			Name:               r.Name,
			Runtime:            r.Runtime,
			SourceCodeObjectId: r.SourceCodeObjectId,
		})
	}

	return &resources.LambdaResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s LambdaService) NewArg() *resources.CloudSpecificLambdaArgs {
	return &resources.CloudSpecificLambdaArgs{}
}

func (s LambdaService) Nil() *resources.LambdaResource {
	return nil
}

func NewLambdaService(database *db.Database) LambdaService {
	ni := LambdaService{
		Service: services.Service[*resources.CloudSpecificLambdaArgs, *resources.LambdaResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
