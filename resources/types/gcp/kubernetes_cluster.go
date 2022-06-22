package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
)

type GcpKubernetesCluster struct {
	*types.KubernetesCluster
}

func InitKubernetesCluster(r *types.KubernetesCluster) resources.ResourceTranslator[*resourcespb.KubernetesClusterResource] {
	return GcpKubernetesCluster{r}
}

func (r GcpKubernetesCluster) FromState(state *output.TfState) (*resourcespb.KubernetesClusterResource, error) {
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
	}
	result.Endpoint = "dryrun"
	if !flags.DryRun {
		//cluster, err := output.GetParsedById[kubernetes_service.GoogleContainerCluster](state, r.ResourceId)
		//if err != nil {
		//	return nil, err
		//}
		//result.Endpoint = cluster.KubeConfig[0].Host
		//result.CaCertificate = cluster.KubeConfig[0].ClusterCaCertificate
		//result.KubeConfigRaw = cluster.KubeConfigRaw
	}

	return result, nil
}

func (r GcpKubernetesCluster) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		&kubernetes_service.GoogleContainerCluster{
			RemoveDefaultNodePool: true,
			InitialNodeCount:      1,
			Network:               fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
			IpAllocationPolicy: kubernetes_service.GoogleContainerClusterIpAllocationPolicy{
				ServicesIpv4CidrBlock: r.Args.ServiceCidr,
			},
		}}, nil
}
func (r GcpKubernetesCluster) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_service.GoogleContainerCluster{}), nil
}
