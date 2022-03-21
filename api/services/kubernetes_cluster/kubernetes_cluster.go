package kubernetes_cluster

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
)

type KubernetesClusterService struct {
	Service services.Service[*resources.CloudSpecificKubernetesClusterArgs, *resources.KubernetesClusterResource]
}

func (s KubernetesClusterService) Convert(resourceId string, args []*resources.CloudSpecificKubernetesClusterArgs) *resources.KubernetesClusterResource {
	var result []*resources.CloudSpecificKubernetesClusterResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificKubernetesClusterResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			SubnetIds:        r.SubnetIds,
		})
	}

	return &resources.KubernetesClusterResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}
}

func (s KubernetesClusterService) NewArg() *resources.CloudSpecificKubernetesClusterArgs {
	return &resources.CloudSpecificKubernetesClusterArgs{}
}

func (s KubernetesClusterService) Nil() *resources.KubernetesClusterResource {
	return nil
}

func NewKubernetesClusterService(database *db.Database) KubernetesClusterService {
	ni := KubernetesClusterService{
		Service: services.Service[*resources.CloudSpecificKubernetesClusterArgs, *resources.KubernetesClusterResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
