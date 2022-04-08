package database

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	output_database "github.com/multycloud/multy/resources/output/database"
)

type DatabaseService struct {
	Service services.Service[*resourcespb.DatabaseArgs, *resourcespb.DatabaseResource]
}

func (s DatabaseService) Convert(resourceId string, args *resourcespb.DatabaseArgs, state *output.TfState) (*resourcespb.DatabaseResource, error) {
	host, err := getHost(resourceId, state, args.CommonParameters.CloudProvider)
	if err != nil {
		return nil, err
	}
	return &resourcespb.DatabaseResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		Engine:           args.Engine,
		EngineVersion:    args.EngineVersion,
		StorageGb:        args.StorageGb,
		Size:             args.Size,
		Username:         args.Username,
		Password:         args.Password,
		SubnetIds:        args.SubnetIds,
		Host:             host,
	}, nil
}

func getHost(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(output_database.AwsDbInstance{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["address"].(string), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(output_database.AzureMySqlServer{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["fqdn"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func NewDatabaseService(database *db.Database) DatabaseService {
	ni := DatabaseService{
		Service: services.Service[*resourcespb.DatabaseArgs, *resourcespb.DatabaseResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
