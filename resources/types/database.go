package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

type Database struct {
	resources.ResourceWithId[*resourcespb.DatabaseArgs]

	Subnets []*Subnet
}

func (r *Database) Create(resourceId string, args *resourcespb.DatabaseArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetIds[0])
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

func (r *Database) Update(args *resourcespb.DatabaseArgs, _ *resources.Resources) error {
	r.Args = args
	return nil
}

func (r *Database) Import(resourceId string, args *resourcespb.DatabaseArgs, others *resources.Resources) error {
	return NewDatabase(r, resourceId, args, others)
}

func (r *Database) Export(_ *resources.Resources) (*resourcespb.DatabaseArgs, bool, error) {
	return r.Args, true, nil
}

func NewDatabase(r *Database, resourceId string, db *resourcespb.DatabaseArgs, others *resources.Resources) error {
	subnets, err := util.MapSliceValuesErr(db.SubnetIds, func(subnetId string) (*Subnet, error) {
		return resources.Get[*Subnet](resourceId, others, subnetId)
	})
	if err != nil {
		return err
	}
	r.Subnets = subnets
	r.ResourceWithId = resources.ResourceWithId[*resourcespb.DatabaseArgs]{
		ResourceId: resourceId,
		Args:       db,
	}

	return nil
}

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
	// TODO regex validate r username && password
	// TODO validate DB Size
	return errs
}
