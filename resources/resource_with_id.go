package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
	"strings"
)

type WithCommonParams interface {
	GetCommonParameters() *commonpb.ResourceCommonArgs
}

type WithChildCommonParams interface {
	GetCommonParameters() *commonpb.ChildResourceCommonArgs
}

type ResourceWithId[T WithCommonParams] struct {
	ResourceId string
	Args       T
}

func (r *ResourceWithId[T]) GetCloudSpecificLocation() string {
	if result, err := common.GetCloudLocationPb(r.GetLocation(), r.GetCloud()); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}

func (r *ResourceWithId[T]) GetLocation() string {
	return strings.ToLower(r.Args.GetCommonParameters().GetLocation().String())
}

func (r *ResourceWithId[T]) GetResourceId() string {
	return r.ResourceId
}

func (r *ResourceWithId[T]) GetCloud() commonpb.CloudProvider {
	return r.Args.GetCommonParameters().CloudProvider
}

func (r *ResourceWithId[T]) Validate() (errs []validate.ValidationError) {
	location := r.GetLocation()
	if _, ok := common.LOCATION[location]; !ok {
		errs = append(errs, validate.ValidationError{
			ErrorMessage: fmt.Sprintf("location %s is not defined", location),
			ResourceId:   r.ResourceId,
			FieldName:    "location",
		})
	}

	return errs
}

type ChildResourceWithId[ParentT Resource, ChildT WithChildCommonParams] struct {
	ResourceId string
	Args       ChildT
	Parent     ParentT
}

func (r ResourceWithId[T]) NewValidationError(err string, field string) validate.ValidationError {
	return validate.ValidationError{
		ErrorMessage: err,
		ResourceId:   r.ResourceId,
		FieldName:    field,
	}
}

func (r ChildResourceWithId[A, B]) NewValidationError(err string, field string) validate.ValidationError {
	return validate.ValidationError{
		ErrorMessage: err,
		ResourceId:   r.ResourceId,
		FieldName:    field,
	}
}

func (r ChildResourceWithId[A, B]) GetResourceId() string {
	return r.ResourceId
}

func (r *ChildResourceWithId[A, B]) GetCloud() commonpb.CloudProvider {
	return r.Parent.GetCloud()
}

func (r *ChildResourceWithId[A, B]) GetCloudSpecificLocation() string {
	return r.Parent.GetCloudSpecificLocation()
}
