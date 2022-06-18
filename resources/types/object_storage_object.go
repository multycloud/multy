package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

// AWS: aws_s3_object
// Azure: azurerm_storage_blob

type ObjectStorageObject struct {
	resources.ChildResourceWithId[*ObjectStorage, *resourcespb.ObjectStorageObjectArgs]

	ObjectStorage *ObjectStorage `mhcl:"ref=object_storage"`
}

func (r *ObjectStorageObject) Create(resourceId string, args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) error {
	return NewObjectStorageObject(r, resourceId, args, others)
}

func (r *ObjectStorageObject) Update(args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) error {
	return NewObjectStorageObject(r, r.ResourceId, args, others)
}

func (r *ObjectStorageObject) Import(resourceId string, args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) error {
	return NewObjectStorageObject(r, resourceId, args, others)
}

func (r *ObjectStorageObject) Export(_ *resources.Resources) (*resourcespb.ObjectStorageObjectArgs, bool, error) {
	return r.Args, true, nil
}

func NewObjectStorageObject(r *ObjectStorageObject, resourceId string, args *resourcespb.ObjectStorageObjectArgs, others *resources.Resources) error {
	obj, err := resources.Get[*ObjectStorage](resourceId, others, args.ObjectStorageId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, r.ResourceId, "object_storage_id"))
	}
	r.ChildResourceWithId = resources.NewChildResource(resourceId, obj, args)
	r.Parent = obj
	r.ObjectStorage = obj
	return nil
}

func (r *ObjectStorageObject) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if len(r.Args.ContentBase64) == 0 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("content_base64 must be set"), ""))
	}
	return errs
}
