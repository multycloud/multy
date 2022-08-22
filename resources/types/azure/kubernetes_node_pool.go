package azure_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/resources/types"
)

type AzureKubernetesNodePool struct {
	*types.KubernetesNodePool
}

func InitKubernetesNodePool(r *types.KubernetesNodePool) resources.ResourceTranslator[*resourcespb.KubernetesNodePoolResource] {
	return AzureKubernetesNodePool{r}
}

func (r AzureKubernetesNodePool) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.KubernetesNodePoolResource, error) {
	out := r.translateToResource()

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[kubernetes_node_pool.AzureKubernetesNodePool](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		r.parseNodePoolResource(out, stateResource)
	} else {
		statuses["azure_kubernetes_node_pool"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureKubernetesNodePool) parseNodePoolResource(out *resourcespb.KubernetesNodePoolResource, stateResource *kubernetes_node_pool.AzureKubernetesNodePool) {
	out.Name = stateResource.Name
	out.DiskSizeGb = int64(stateResource.OsDiskSizeGb)
	out.MaxNodeCount = int32(stateResource.MaxSize)
	out.MinNodeCount = int32(stateResource.MinSize)
	out.Labels = stateResource.Labels
	if len(r.Args.GetAzureOverride().GetVmSize()) > 0 {
		out.AzureOverride.VmSize = stateResource.VmSize
	} else {
		out.VmSize = common.ParseVmSize(stateResource.VmSize, common.AZURE)
	}

	// its null for default node pool
	if stateResource.AzResource != nil {
		out.AzureOutputs = &resourcespb.KubernetesNodePoolAzureOutputs{
			AksNodePoolId: stateResource.ResourceId,
		}
	}

}

func (r AzureKubernetesNodePool) translateToResource() *resourcespb.KubernetesNodePoolResource {
	return &resourcespb.KubernetesNodePoolResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		Name:              r.Args.Name,
		SubnetId:          r.Args.SubnetId,
		ClusterId:         r.Args.ClusterId,
		StartingNodeCount: r.Args.StartingNodeCount,
		MinNodeCount:      r.Args.MinNodeCount,
		MaxNodeCount:      r.Args.MaxNodeCount,
		VmSize:            r.Args.VmSize,
		DiskSizeGb:        r.Args.DiskSizeGb,
		AvailabilityZone:  r.Args.AvailabilityZone,
		AwsOverride:       r.Args.AwsOverride,
		AzureOverride:     r.Args.AzureOverride,
		GcpOverride:       r.Args.GcpOverride,
		Labels:            r.Args.Labels,
	}
}

func (r AzureKubernetesNodePool) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	pool, err := r.translateAzNodePool()
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		pool,
	}, nil
}

func (r AzureKubernetesNodePool) translateAzNodePool() (*kubernetes_node_pool.AzureKubernetesNodePool, error) {
	clusterId, err := resources.GetMainOutputId(AzureKubernetesCluster{r.KubernetesCluster})
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(AzureSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}

	var vmSize string
	if r.Args.AzureOverride.GetVmSize() != "" {
		vmSize = r.Args.AzureOverride.GetVmSize()
	} else {
		vmSize = common.VMSIZE[r.Args.VmSize][r.GetCloud()]
	}

	var zones []string
	for _, zone := range r.Args.AvailabilityZone {
		availabilityZone, err := common.GetAvailabilityZone(r.KubernetesCluster.GetLocation(), int(zone), r.GetCloud())
		if err != nil {
			return nil, err
		}
		zones = append(zones, availabilityZone)
	}

	return &kubernetes_node_pool.AzureKubernetesNodePool{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			Name:              r.Args.Name,
		},
		ClusterId:              clusterId,
		NodeCount:              int(r.Args.StartingNodeCount),
		MaxSize:                int(r.Args.MaxNodeCount),
		MinSize:                int(r.Args.MinNodeCount),
		Labels:                 r.Args.Labels,
		EnableAutoScaling:      true,
		VmSize:                 vmSize,
		VirtualNetworkSubnetId: subnetId,
		Zones:                  zones,
		OsDiskSizeGb:           int(r.Args.DiskSizeGb),
	}, nil
}

func (r AzureKubernetesNodePool) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_node_pool.AzureKubernetesNodePool{}), nil
}
