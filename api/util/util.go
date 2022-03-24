package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/config"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
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
		return "", errors.PermissionDenied(fmt.Sprintf("user id must be set"))
	}
	if len(userIds) > 1 {
		return "", errors.PermissionDenied(fmt.Sprintf("only expected 1 user id, found %d: %s", len(userIds), strings.Join(userIds, ", ")))
	}
	return userIds[0], nil
}

func InsertIntoConfig[Arg proto.Message](in []Arg, c *config.Config) (*config.Resource, error) {
	args, err := convert(in)
	if err != nil {
		return nil, err
	}
	resource := config.Resource{
		ResourceId:   base64.StdEncoding.EncodeToString([]byte(uuid.New().String())),
		ResourceArgs: args,
	}
	c.Resources = append(c.Resources, &resource)
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

func DeleteResourceFromConfig(c *config.Config, resourceId string) error {
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		c.Resources = append(c.Resources[:i], c.Resources[i+1:]...)
	}
	return nil
}

func UpdateInConfig[Arg proto.Message](c *config.Config, resourceId string, in []Arg) error {
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		a, err := convert(in)
		if err != nil {
			return err
		}
		c.Resources[i].ResourceArgs = a
	} else {
		return fmt.Errorf("resource with id %s not found", resourceId)
	}
	return nil
}
