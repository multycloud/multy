package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type WithCommonParams interface {
	GetCommonParameters() *commonpb.ResourceCommonArgs
	proto.Message
}

type WithChildCommonParams interface {
	GetCommonParameters() *commonpb.ChildResourceCommonArgs
	proto.Message
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

func (r *ResourceWithId[T]) ParseCloud(args T) commonpb.CloudProvider {
	return args.GetCommonParameters().CloudProvider
}

func (r *ResourceWithId[T]) GetMetadata(m ResourceMetadatas) (ResourceMetadataInterface, error) {
	converter, err := m.GetConverter(proto.MessageName(r.Args))
	if err != nil {
		return nil, err
	}
	return converter, nil
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

	if r.Args.GetCommonParameters().ResourceGroupId == "" {
		errs = append(errs, validate.ValidationError{
			ErrorMessage: "resource group id cannot be empty",
			ResourceId:   r.ResourceId,
			FieldName:    "resource_group_id",
		})
	}

	return errs
}

type ChildResourceWithId[ParentT Resource, ChildT WithChildCommonParams] struct {
	ResourceId string
	Args       ChildT
	Parent     ParentT
	Metadata   ResourceMetadataInterface
}

func NewResource[T WithCommonParams](resourceId string, args T) ResourceWithId[T] {
	return ResourceWithId[T]{
		ResourceId: resourceId,
		Args:       args,
	}
}

func NewChildResource[ParentT Resource, ChildT WithChildCommonParams](resourceId string, parent ParentT, args ChildT) ChildResourceWithId[ParentT, ChildT] {
	return ChildResourceWithId[ParentT, ChildT]{
		ResourceId: resourceId,
		Args:       args,
		Parent:     parent,
	}
}

func (r *ResourceWithId[T]) NewValidationError(err error, field string) validate.ValidationError {
	return NewError(err, r.ResourceId, field)
}

func NewError(err error, resourceId string, fieldName string) validate.ValidationError {
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

func (r *ChildResourceWithId[A, B]) NewValidationError(err error, field string) validate.ValidationError {
	return NewError(err, r.ResourceId, field)
}

func (r *ChildResourceWithId[A, B]) GetResourceId() string {
	return r.ResourceId
}

func (r *ChildResourceWithId[A, B]) GetCloud() commonpb.CloudProvider {
	return r.Parent.GetCloud()
}

func (r *ChildResourceWithId[A, B]) GetCloudSpecificLocation() string {
	return r.Parent.GetCloudSpecificLocation()
}

func (r *ChildResourceWithId[A, B]) GetMetadata(m ResourceMetadatas) (ResourceMetadataInterface, error) {
	converter, err := m.GetConverter(proto.MessageName(r.Args))
	if err != nil {
		return nil, err
	}
	return converter, nil
}
