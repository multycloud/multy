package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if result, err := common.GetCloudLocation(r.GetLocation(), r.GetCloud()); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}

func (r *ResourceWithId[T]) GetLocation() commonpb.Location {
	return r.Args.GetCommonParameters().GetLocation()
}

func (r *ResourceWithId[T]) GetResourceId() string {
	return r.ResourceId
}

func (r *ResourceWithId[T]) GetResourceGroupId() string {
	return r.Args.GetCommonParameters().GetResourceGroupId()
}

func (r *ResourceWithId[T]) GetCloud() commonpb.CloudProvider {
	return r.Args.GetCommonParameters().CloudProvider
}

func (r *ResourceWithId[T]) GetCommonArgs() any {
	return r.Args.GetCommonParameters()
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

func (r ResourceWithId[T]) NewValidationError(err error, field string) validate.ValidationError {
	return newError(err, r.ResourceId, field)
}

func newError(err error, resourceId string, fieldName string) validate.ValidationError {
	result := validate.ValidationError{
		ErrorMessage: err.Error(),
		ResourceId:   resourceId,
		FieldName:    fieldName,
	}
	if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
		details := s.Details()[0].(*pberr.ResourceNotFoundDetails)
		result.ResourceNotFound = true
		result.ResourceNotFoundId = details.ResourceId
	}
	return result
}

func (r ChildResourceWithId[A, B]) NewValidationError(err error, field string) validate.ValidationError {
	return newError(err, r.ResourceId, field)
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

func (r *ChildResourceWithId[A, B]) GetCommonArgs() any {
	return r.Args.GetCommonParameters()
}
