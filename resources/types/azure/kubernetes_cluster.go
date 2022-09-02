package azure_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	"github.com/multycloud/multy/resources/types"
)

type AzureKubernetesCluster struct {
	*types.KubernetesCluster
}

func InitKubernetesCluster(r *types.KubernetesCluster) resources.ResourceTranslator[*resourcespb.KubernetesClusterResource] {
	return AzureKubernetesCluster{r}
}

func (r AzureKubernetesCluster) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.KubernetesClusterResource, error) {
	result := &resourcespb.KubernetesClusterResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:             r.Args.Name,
		ServiceCidr:      r.Args.ServiceCidr,
		VirtualNetworkId: r.Args.VirtualNetworkId,
		Version:          r.Args.Version,
		GcpOverride:      r.Args.GcpOverride,
		Endpoint:         "dryrun",
	}

	result.DefaultNodePool = AzureKubernetesNodePool{r.DefaultNodePool}.translateToResource()

	if flags.DryRun {
		return result, nil
	}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if cluster, exists, err := output.MaybeGetParsedById[kubernetes_service.AzureAksCluster](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		result.Name = cluster.Name
		if len(cluster.NetworkProfile) == 0 {
			result.ServiceCidr = ""
		} else {
			result.ServiceCidr = cluster.NetworkProfile[0].ServiceCidr
		}

		result.Endpoint = cluster.KubeConfig[0].Host
		result.CaCertificate = cluster.KubeConfig[0].ClusterCaCertificate
		result.KubeConfigRaw = cluster.KubeConfigRaw
		result.Version = cluster.KubernetesVersion
		if len(cluster.DefaultNodePoolOut) == 0 {
			result.DefaultNodePool = nil
		} else {
			AzureKubernetesNodePool{r.DefaultNodePool}.parseNodePoolResource(result.DefaultNodePool, cluster.DefaultNodePoolOut[0])
		}

		result.AzureOutputs = &resourcespb.KubernetesClusterAzureOutputs{
			AksClusterId: cluster.ResourceId,
		}
		output.AddToStatuses(statuses, "azure_kubernetes_cluster", output.MaybeGetPlannedChageById[kubernetes_service.AzureAksCluster](plan, r.ResourceId))
	} else {
		statuses["azure_kubernetes_cluster"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		result.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return result, nil
}

func (r AzureKubernetesCluster) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	defaultPool, err := AzureKubernetesNodePool{r.DefaultNodePool}.translateAzNodePool()
	if err != nil {
		return nil, err
	}
	defaultPool.Name = defaultPool.AzResource.Name
	defaultPool.AzResource = nil
	defaultPool.ClusterId = ""

	return []output.TfBlock{
		&kubernetes_service.AzureAksCluster{
			AzResource:        common.NewAzResource(r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId), r.GetCloudSpecificLocation()),
			DefaultNodePool:   defaultPool,
			DnsPrefix:         common.UniqueId(r.Args.Name, "aks", common.LowercaseAlphanumericFormatFunc),
			Identity:          []kubernetes_service.AzureIdentity{{Type: "SystemAssigned"}},
			KubernetesVersion: r.Args.Version,
			NetworkProfile: []kubernetes_service.NetworkProfile{{
				NetworkPlugin:    "azure",
				DnsServiceIp:     "10.100.0.10",
				DockerBridgeCidr: "172.17.0.1/16",
				ServiceCidr:      r.Args.ServiceCidr,
			}},
		},
	}, nil
}
func (r AzureKubernetesCluster) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_service.AzureAksCluster{}), nil
}
