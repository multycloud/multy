package services

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"google.golang.org/protobuf/proto"
)

type CreateRequest[Arg proto.Message] interface {
	GetResources() []Arg
	proto.Message
}

type UpdateRequest[Arg proto.Message] interface {
	GetResources() []Arg
	WithResourceId
}

type WithResourceId interface {
	GetResourceId() string
	proto.Message
}

type Service[Arg proto.Message, OutT proto.Message] struct {
	Db         *db.Database
	Converters converter.ResourceConverters[Arg, OutT]
}

func (s Service[Arg, OutT]) Create(ctx context.Context, in CreateRequest[Arg]) (OutT, error) {
	fmt.Println("Service create")
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return s.Converters.Nil(), err
	}
	c, err := s.Db.LoadUserConfig(userId)
	if err != nil {
		return s.Converters.Nil(), err
	}
	resource, err := util.InsertIntoConfig(in.GetResources(), c)
	if err != nil {
		return s.Converters.Nil(), err
	}

	fmt.Printf("Deploying %s\n", resource.ResourceId)
	err = deploy.Deploy(c, resource.ResourceId)
	if err != nil {
		return s.Converters.Nil(), err
	}

	err = s.Db.StoreUserConfig(c)
	if err != nil {
		return s.Converters.Nil(), err
	}
	return s.Read(ctx, &resources.ReadVirtualNetworkRequest{ResourceId: resource.ResourceId})
}

func (s Service[Arg, OutT]) Read(ctx context.Context, in WithResourceId) (OutT, error) {
	fmt.Printf("Service read: %s\n", in.GetResourceId())
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return s.Converters.Nil(), err
	}

	c, err := s.Db.LoadUserConfig(userId)
	if err != nil {
		return s.Converters.Nil(), err
	}
	for _, r := range c.Resources {
		if r.ResourceId == in.GetResourceId() {
			var convertedArgs []Arg
			for _, arg := range r.ResourceArgs.ResourceArgs {
				converted := s.Converters.NewArg()
				err := arg.UnmarshalTo(converted)
				if err != nil {
					return s.Converters.Nil(), err
				}
				convertedArgs = append(convertedArgs, converted)
			}
			return s.Converters.Convert(in.GetResourceId(), convertedArgs), nil
		}
	}

	return s.Converters.Nil(), fmt.Errorf("resource with id %s not found", in.GetResourceId())
}

func (s Service[Arg, OutT]) Update(ctx context.Context, in UpdateRequest[Arg]) (OutT, error) {
	fmt.Printf("Service update: %s\n", in.GetResourceId())
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return s.Converters.Nil(), err
	}
	c, err := s.Db.LoadUserConfig(userId)
	if err != nil {
		return s.Converters.Nil(), err
	}
	err = util.UpdateInConfig(c, in.GetResourceId(), in.GetResources())
	if err != nil {
		return s.Converters.Nil(), err
	}

	err = deploy.Deploy(c, in.GetResourceId())
	if err != nil {
		return s.Converters.Nil(), err
	}

	err = s.Db.StoreUserConfig(c)
	if err != nil {
		return s.Converters.Nil(), err
	}
	return s.Read(ctx, in)
}

func (s Service[Arg, OutT]) Delete(ctx context.Context, in WithResourceId) (*common.Empty, error) {
	fmt.Printf("Service delete: %s\n", in.GetResourceId())
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}
	c, err := s.Db.LoadUserConfig(userId)
	if err != nil {
		return nil, err
	}
	err = util.DeleteResourceFromConfig(c, in.GetResourceId())
	if err != nil {
		return nil, err
	}

	err = deploy.Deploy(c, in.GetResourceId())
	if err != nil {
		return nil, err
	}

	err = s.Db.StoreUserConfig(c)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
