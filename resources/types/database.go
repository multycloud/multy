package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
	"golang.org/x/exp/slices"
)

type Database struct {
	resources.ResourceWithId[*resourcespb.DatabaseArgs]

	Subnet *Subnet
}

func (r *Database) Create(resourceId string, args *resourcespb.DatabaseArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return err
		}
		rgId, err := NewRgFromParent("db", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewDatabase(r, resourceId, args, others)
}

func (r *Database) Update(args *resourcespb.DatabaseArgs, others *resources.Resources) error {
	return NewDatabase(r, r.ResourceId, args, others)
}

func (r *Database) Import(resourceId string, args *resourcespb.DatabaseArgs, others *resources.Resources) error {
	return NewDatabase(r, resourceId, args, others)
}

func (r *Database) Export(_ *resources.Resources) (*resourcespb.DatabaseArgs, bool, error) {
	return r.Args, true, nil
}

func NewDatabase(r *Database, resourceId string, db *resourcespb.DatabaseArgs, others *resources.Resources) error {
	subnet, err := resources.Get[*Subnet](resourceId, others, db.SubnetId)
	if err != nil {
		return err
	}
	r.Subnet = subnet

	r.ResourceWithId = resources.ResourceWithId[*resourcespb.DatabaseArgs]{
		ResourceId: resourceId,
		Args:       db,
	}

	return nil
}

var MysqlVersions = []string{"5.6", "5.7", "8.0"}
var PostgresVersions = []string{"10", "11", "12", "13", "14"}
var MariaDBVersions = []string{"10.2", "10.3"} // https://docs.microsoft.com/bs-latn-ba/azure/mariadb/concepts-supported-versions

func (r *Database) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if r.Args.Engine == resourcespb.DatabaseEngine_UNKNOWN_ENGINE {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown database engine provided"), "engine"))
	}
	if r.Args.StorageGb < 10 || r.Args.StorageGb > 20 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("storage must be between 10 and 20"), "storage"))
	}
	if r.GetCloud() == commonpb.CloudProvider_AZURE && r.Args.Port != 0 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("azure doesn't support custom ports"), "port"))
	}
	if r.Args.Engine == resourcespb.DatabaseEngine_MYSQL {
		if !slices.Contains(MysqlVersions, r.Args.EngineVersion) {
			errs = append(errs, r.NewValidationError(fmt.Errorf("'%s' is an unsupported engine version for mysql, must be one of %+q", r.Args.EngineVersion, MysqlVersions), "engine_version"))
		}
	}
	if r.Args.Engine == resourcespb.DatabaseEngine_POSTGRES {
		if !slices.Contains(PostgresVersions, r.Args.EngineVersion) {
			errs = append(errs, r.NewValidationError(fmt.Errorf("'%s' is an unsupported engine version for postgres, must be one of %+q", r.Args.EngineVersion, PostgresVersions), "engine_version"))
		}
	}
	if r.Args.Engine == resourcespb.DatabaseEngine_MARIADB {
		if r.GetCloud() == commonpb.CloudProvider_GCP {
			errs = append(errs, r.NewValidationError(fmt.Errorf("mariadb is not supported in gcp"), "engine"))
		}
		if !slices.Contains(MariaDBVersions, r.Args.EngineVersion) {
			errs = append(errs, r.NewValidationError(fmt.Errorf("'%s' is an unsupported engine version for mariadb, must be one of %+q", r.Args.EngineVersion, MariaDBVersions), "engine_version"))
		}
	}
	if r.Args.Size == commonpb.DatabaseSize_UNKNOWN_VM_SIZE {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown database size"), "size"))
	}
	usernameValidator := validate.NewDbUsernameValidator(r.Args.Engine, r.Args.EngineVersion)
	passwordValidator := validate.NewDbPasswordValidator(r.Args.Engine)
	if err := usernameValidator.Check(r.Args.Username, r.ResourceId); err != nil {
		errs = append(errs, r.NewValidationError(err, "username"))
	}
	if err := passwordValidator.Check(r.Args.Password, r.ResourceId); err != nil {
		errs = append(errs, r.NewValidationError(err, "password"))
	}
	return errs
}
