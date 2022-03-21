package kubernetes_node_pool

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type KubernetesNodePoolService struct {
	Service services.Service[*resources.CloudSpecificKubernetesNodePoolArgs, *resources.KubernetesNodePoolResource]
}

func (s KubernetesNodePoolService) Convert(resourceId string, args []*resources.CloudSpecificKubernetesNodePoolArgs) *resources.KubernetesNodePoolResource {
	var result []*resources.CloudSpecificKubernetesNodePoolResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificKubernetesNodePoolResource{
			CommonParameters:  util.ConvertCommonParams(r.CommonParameters),
			Name:              r.Name,
			SubnetIds:         r.SubnetIds,
			ClusterId:         r.ClusterId,
			IsDefaultPool:     r.IsDefaultPool,
			StartingNodeCount: r.StartingNodeCount,
			MinNodeCount:      r.MinNodeCount,
			MaxNodeCount:      r.MaxNodeCount,
			VmSize:            r.VmSize,
			DiskSizeGb:        r.DiskSizeGb,
			Labels:            r.Labels,
		})
	}

	return &resources.KubernetesNodePoolResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s KubernetesNodePoolService) NewArg() *resources.CloudSpecificKubernetesNodePoolArgs {
	return &resources.CloudSpecificKubernetesNodePoolArgs{}
}

func (s KubernetesNodePoolService) Nil() *resources.KubernetesNodePoolResource {
	return nil
}

func NewKubernetesNodePoolService(database *db.Database) KubernetesNodePoolService {
	ni := KubernetesNodePoolService{
		Service: services.Service[*resources.CloudSpecificKubernetesNodePoolArgs, *resources.KubernetesNodePoolResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
