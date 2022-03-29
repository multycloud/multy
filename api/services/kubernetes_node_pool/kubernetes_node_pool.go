package kubernetes_node_pool

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type KubernetesNodePoolService struct {
	Service services.Service[*resources.KubernetesNodePoolArgs, *resources.KubernetesNodePoolResource]
}

func (s KubernetesNodePoolService) Convert(resourceId string, args *resources.KubernetesNodePoolArgs, state *output.TfState) (*resources.KubernetesNodePoolResource, error) {
	return &resources.KubernetesNodePoolResource{
		CommonParameters:  util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		Name:              args.Name,
		SubnetIds:         args.SubnetIds,
		ClusterId:         args.ClusterId,
		IsDefaultPool:     args.IsDefaultPool,
		StartingNodeCount: args.StartingNodeCount,
		MinNodeCount:      args.MinNodeCount,
		MaxNodeCount:      args.MaxNodeCount,
		VmSize:            args.VmSize,
		DiskSizeGb:        args.DiskSizeGb,
		Labels:            args.Labels,
	}, nil
}

func NewKubernetesNodePoolService(database *db.Database) KubernetesNodePoolService {
	ni := KubernetesNodePoolService{
		Service: services.Service[*resources.KubernetesNodePoolArgs, *resources.KubernetesNodePoolResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
