package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strings"
)

func ConvertCommonParams(resourceId string, parameters *commonpb.ResourceCommonArgs) *commonpb.CommonResourceParameters {
	return &commonpb.CommonResourceParameters{
		ResourceId:      resourceId,
		ResourceGroupId: parameters.ResourceGroupId,
		Location:        parameters.Location,
		CloudProvider:   parameters.CloudProvider,
		NeedsUpdate:     false,
	}
}

func ConvertCommonChildParams(resourceId string, parameters *commonpb.ChildResourceCommonArgs) *commonpb.CommonChildResourceParameters {
	return &commonpb.CommonChildResourceParameters{
		ResourceId:  resourceId,
		NeedsUpdate: false,
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

func ExtractCloudCredentials(ctx context.Context) (*credspb.CloudCredentials, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	cloudCredsBin := md.Get("cloud-creds-bin")
	if len(cloudCredsBin) == 0 {
		return nil, errors.PermissionDenied(fmt.Sprintf("cloud credentials must be set"))
	}
	if len(cloudCredsBin) > 1 {
		return nil, errors.PermissionDenied(fmt.Sprintf("only expected 1 cloud creds id, found %d", len(cloudCredsBin)))
	}
	res := credspb.CloudCredentials{}
	err := proto.Unmarshal([]byte(cloudCredsBin[0]), &res)
	return &res, err
}

func InsertIntoConfig[Arg proto.Message](in Arg, c *configpb.Config) (*configpb.Resource, error) {
	args, err := convert(in)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error marhsalling resource", err)
	}
	resource := configpb.Resource{
		ResourceId:   base64.StdEncoding.EncodeToString([]byte(uuid.New().String())),
		ResourceArgs: args,
	}
	c.Resources = append(c.Resources, &resource)
	return &resource, nil
}

func convert[Arg proto.Message](in Arg) (*configpb.ResourceArgs, error) {
	args := configpb.ResourceArgs{}

	a, err := anypb.New(in)
	if err != nil {
		return nil, err
	}
	args.ResourceArgs = a
	return &args, nil
}

func DeleteResourceFromConfig(c *configpb.Config, resourceId string) error {
	if i := slices.IndexFunc(c.Resources, func(r *configpb.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		c.Resources = append(c.Resources[:i], c.Resources[i+1:]...)
	}
	return nil
}

func UpdateInConfig[Arg proto.Message](c *configpb.Config, resourceId string, in Arg) error {
	if i := slices.IndexFunc(c.Resources, func(r *configpb.Resource) bool { return r.ResourceId == resourceId }); i != -1 {
		a, err := convert(in)
		if err != nil {
			return err
		}
		c.Resources[i].ResourceArgs = a
	} else {
		return errors.ResourceNotFound(resourceId)
	}
	return nil
}
