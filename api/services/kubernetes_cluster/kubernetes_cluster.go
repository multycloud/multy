package kubernetes_cluster

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/services/kubernetes_node_pool"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
)

type KubernetesClusterService struct {
	Service services.Service[*resourcespb.KubernetesClusterArgs, *resourcespb.KubernetesClusterResource]
}

func (s KubernetesClusterService) Convert(resourceId string, args *resourcespb.KubernetesClusterArgs, state *output.TfState) (*resourcespb.KubernetesClusterResource, error) {
	var err error
	endpoint := "dryrun"
	if !flags.DryRun {
		endpoint, err = getEndpoint(resourceId, state, args.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
	}

	defaultNodePool, err := kubernetes_node_pool.ConvertNodePool("", args.DefaultNodePool)
	if err != nil {
		return nil, err
	}

	return &resourcespb.KubernetesClusterResource{
		CommonParameters: util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:             args.Name,
		SubnetIds:        args.SubnetIds,
		Endpoint:         endpoint,
		DefaultNodePool:  defaultNodePool,
	}, nil
}

func getEndpoint(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(kubernetes_service.AwsEksCluster{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["endpoint"].(string), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(kubernetes_service.AzureEksCluster{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["kube_config"].([]interface{})[0].(map[string]interface{})["host"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func NewKubernetesClusterService(database *db.Database) KubernetesClusterService {
	ni := KubernetesClusterService{
		Service: services.Service[*resourcespb.KubernetesClusterArgs, *resourcespb.KubernetesClusterResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "kubernetes_cluster",
		},
	}
	ni.Service.Converters = &ni
	return ni
}
