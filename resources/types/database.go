package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
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
	SubnetIds     []string `hcl:"subnet_ids,optional"`
}

func (db *Database) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []any {
	if cloud == common.AWS {
		name := common.RemoveSpecialChars(db.Name)
		dbSubnetGroup := database.AwsDbSubnetGroup{
			AwsResource: common.NewAwsResource("aws_db_subnet_group", db.GetTfResourceId(cloud), db.Name),
			SubnetIds:   db.SubnetIds,
		}
		return []any{
			dbSubnetGroup,
			database.AwsDbInstance{
				AwsResource:        common.NewAwsResource("aws_db_instance", db.GetTfResourceId(cloud), name),
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
				AzResource: common.AzResource{
					ResourceId:        db.GetTfResourceId(cloud),
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
