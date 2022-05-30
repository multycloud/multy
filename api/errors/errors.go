package errors

import (
	"context"
	"fmt"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/validate"
	"google.golang.org/genproto/googleapis/rpc/code"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
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
		details := &pberr.ResourceValidationError{
			ResourceId:   e.ResourceId,
			ErrorMessage: e.ErrorMessage,
			FieldName:    e.FieldName,
		}

		if e.ResourceNotFound {
			details.ErrorDetails = &pberr.ResourceValidationError_NotFoundDetails_{
				NotFoundDetails: &pberr.ResourceValidationError_NotFoundDetails{
					ResourceId: e.ResourceNotFoundId,
				},
			}
		}
		st, _ = st.WithDetails(details)
	}
	return st.Err()
}

func ResourceNotFound(resourceId string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("resource with id %s not found", resourceId))
	st, _ = st.WithDetails(&pberr.ResourceNotFoundDetails{ResourceId: resourceId})
	return st.Err()
}

func UserAlreadyExists(emailAddress string) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("user with email address %s already exists", emailAddress))
	return st.Err()
}

func ResourceInUseError(resourceId string, usedByResourceId string) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("resource with id %s cannot be deleted, because it is being used by resource with id %s", resourceId, usedByResourceId))
	st, _ = st.WithDetails(&pberr.ResourceNotFoundDetails{ResourceId: resourceId})
	return st.Err()
}

func WrappingErrors[InT any, OutT any](f func(context.Context, InT) (OutT, error)) func(context.Context, InT) (OutT, error) {
	return func(ctx context.Context, in InT) (out OutT, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] server panic: %v\n", r)
				debug.PrintStack()
				err = InternalServerErrorWithMessage("server panic", fmt.Errorf("%+v", r))
			}
		}()
		out, err = f(ctx, in)
		if err != nil {
			return out, InternalServerError(err)
		}
		return out, err
	}
}
