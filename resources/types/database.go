package types

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/database"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

type Database struct {
	*resources.CommonResourceParams
	Name          string   `hcl:"name"`
	Engine        string   `hcl:"engine"`
	EngineVersion string   `hcl:"engine_version"`
	Storage       int      `hcl:"storage"`
	Size          string   `hcl:"size"`
	DbUsername    string   `hcl:"db_username"`
	DbPassword    string   `hcl:"db_password"`
	SubnetIds     []string `hcl:"subnet_ids"`
}

func (db *Database) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		name := common.RemoveSpecialChars(db.Name)
		// TODO validate subnet configuration (minimum 2 different AZs)
		dbSubnetGroup := database.AwsDbSubnetGroup{
			AwsResource: common.NewAwsResource(db.GetTfResourceId(cloud), db.Name),
			Name:        db.Name,
			SubnetIds:   db.SubnetIds,
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
		}
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
				SubnetIds:                  db.SubnetIds,
			},
		)
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (db *Database) Validate(ctx resources.MultyContext) {
	if db.Engine != "mysql" {
		db.LogFatal(db.ResourceId, "engine", fmt.Sprintf("\"%s\" is not valid a valid Engine", db.Engine))
	}
	if db.Storage < 10 && db.Storage < 20 {
		db.LogFatal(db.ResourceId, "storage", "storage must be between 10 and 20")
	}
	// TODO regex validate db username && password
	// TODO validate DB Size
	return
}

func (db *Database) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return database.AwsResourceName
	case common.AZURE:
		if db.Engine == "mysql" {
			return database.AzureMysqlResourceName
		}
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}

func (db *Database) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"password": cty.StringVal(db.DbPassword),
			"host": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.address}", common.GetResourceName(database.AwsDbInstance{}),
					db.GetTfResourceId(cloud),
				),
			),
			"username": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.username}", common.GetResourceName(database.AwsDbInstance{}),
					db.GetTfResourceId(cloud),
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"password": cty.StringVal(db.DbPassword),
			"host": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.fqdn}", common.GetResourceName(database.AzureMySqlServer{}),
					db.GetTfResourceId(cloud),
				),
			),
			"username": cty.StringVal(
				fmt.Sprintf("%s@%s", db.DbUsername, db.Name),
			),
		}
	}

	validate.LogInternalError("unknown cloud %s", cloud)
	return nil
}
