package errors

import (
	"fmt"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/validate"
	"google.golang.org/genproto/googleapis/rpc/code"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PermissionDenied(msg string) error {
	return status.ErrorProto(&spb.Status{
		Code:    int32(code.Code_PERMISSION_DENIED),
		Message: msg,
		Details: nil,
	})
}

func ErrorCode(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Code().String()
	}

	return codes.Internal.String()
}

func InternalServerError(err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	return status.New(codes.Internal, err.Error()).Err()
}

func DeployError(err error) error {
	st := status.New(codes.FailedPrecondition, "error while deploying resources")
	st, _ = st.WithDetails(&pberr.DeploymentErrorDetails{ErrorMessage: err.Error()})
	return st.Err()
}

func InternalServerErrorWithMessage(msg string, err error) error {
	st := status.New(codes.Internal, msg)
	st, _ = st.WithDetails(&pberr.InternalErrorDetails{ErrorMessage: err.Error()})
	return st.Err()
}

func ValidationErrors(errs []validate.ValidationError) error {
	st := status.New(codes.InvalidArgument, fmt.Sprintf("%d validation errors found", len(errs)))
	for _, e := range errs {
		st, _ = st.WithDetails(&pberr.ResourceValidationError{
			ResourceId:   e.ResourceId,
			ErrorMessage: e.ErrorMessage,
			FieldName:    e.FieldName,
		})
	}
	return st.Err()
}

func ResourceNotFound(resourceId string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("resource with id %s not found", resourceId))
	st, _ = st.WithDetails(&pberr.ResourceNotFoundDetails{ResourceId: resourceId})
	return st.Err()
}
