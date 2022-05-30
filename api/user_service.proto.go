package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/userpb"
	"github.com/multycloud/multy/api/service_context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
)

type UserService struct {
	proto.UnimplementedMultyUserServiceServer
	*service_context.UserServiceContext
}

func CreateUserServer(serviceContext *service_context.UserServiceContext) UserService {
	return UserService{proto.UnimplementedMultyUserServiceServer{}, serviceContext}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.User, error) {
	_, err := mail.ParseAddress(req.EmailAddress)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, fmt.Sprintf("invalid email address: %s", req.EmailAddress)).Err()
	}
	apiKey, err := s.UserServiceContext.Database.CreateUser(ctx, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	return &userpb.User{
		UserId: req.EmailAddress,
		ApiKey: apiKey,
	}, nil
}
