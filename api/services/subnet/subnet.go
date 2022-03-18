package subnet

import (
	"context"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type SubnetService struct {
	s services.Service[*resources.CloudSpecificSubnetArgs, *resources.SubnetResource]
}

func (s SubnetService) Convert(resourceId string, args []*resources.CloudSpecificSubnetArgs) *resources.SubnetResource {
	var result []*resources.CloudSpecificSubnetResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificSubnetResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
			AvailabilityZone: r.AvailabilityZone,
			VirtualNetworkId: r.VirtualNetworkId,
		})
	}

	return &resources.SubnetResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s SubnetService) NewArg() *resources.CloudSpecificSubnetArgs {
	return &resources.CloudSpecificSubnetArgs{}
}

func (s SubnetService) Nil() *resources.SubnetResource {
	return nil
}

func (s SubnetService) Create(ctx context.Context, in *resources.CreateSubnetRequest) (*resources.SubnetResource, error) {
	return s.s.Create(ctx, in)
}

func (s SubnetService) Update(ctx context.Context, in *resources.UpdateSubnetRequest) (*resources.SubnetResource, error) {
	return s.s.Update(ctx, in)
}

func (s SubnetService) Delete(ctx context.Context, in *resources.DeleteSubnetRequest) (*common.Empty, error) {
	return s.s.Delete(ctx, in)
}

func (s SubnetService) Read(ctx context.Context, in *resources.ReadSubnetRequest) (*resources.SubnetResource, error) {
	return s.s.Read(ctx, in)
}

func NewSubnetServiceService(database *db.Database) SubnetService {
	vn := SubnetService{
		s: services.Service[*resources.CloudSpecificSubnetArgs, *resources.SubnetResource]{
			Db:         database,
			Converters: nil,
		},
	}
	vn.s.Converters = &vn
	return vn
}
