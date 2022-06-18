package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type ObjectStorage struct {
	resources.ResourceWithId[*resourcespb.ObjectStorageArgs]
}

func (r *ObjectStorage) Create(resourceId string, args *resourcespb.ObjectStorageArgs, others *resources.Resources) error {
	return CreateObjectStorage(r, resourceId, args, others)
}

func (r *ObjectStorage) Update(args *resourcespb.ObjectStorageArgs, others *resources.Resources) error {
	return NewObjectStorage(r, r.ResourceId, args, others)
}

func (r *ObjectStorage) Import(resourceId string, args *resourcespb.ObjectStorageArgs, others *resources.Resources) error {
	return NewObjectStorage(r, resourceId, args, others)
}

func (r *ObjectStorage) Export(_ *resources.Resources) (*resourcespb.ObjectStorageArgs, bool, error) {
	return r.Args, true, nil
}

func NewObjectStorage(r *ObjectStorage, resourceId string, db *resourcespb.ObjectStorageArgs, _ *resources.Resources) error {
	r.ResourceWithId = resources.ResourceWithId[*resourcespb.ObjectStorageArgs]{
		ResourceId: resourceId,
		Args:       db,
	}
	return nil
}

func CreateObjectStorage(r *ObjectStorage, resourceId string, args *resourcespb.ObjectStorageArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		rgId, err := NewRg("st", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewObjectStorage(r, resourceId, args, others)
}

func (r *ObjectStorage) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}
