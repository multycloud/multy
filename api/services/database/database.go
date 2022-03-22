package database

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	output_database "github.com/multycloud/multy/resources/output/database"
	common_util "github.com/multycloud/multy/util"
)

type DatabaseService struct {
	Service services.Service[*resources.CloudSpecificDatabaseArgs, *resources.DatabaseResource]
}

func (s DatabaseService) Convert(resourceId string, args []*resources.CloudSpecificDatabaseArgs, state *output.TfState) (*resources.DatabaseResource, error) {
	var result []*resources.CloudSpecificDatabaseResource
	for _, r := range args {
		host, err := getHost(resourceId, state, r.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
		result = append(result, &resources.CloudSpecificDatabaseResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			Engine:           r.Engine,
			EngineVersion:    r.EngineVersion,
			StorageGb:        r.StorageGb,
			Size:             r.Size,
			Username:         r.Username,
			Password:         r.Password,
			SubnetIds:        r.SubnetIds,
			Host:             host,
		})
	}

	return &resources.DatabaseResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func getHost(resourceId string, state *output.TfState, cloud common.CloudProvider) (string, error) {
	rId := common_util.GetTfResourceId(resourceId, cloud.String())
	switch cloud {
	case common.CloudProvider_AWS:
		values, err := state.GetValues(output_database.AwsDbInstance{}, rId)
		if err != nil {
			return "", err
		}
		return values["address"].(string), nil
	case common.CloudProvider_AZURE:
		values, err := state.GetValues(output_database.AzureMySqlServer{}, rId)
		if err != nil {
			return "", err
		}
		return values["fqdn"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())

}

func (s DatabaseService) NewArg() *resources.CloudSpecificDatabaseArgs {
	return &resources.CloudSpecificDatabaseArgs{}
}

func (s DatabaseService) Nil() *resources.DatabaseResource {
	return nil
}

func NewDatabaseService(database *db.Database) DatabaseService {
	ni := DatabaseService{
		Service: services.Service[*resources.CloudSpecificDatabaseArgs, *resources.DatabaseResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
