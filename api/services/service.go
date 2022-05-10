package services

import (
	"context"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	pberr "github.com/multycloud/multy/api/proto/errorspb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
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
	ResourceName string
}

func NewService[Arg proto.Message, OutT proto.Message](resourceName string, db *db.Database) Service[Arg, OutT] {
	return Service[Arg, OutT]{ResourceName: resourceName, Db: db}
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

	c, err := s.getConfig(userId, lock)
	if err != nil {
		return *new(OutT), err
	}
	resource, err := c.CreateResource(in.GetResource())
	if err != nil {
		return *new(OutT), err
	}

	log.Printf("[INFO] Deploying %s\n", resource.GetResourceId())
	_, err = deploy.Deploy(ctx, c, nil, resource)
	if err != nil {
		return *new(OutT), err
	}

	err = s.saveConfig(c, lock)
	if err != nil {
		return *new(OutT), err
	}
	return s.readFromConfig(ctx, c, &resourcespb.ReadVirtualNetworkRequest{ResourceId: resource.GetResourceId()}, false)
}

func (s Service[Arg, OutT]) getConfig(userId string, lock *db.ConfigLock) (*resources.MultyConfig, error) {
	c, err := s.Db.LoadUserConfig(userId, lock)
	if err != nil {
		return nil, err
	}
	mconfig, err := resources.LoadConfig(c, types.Metadatas)
	if err != nil {
		return nil, err
	}
	return mconfig, err
}

func (s Service[Arg, OutT]) saveConfig(c *resources.MultyConfig, lock *db.ConfigLock) error {
	exportedConfig, err := c.ExportConfig()
	if err != nil {
		return err
	}

	return s.Db.StoreUserConfig(exportedConfig, lock)
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

	c, err := s.getConfig(userId, nil)
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

func (s Service[Arg, OutT]) readFromConfig(ctx context.Context, c *resources.MultyConfig, in WithResourceId, readonly bool) (OutT, error) {
	for _, r := range c.Resources.GetAll() {
		if r.GetResourceId() == in.GetResourceId() {
			err := deploy.MaybeInit(ctx, c.GetUserId(), readonly)
			if err != nil {
				return *new(OutT), err
			}
			state, err := deploy.GetState(ctx, c.GetUserId(), readonly)
			if err != nil {
				return *new(OutT), err
			}
			if parser, ok := r.(stateParser[OutT]); ok && !flags.DryRun {
				return parser.FromState(state)
			}

			out, err := r.GetMetadata().ReadFromState(r, state)
			if err != nil {
				return *new(OutT), err
			}
			return out.(OutT), err
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

	c, err := s.getConfig(userId, lock)
	if err != nil {
		return *new(OutT), err
	}

	r, err := c.UpdateResource(in.GetResourceId(), in.GetResource())
	if err != nil {
		return *new(OutT), err
	}

	_, err = deploy.Deploy(ctx, c, r, r)
	if err != nil {
		return *new(OutT), err
	}

	err = s.saveConfig(c, lock)
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
	c, err := s.getConfig(userId, lock)
	if err != nil {
		return nil, err
	}
	previousResource, err := c.DeleteResource(in.GetResourceId())
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

	err = s.saveConfig(c, lock)
	if err != nil {
		return nil, err
	}
	return &commonpb.Empty{}, nil
}
