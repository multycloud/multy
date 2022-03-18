package virtual_network

import (
	"context"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type VnService struct {
	s services.Service[*resources.CloudSpecificVirtualNetworkArgs, *resources.VirtualNetworkResource]
}

func (s VnService) Convert(resourceId string, args []*resources.CloudSpecificVirtualNetworkArgs) *resources.VirtualNetworkResource {
	var result []*resources.CloudSpecificVirtualNetworkResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificVirtualNetworkResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
		})
	}

	return &resources.VirtualNetworkResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s VnService) NewArg() *resources.CloudSpecificVirtualNetworkArgs {
	return &resources.CloudSpecificVirtualNetworkArgs{}
}

func (s VnService) Nil() *resources.VirtualNetworkResource {
	return nil
}

func (s VnService) Create(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.s.Create(ctx, in)
}

func (s VnService) Update(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.s.Update(ctx, in)
}

func (s VnService) Delete(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	return s.s.Delete(ctx, in)
}

func (s VnService) Read(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.s.Read(ctx, in)
}

func NewVnService(database *db.Database) VnService {
	vn := VnService{
		s: services.Service[*resources.CloudSpecificVirtualNetworkArgs, *resources.VirtualNetworkResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.s.Converters = &vn
	return vn
}
