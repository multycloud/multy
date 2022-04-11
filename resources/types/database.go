package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"strings"
)

type Database struct {
	resources.ResourceWithId[*resourcespb.DatabaseArgs]

	Subnets []*Subnet
}

func NewDatabase(resourceId string, db *resourcespb.DatabaseArgs, others resources.Resources) (*Database, error) {
	subnets, err := util.MapSliceValuesErr(db.SubnetIds, func(subnetId string) (*Subnet, error) {
		return resources.Get[*Subnet](resourceId, others, subnetId)
	})
	if err != nil {
		return nil, err
	}
	return &Database{
		ResourceWithId: resources.ResourceWithId[*resourcespb.DatabaseArgs]{
			ResourceId: resourceId,
			Args:       db,
		},
		Subnets: subnets,
	}, nil
}

func (r *Database) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds, err := util.MapSliceValuesErr(r.Subnets, func(v *Subnet) (string, error) {
		return resources.GetMainOutputId(v)
	})
	if err != nil {
		return nil, err
	}
	// TODO validate subnet configuration (minimum 2 different AZs)
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		name := common.RemoveSpecialChars(r.Args.Name)
		dbSubnetGroup := database.AwsDbSubnetGroup{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			Name:        r.Args.Name,
			Description: "Managed by Multy",
			SubnetIds:   subnetIds,
		}
		return []output.TfBlock{
			dbSubnetGroup,
			database.AwsDbInstance{
				AwsResource:        common.NewAwsResource(r.ResourceId, name),
				Name:               name,
				AllocatedStorage:   int(r.Args.StorageGb),
				Engine:             strings.ToLower(r.Args.Engine.String()),
				EngineVersion:      r.Args.EngineVersion,
				Username:           r.Args.Username,
				Password:           r.Args.Password,
				InstanceClass:      common.DBSIZE[r.Args.Size][r.GetCloud()],
				Identifier:         r.Args.Name,
				SkipFinalSnapshot:  true,
				DbSubnetGroupName:  dbSubnetGroup.GetResourceName(),
				PubliclyAccessible: true,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return database.NewAzureDatabase(
			database.AzureDbServer{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					Name:              r.Args.Name,
					ResourceGroupName: rg.GetResourceGroupName(r.Args.GetCommonParameters().ResourceGroupId),
					Location:          r.GetCloudSpecificLocation(),
				},
				Engine:                     strings.ToLower(r.Args.Engine.String()),
				Version:                    r.Args.EngineVersion,
				StorageMb:                  int(r.Args.StorageGb * 1024),
				AdministratorLogin:         r.Args.Username,
				AdministratorLoginPassword: r.Args.Password,
				SkuName:                    common.DBSIZE[r.Args.Size][r.GetCloud()],
				SubnetIds:                  subnetIds,
			},
		), nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *Database) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if r.Args.Engine == resourcespb.DatabaseEngine_UNKNOWN_ENGINE {
		errs = append(errs, r.NewValidationError("unknown database engine provided", "engine"))
	}
	if r.Args.StorageGb < 10 || r.Args.StorageGb > 20 {
		errs = append(errs, r.NewValidationError("storage must be between 10 and 20", "storage"))
	}
	// TODO regex validate r username && password
	// TODO validate DB Size
	return errs
}

func (r *Database) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return database.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		if r.Args.Engine == resourcespb.DatabaseEngine_MYSQL {
			return database.AzureMysqlResourceName, nil
		}
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
	return "", nil
}
