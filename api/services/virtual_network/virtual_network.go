package virtual_network

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type VnService struct {
	Db *db.Database
}

func (s VnService) Create(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}
	resource, err := util.StoreResourceInDb(ctx, in, s.Db)
	if err != nil {
		return nil, err
	}
	c, err := s.Db.Load(userId)
	if err != nil {
		return nil, err
	}
	err = deploy.Deploy(c, resource.ResourceId)
	if err != nil {
		return nil, err
	}
	return s.Read(ctx, &resources.ReadVirtualNetworkRequest{
		ResourceId: resource.ResourceId,
	})
}

func (s VnService) Update(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = util.UpdateResourceInDb(ctx, in.ResourceId, &resources.CreateVirtualNetworkRequest{Resources: in.Resources}, s.Db)
	if err != nil {
		return nil, err
	}
	c, err := s.Db.Load(userId)
	if err != nil {
		return nil, err
	}
	err = deploy.Deploy(c, in.ResourceId)
	if err != nil {
		return nil, err
	}

	return s.Read(ctx, &resources.ReadVirtualNetworkRequest{
		ResourceId: in.ResourceId,
	})
}

func (s VnService) Delete(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}
	err = util.DeleteResourceFromDb(ctx, in.ResourceId, s.Db)
	if err != nil {
		return nil, err
	}
	c, err := s.Db.Load(userId)
	if err != nil {
		return nil, err
	}
	err = deploy.Deploy(c, in.ResourceId)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (s VnService) Read(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := util.ExtractUserId(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.Db.Load(userId)
	if err != nil {
		return nil, err
	}
	for _, r := range c.Resources {
		if r.ResourceId == in.ResourceId {
			converted, err := convertVirtualNetworks(r.Resource)
			if err != nil {
				return nil, err
			}
			return &resources.VirtualNetworkResource{
				CommonParameters: &common.CommonResourceParameters{ResourceId: r.ResourceId},
				Resources:        converted,
			}, nil
		}
	}

	return nil, fmt.Errorf("resource with id %s not found", in.ResourceId)
}

func convertVirtualNetworks(resource *any.Any) ([]*resources.CloudSpecificVirtualNetworkResource, error) {
	out := resources.CreateVirtualNetworkRequest{}
	err := resource.UnmarshalTo(&out)
	if err != nil {
		return nil, err
	}

	var result []*resources.CloudSpecificVirtualNetworkResource
	for _, r := range out.Resources {
		result = append(result, &resources.CloudSpecificVirtualNetworkResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
		})
	}

	return result, nil
}
