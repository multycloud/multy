package services

import (
	"context"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources/output"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
)

type CreateRequest[Arg proto.Message] interface {
	GetResource() Arg
	proto.Message
}

type UpdateRequest[Arg proto.Message] interface {
	GetResource() Arg
	WithResourceId
}

type WithResourceId interface {
	GetResourceId() string
	proto.Message
}

type Service[Arg proto.Message, OutT proto.Message] struct {
	Db           *db.Database
	Converters   converter.ResourceConverters[Arg, OutT]
	ResourceName string
}

func (s Service[Arg, OutT]) updateErrorMetric(err error, method string) {
	if err != nil {
		go s.Db.AwsClient.UpdateErrorMetric(s.ResourceName, method, errors.ErrorCode(err))
	}
}

func (s Service[Arg, OutT]) Create(ctx context.Context, in CreateRequest[Arg]) (out OutT, err error) {
	defer func() { s.updateErrorMetric(err, "create") }()
	return errors.WrappingErrors(s.create)(ctx, in)
}

func (s Service[Arg, OutT]) create(ctx context.Context, in CreateRequest[Arg]) (OutT, error) {
	go s.Db.AwsClient.UpdateQPSMetric(s.ResourceName, "create")

	log.Println("Service create")
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return *new(OutT), err
	}

	userId, err := s.Db.GetUserId(ctx, key)
	if err != nil {
		return *new(OutT), err
	}

	lock, err := s.Db.LockConfig(ctx, userId)
	if err != nil {
		return *new(OutT), err
	}
	defer s.Db.UnlockConfig(ctx, lock)

	c, err := s.Db.LoadUserConfig(userId, lock)
	if err != nil {
		return *new(OutT), err
	}
	resource, err := util.InsertIntoConfig(in.GetResource(), c)
	if err != nil {
		return *new(OutT), err
	}

	log.Printf("Deploying %s\n", resource.ResourceId)
	_, err = deploy.Deploy(ctx, c, nil, resource)
	if err != nil {
		return *new(OutT), err
	}

	err = s.Db.StoreUserConfig(c, lock)
	if err != nil {
		return *new(OutT), err
	}
	return s.readFromConfig(ctx, c, &resourcespb.ReadVirtualNetworkRequest{ResourceId: resource.ResourceId}, false)
}

func (s Service[Arg, OutT]) Read(ctx context.Context, in WithResourceId) (out OutT, err error) {
	defer func() { s.updateErrorMetric(err, "read") }()
	return errors.WrappingErrors(s.read)(ctx, in)
}

func (s Service[Arg, OutT]) read(ctx context.Context, in WithResourceId) (OutT, error) {
	go s.Db.AwsClient.UpdateQPSMetric(s.ResourceName, "read")

	log.Printf("Service read: %s\n", in.GetResourceId())
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return *new(OutT), err
	}

	userId, err := s.Db.GetUserId(ctx, key)
	if err != nil {
		return *new(OutT), err
	}

	c, err := s.Db.LoadUserConfig(userId, nil)
	if err != nil {
		return *new(OutT), err
	}

	_, err = deploy.EncodeAndStoreTfFile(ctx, c, nil, nil, true)
	if err != nil {
		return *new(OutT), err
	}

	return s.readFromConfig(ctx, c, in, true)
}

type stateParser[OutT any] interface {
	FromState(state *output.TfState) (OutT, error)
}

