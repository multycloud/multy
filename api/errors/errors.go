package errors

import (
	"fmt"
	pberr "github.com/multycloud/multy/api/proto/errors"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"google.golang.org/genproto/googleapis/rpc/code"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func PermissionDenied(msg string) error {
	return status.ErrorProto(&spb.Status{
		Code:    int32(code.Code_PERMISSION_DENIED),
		Message: msg,
		Details: nil,
	})
}

func InternalServerError(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	return status.New(codes.Internal, err.Error()).Err()
}

func InternalServerErrorWithMessage(msg string, err error) error {
	st := status.New(codes.Internal, msg)
	st, _ = st.WithDetails(&pberr.InternalErrorDetails{ErrorMessage: err.Error()})
	return st.Err()
}

func ValidationErrors(errs []validate.ValidationError) error {
	return status.ErrorProto(&spb.Status{
		Code:    int32(code.Code_INVALID_ARGUMENT),
		Message: fmt.Sprintf("%d validation errors found", len(errs)),
		Details: util.MapSliceValues(errs, func(e validate.ValidationError) *anypb.Any {
			a, _ := anypb.New(&pberr.ResourceValidationError{
				ResourceId:   e.ResourceId,
				ErrorMessage: e.ErrorMessage,
				FieldName:    e.FieldName,
			})
			return a
		}),
	})
}

func ResourceNotFound(resourceId string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("resource with id %s not found", resourceId))
	st, _ = st.WithDetails(&pberr.ResourceNotFoundDetails{ResourceId: resourceId})
	return st.Err()
}
