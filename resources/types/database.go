package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/database"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

/*
aws_db_instance
aws_db_subnet_group

azurerm_*_server
*/

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

func (db *Database) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	var subnetList []*Subnet
	for _, id := range db.SubnetIds {
		if s, err := ctx.GetResource(id); err != nil {
			db.LogFatal(db.ResourceId, "subnet_ids", err.Error())
		} else {
			subnet := s.Resource.(*Subnet)
			subnetList = append(subnetList, subnet)
		}
	}

	if cloud == common.AWS {
		var awsSubnetIds []string
		for _, sub := range subnetList {
			awsSubnetIds = append(awsSubnetIds, sub.GetId(cloud))
		}
		name := common.RemoveSpecialChars(db.Name)
		dbSubnetGroup := database.AwsDbSubnetGroup{
			AwsResource: common.AwsResource{
				ResourceName: "aws_db_subnet_group",
				ResourceId:   db.GetTfResourceId(cloud),
				Name:         db.Name,
			},
			SubnetIds: awsSubnetIds,
		}
		return []interface{}{
			dbSubnetGroup,
			database.AwsDbInstance{
				AwsResource: common.AwsResource{
					ResourceName: "aws_db_instance",
					ResourceId:   db.GetTfResourceId(cloud),
					Name:         name,
				},
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
		var azSubnets []string
		for _, sub := range subnetList {
			azSubnets = append(azSubnets, sub.GetId(cloud))
		}

		return database.NewAzureDatabase(database.AzureDbServer{
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
			SubnetIds:                  azSubnets,
		})
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