func (s Service[Arg, OutT]) readFromConfig(ctx context.Context, c *configpb.Config, in WithResourceId, readonly bool) (OutT, error) {
	allResources, err := deploy.GetResources(c)
	if err != nil {
		return *new(OutT), err
	}
	for _, r := range c.Resources {
		if r.ResourceId == in.GetResourceId() {
			converted, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
			if err != nil {
				return *new(OutT), err
			}
			err = deploy.MaybeInit(ctx, c.UserId, readonly)
			if err != nil {
				return *new(OutT), err
			}
			state, err := deploy.GetState(ctx, c.UserId)
			if err != nil {
				return *new(OutT), err
			}
			if parser, ok := allResources.ResourceMap[r.ResourceId].(stateParser[OutT]); ok && !flags.DryRun {
				return parser.FromState(state)
			}
			return s.Converters.Convert(in.GetResourceId(), converted.(Arg), state)
		}
	}

	return *new(OutT), errors.ResourceNotFound(in.GetResourceId())
}

func (s Service[Arg, OutT]) Update(ctx context.Context, in UpdateRequest[Arg]) (out OutT, err error) {
	defer func() { s.updateErrorMetric(err, "update") }()
	return errors.WrappingErrors(s.update)(ctx, in)
}

func (s Service[Arg, OutT]) update(ctx context.Context, in UpdateRequest[Arg]) (OutT, error) {
	go s.Db.AwsClient.UpdateQPSMetric(s.ResourceName, "update")

	log.Printf("Service update: %s\n", in.GetResourceId())
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return *new(OutT), err
	}
	userId, err := s.Db.GetUserId(ctx, key)
	if err != nil {
		return *new(OutT), err
	}
	lock, err := s.Db.LockConfig(ctx, userId)
	if err != nil {
		return *new(OutT), err
	}
	defer s.Db.UnlockConfig(ctx, lock)

	c, err := s.Db.LoadUserConfig(userId, lock)
	if err != nil {
		return *new(OutT), err
	}
	previousResource := find(c, in.GetResourceId())
	err = util.UpdateInConfig(c, in.GetResourceId(), in.GetResource())
	if err != nil {
		return *new(OutT), err
	}

	_, err = deploy.Deploy(ctx, c, previousResource, find(c, in.GetResourceId()))
	if err != nil {
		return *new(OutT), err
	}

	err = s.Db.StoreUserConfig(c, lock)
	if err != nil {
		return *new(OutT), err
	}
	return s.readFromConfig(ctx, c, in, false)
}

func (s Service[Arg, OutT]) Delete(ctx context.Context, in WithResourceId) (_ *commonpb.Empty, err error) {
	defer func() { s.updateErrorMetric(err, "delete") }()
	return errors.WrappingErrors(s.delete)(ctx, in)
}

func (s Service[Arg, OutT]) delete(ctx context.Context, in WithResourceId) (*commonpb.Empty, error) {
	go s.Db.AwsClient.UpdateQPSMetric(s.ResourceName, "delete")

	log.Printf("Service delete: %s\n", in.GetResourceId())
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return nil, err
	}
	userId, err := s.Db.GetUserId(ctx, key)
	if err != nil {
		return nil, err
	}
	lock, err := s.Db.LockConfig(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer s.Db.UnlockConfig(ctx, lock)
	c, err := s.Db.LoadUserConfig(userId, lock)
	if err != nil {
		return nil, err
	}
	previousResource := find(c, in.GetResourceId())
	err = util.DeleteResourceFromConfig(c, in.GetResourceId())
	if err != nil {
		return nil, err
	}

	_, err = deploy.Deploy(ctx, c, previousResource, nil)
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.InvalidArgument {
			for _, details := range s.Details() {
				v := details.(*pberr.ResourceValidationError)
				if v.GetNotFoundDetails() != nil && v.GetNotFoundDetails().ResourceId == in.GetResourceId() {
					return nil, errors.ResourceInUseError(in.GetResourceId(), v.ResourceId)
				}
			}
		}
		return nil, err
	}

	err = s.Db.StoreUserConfig(c, lock)
	if err != nil {
		return nil, err
	}
	return &commonpb.Empty{}, nil
}

func find(c *configpb.Config, resourceId string) *configpb.Resource {
	for _, r := range c.Resources {
		if r.ResourceId == resourceId {
			return r
		}
	}
	return nil
}
