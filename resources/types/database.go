package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
)

type Database struct {
	*resources.CommonResourceParams
	Name          string    `hcl:"name"`
	Engine        string    `hcl:"engine"`
	EngineVersion string    `hcl:"engine_version"`
	Storage       int       `hcl:"storage"`
	Size          string    `hcl:"size"`
	DbUsername    string    `hcl:"db_username"`
	DbPassword    string    `hcl:"db_password"`
	SubnetIds     []*Subnet `mhcl:"ref=subnet_ids"`
}

func (db *Database) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds, err := util.MapSliceValuesErr(db.SubnetIds, func(v *Subnet) (string, error) {
		return resources.GetMainOutputId(v, cloud)
	})
	if err != nil {
		return nil, err
	}
	// TODO validate subnet configuration (minimum 2 different AZs)
	if cloud == common.AWS {
		name := common.RemoveSpecialChars(db.Name)
		dbSubnetGroup := database.AwsDbSubnetGroup{
			AwsResource: common.NewAwsResource(db.GetTfResourceId(cloud), db.Name),
			Name:        db.Name,
			SubnetIds:   subnetIds,
		}
		return []output.TfBlock{
			dbSubnetGroup,
			database.AwsDbInstance{
				AwsResource:        common.NewAwsResource(db.GetTfResourceId(cloud), name),
				Name:               name,
				AllocatedStorage:   db.Storage,
				Engine:             db.Engine,
				EngineVersion:      db.EngineVersion,
				Username:           db.DbUsername,
				Password:           db.DbPassword,
				InstanceClass:      common.DBSIZE[db.Size][cloud],
				Identifier:         db.Name,
				SkipFinalSnapshot:  true,
				DbSubnetGroupName:  dbSubnetGroup.GetResourceName(),
				PubliclyAccessible: true,
			},
		}, nil
	} else if cloud == common.AZURE {
		return database.NewAzureDatabase(
			database.AzureDbServer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: db.GetTfResourceId(cloud)},
					Name:              db.Name,
					ResourceGroupName: rg.GetResourceGroupName(db.ResourceGroupId, cloud),
					Location:          ctx.GetLocationFromCommonParams(db.CommonResourceParams, cloud),
				},
				Engine:                     db.Engine,
				Version:                    db.EngineVersion,
				StorageMb:                  db.Storage * 1024,
				AdministratorLogin:         db.DbUsername,
				AdministratorLoginPassword: db.DbPassword,
				SkuName:                    common.DBSIZE[db.Size][cloud],
				SubnetIds:                  subnetIds,
			},
		), nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (db *Database) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, db.CommonResourceParams.Validate(ctx, cloud)...)
	if db.Engine != "mysql" {
		errs = append(errs, db.NewError("engine", fmt.Sprintf("\"%s\" is not valid a valid Engine", db.Engine)))
	}
	if db.Storage < 10 || db.Storage > 20 {
		errs = append(errs, db.NewError("storage", "storage must be between 10 and 20"))
	}
	// TODO regex validate db username && password
	// TODO validate DB Size
	return errs
}

func (db *Database) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return database.AwsResourceName, nil
	case common.AZURE:
		if db.Engine == "mysql" {
			return database.AzureMysqlResourceName, nil
		}
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
	return "", nil
}

func (db *Database) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"password": cty.StringVal(db.DbPassword),
			"host": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.address}", output.GetResourceName(database.AwsDbInstance{}),
					db.GetTfResourceId(cloud),
				),
			),
			"username": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.username}", output.GetResourceName(database.AwsDbInstance{}),
					db.GetTfResourceId(cloud),
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"password": cty.StringVal(db.DbPassword),
			"host": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.fqdn}", output.GetResourceName(database.AzureMySqlServer{}),
					db.GetTfResourceId(cloud),
				),
			),
			"username": cty.StringVal(
				fmt.Sprintf("%s@%s", db.DbUsername, db.Name),
			),
		}
	}
	return nil
}
