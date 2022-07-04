package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	"github.com/multycloud/multy/resources/types"
	"strings"
)

type AzureDatabase struct {
	*types.Database
}

func InitDatabase(r *types.Database) resources.ResourceTranslator[*resourcespb.DatabaseResource] {
	return AzureDatabase{r}
}

func (r AzureDatabase) FromState(state *output.TfState) (*resourcespb.DatabaseResource, error) {
	host := "dyrun"
	if !flags.DryRun {
		var azureDatabaseEngine database.AzureDatabaseEngine

		if r.Args.Engine == resourcespb.DatabaseEngine_MYSQL {
			azureDatabaseEngine = database.AzureMySqlServer{}
		} else if r.Args.Engine == resourcespb.DatabaseEngine_POSTGRES {
			azureDatabaseEngine = database.AzurePostgreSqlServer{}
		} else if r.Args.Engine == resourcespb.DatabaseEngine_MARIADB {
			azureDatabaseEngine = database.AzureMariaDbServer{}
		} else {
			return nil, fmt.Errorf("unhandled engine %s", r.Args.Engine.String())
		}

		values, err := state.GetValues(azureDatabaseEngine, r.ResourceId)
		if err != nil {
			return nil, err
		}
		host = values["fqdn"].(string)
	}

	return &resourcespb.DatabaseResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:               r.Args.Name,
		Engine:             r.Args.Engine,
		EngineVersion:      r.Args.EngineVersion,
		StorageGb:          r.Args.StorageGb,
		Size:               r.Args.Size,
		Username:           r.Args.Username,
		Password:           r.Args.Password,
		SubnetIds:          r.Args.SubnetIds,
		Port:               r.Args.Port,
		SubnetId:           r.Args.SubnetId,
		Host:               host,
		ConnectionUsername: fmt.Sprintf("%s@%s", r.Args.Username, host),
	}, nil
}

func (r AzureDatabase) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	subnetId, err := resources.GetMainOutputId(AzureSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}
	return database.NewAzureDatabase(
		database.AzureDbServer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				Name:              r.Args.Name,
				ResourceGroupName: GetResourceGroupName(r.Args.GetCommonParameters().ResourceGroupId),
				Location:          r.GetCloudSpecificLocation(),
			},
			Engine:                     strings.ToLower(r.Args.Engine.String()),
			Version:                    r.Args.EngineVersion,
			StorageMb:                  int(r.Args.StorageGb * 1024),
			AdministratorLogin:         r.Args.Username,
			AdministratorLoginPassword: r.Args.Password,
			SkuName:                    common.DBSIZE[r.Args.Size][r.GetCloud()],
			SubnetId:                   subnetId,
		},
	), nil
}

func (r AzureDatabase) GetMainResourceName() (string, error) {
	switch r.Args.Engine {
	case resourcespb.DatabaseEngine_MYSQL:
		return database.AzureMysqlResourceName, nil
	case resourcespb.DatabaseEngine_MARIADB:
		return database.AzureMariaDbResourceName, nil
	case resourcespb.DatabaseEngine_POSTGRES:
		return database.AzurePostgresqlResourceName, nil
	}
	return "", fmt.Errorf("unhandled engine %s", r.Args.Engine.String())
}
