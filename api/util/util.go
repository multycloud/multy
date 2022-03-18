package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/config"
	"github.com/multycloud/multy/db"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertCommonParams(parameters *common.CloudSpecificResourceCommonArgs) *common.CloudSpecificCommonResourceParameters {
	return &common.CloudSpecificCommonResourceParameters{
		ResourceGroupId: parameters.ResourceGroupId,
		Location:        parameters.Location,
		CloudProvider:   parameters.CloudProvider,
		NeedsUpdate:     false,
	}
}

func ExtractUserId(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	userIds := md.Get("user_id")
	if len(userIds) == 0 {
		return "", fmt.Errorf("user id must be set")
	}
	if len(userIds) > 1 {
		return "", fmt.Errorf("only expected 1 user id, found %d", len(userIds))
	}
	return userIds[0], nil
}

func StoreResourceInDb[Arg proto.Message](ctx context.Context, in []Arg, database *db.Database) (*config.Resource, error) {
	userId, err := ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}

	args, err := convert(in)
	if err != nil {
		return nil, err
	}
	resource := config.Resource{
		ResourceId:   base64.StdEncoding.EncodeToString([]byte(uuid.New().String())),
		ResourceArgs: args,
	}

	c, err := database.Load(userId)
	if err != nil {
		return nil, err
	}
	c.Resources = append(c.Resources, &resource)
	err = database.Store(c)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func convert[Arg proto.Message](in []Arg) (*config.ResourceArgs, error) {
	args := config.ResourceArgs{}

	for _, arg := range in {
		a, err := anypb.New(arg)
		if err != nil {
			return nil, err
		}
		args.ResourceArgs = append(args.ResourceArgs, a)
	}
	return &args, nil
}

func DeleteResourceFromDb(ctx context.Context, resourceId string, database *db.Database) error {
	userId, err := ExtractUserId(ctx)
	if err != nil {
		return err
	}

	c, err := database.Load(userId)
	if err != nil {
		return err
	}
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		c.Resources = append(c.Resources[:i], c.Resources[i+1:]...)
		err = database.Store(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateResourceInDb[Arg proto.Message](ctx context.Context, resourceId string, in []Arg, database *db.Database) error {
	userId, err := ExtractUserId(ctx)
	if err != nil {
		return err
	}

	c, err := database.Load(userId)
	if err != nil {
		return err
	}
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		a, err := convert(in)
		if err != nil {
			return err
		}
		c.Resources[i].ResourceArgs = a
		err = database.Store(c)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("resource with id %s not found", resourceId)
	}
	return nil
}
