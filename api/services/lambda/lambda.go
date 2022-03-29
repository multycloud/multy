package lambda

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type LambdaService struct {
	Service services.Service[*resources.LambdaArgs, *resources.LambdaResource]
}

func (s LambdaService) Convert(resourceId string, args *resources.LambdaArgs, state *output.TfState) (*resources.LambdaResource, error) {
	return &resources.LambdaResource{
		CommonParameters:   util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:               args.Name,
		Runtime:            args.Runtime,
		SourceCodeObjectId: args.SourceCodeObjectId,
	}, nil
}

func NewLambdaService(database *db.Database) LambdaService {
	ni := LambdaService{
		Service: services.Service[*resources.LambdaArgs, *resources.LambdaResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
