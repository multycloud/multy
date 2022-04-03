package lambda

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type LambdaService struct {
	Service services.Service[*resourcespb.LambdaArgs, *resourcespb.LambdaResource]
}

func (s LambdaService) Convert(resourceId string, args *resourcespb.LambdaArgs, state *output.TfState) (*resourcespb.LambdaResource, error) {
	return &resourcespb.LambdaResource{
		CommonParameters:   util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:               args.Name,
		Runtime:            args.Runtime,
		SourceCodeObjectId: args.SourceCodeObjectId,
	}, nil
}

func NewLambdaService(database *db.Database) LambdaService {
	ni := LambdaService{
		Service: services.Service[*resourcespb.LambdaArgs, *resourcespb.LambdaResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
