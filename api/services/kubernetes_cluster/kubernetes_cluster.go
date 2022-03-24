package kubernetes_cluster

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type KubernetesClusterService struct {
	Service services.Service[*resources.CloudSpecificKubernetesClusterArgs, *resources.KubernetesClusterResource]
}

func (s KubernetesClusterService) Convert(resourceId string, args []*resources.CloudSpecificKubernetesClusterArgs, state *output.TfState) (*resources.KubernetesClusterResource, error) {
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
	}, nil
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
