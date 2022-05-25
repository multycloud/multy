package types

import (
	"fmt"
	"strings"

	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

var dbMetadata = resources.ResourceMetadata[*resourcespb.DatabaseArgs, *Database, *resourcespb.DatabaseResource]{
	CreateFunc:        CreateDatabase,
	UpdateFunc:        UpdateDatabase,
	ReadFromStateFunc: DatabaseFromState,
	ExportFunc: func(r *Database, _ *resources.Resources) (*resourcespb.DatabaseArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewDatabase,
	AbbreviatedName: "db",
}

type Database struct {
	resources.ResourceWithId[*resourcespb.DatabaseArgs]

	Subnets []*Subnet
}

func (r *Database) GetMetadata() resources.ResourceMetadataInterface {
	return &dbMetadata
}

func NewDatabase(resourceId string, db *resourcespb.DatabaseArgs, others *resources.Resources) (*Database, error) {
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

func CreateDatabase(resourceId string, args *resourcespb.DatabaseArgs, others *resources.Resources) (*Database, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetIds[0])
		if err != nil {
			return nil, err
		}
		rgId := NewRgFromParent("db", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewDatabase(resourceId, args, others)
}

func UpdateDatabase(resource *Database, vn *resourcespb.DatabaseArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func DatabaseFromState(resource *Database, state *output.TfState) (*resourcespb.DatabaseResource, error) {
	var err error
	host := "dyrun"
	if !flags.DryRun {
		host, err = getHost(resource.ResourceId, state, resource.Args.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
	}
	connectionUsername := resource.Args.Username
	if resource.GetCloud() == commonpb.CloudProvider_AZURE {
		connectionUsername = fmt.Sprintf("%s@%s", resource.Args.Username, host)
	}

	return &resourcespb.DatabaseResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:               resource.Args.Name,
		Engine:             resource.Args.Engine,
		EngineVersion:      resource.Args.EngineVersion,
		StorageGb:          resource.Args.StorageGb,
		Size:               resource.Args.Size,
		Username:           resource.Args.Username,
		Password:           resource.Args.Password,
		SubnetIds:          resource.Args.SubnetIds,
		Host:               host,
		ConnectionUsername: connectionUsername,
	}, nil
}

func getHost(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(database.AwsDbInstance{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["address"].(string), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(database.AzureMySqlServer{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["fqdn"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
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
				Port:               int(r.Args.Port),
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
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
				SubnetIds:                  subnetIds,
				Port:                       int(r.Args.Port),
			},
		), nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *Database) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if r.Args.Engine == resourcespb.DatabaseEngine_UNKNOWN_ENGINE {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown database engine provided"), "engine"))
	}
	if r.Args.StorageGb < 10 || r.Args.StorageGb > 20 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("storage must be between 10 and 20"), "storage"))
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
