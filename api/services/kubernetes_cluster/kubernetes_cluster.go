package kubernetes_cluster

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	common_util "github.com/multycloud/multy/util"
)

type KubernetesClusterService struct {
	Service services.Service[*resources.CloudSpecificKubernetesClusterArgs, *resources.KubernetesClusterResource]
}

func (s KubernetesClusterService) Convert(resourceId string, args []*resources.CloudSpecificKubernetesClusterArgs, state *output.TfState) (*resources.KubernetesClusterResource, error) {
	var result []*resources.CloudSpecificKubernetesClusterResource
	for _, r := range args {
		endpoint, err := getEndpoint(resourceId, state, r.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
		result = append(result, &resources.CloudSpecificKubernetesClusterResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			Name:             r.Name,
			SubnetIds:        r.SubnetIds,

			Endpoint: endpoint,
		})
	}

	return &resources.KubernetesClusterResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func getEndpoint(resourceId string, state *output.TfState, cloud common.CloudProvider) (string, error) {
	rId := common_util.GetTfResourceId(resourceId, cloud.String())
	switch cloud {
	case common.CloudProvider_AWS:
		values, err := state.GetValues(kubernetes_service.AwsEksCluster{}, rId)
		if err != nil {
			return "", err
		}
		return values["endpoint"].(string), nil
	case common.CloudProvider_AZURE:
		values, err := state.GetValues(kubernetes_service.AzureEksCluster{}, rId)
		if err != nil {
			return "", err
		}
		return values["kube_config"].([]interface{})[0].(map[string]interface{})["host"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
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
