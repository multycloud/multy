package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/resources/types"
)

type GcpKubernetesNodePool struct {
	*types.KubernetesNodePool
}

func InitKubernetesNodePool(r *types.KubernetesNodePool) resources.ResourceTranslator[*resourcespb.KubernetesNodePoolResource] {
	return GcpKubernetesNodePool{r}
}

func (r GcpKubernetesNodePool) FromState(state *output.TfState) (*resourcespb.KubernetesNodePoolResource, error) {
	out := &resourcespb.KubernetesNodePoolResource{
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

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[kubernetes_node_pool.GoogleContainerNodePool](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		numZones := len(stateResource.NodeLocations)
		if numZones == 0 {
			numZones = 3
		}
		out.Name = stateResource.Name

		out.StartingNodeCount = int32(stateResource.InitialNodeCount * numZones)
		if len(stateResource.Autoscaling) == 0 {
			out.MaxNodeCount = 0
			out.MinNodeCount = 0
		} else {
			out.MaxNodeCount = int32(stateResource.Autoscaling[0].MaxNodeCount * numZones)
			out.MinNodeCount = int32(stateResource.Autoscaling[0].MinNodeCount * numZones)
		}

		if len(stateResource.NodeConfig) == 0 {
			out.Labels = map[string]string{}
			out.DiskSizeGb = 0
		} else {
			out.Labels = stateResource.NodeConfig[0].Labels
			out.DiskSizeGb = int64(stateResource.NodeConfig[0].DiskSizeGb)
			if len(r.Args.GetGcpOverride().GetMachineType()) > 0 {
				out.GcpOverride.MachineType = stateResource.NodeConfig[0].MachineType
			} else {
				out.VmSize = common.ParseVmSize(stateResource.NodeConfig[0].MachineType, common.GCP)
			}
		}

		out.GcpOutputs = &resourcespb.KubernetesNodePoolGcpOutputs{
			GkeNodePoolId: stateResource.SelfLink,
		}
	} else {
		statuses["gcp_container_node_pool"] = commonpb.ResourceStatus_NEEDS_CREATE

	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r GcpKubernetesNodePool) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	clusterId, err := resources.GetMainOutputId(GcpKubernetesCluster{r.KubernetesCluster})
	if err != nil {
		return nil, err
	}
	var size string
	if r.Args.GetGcpOverride().GetMachineType() != "" {
		size = r.Args.GetGcpOverride().GetMachineType()
	} else {
		size, err = common.GetVmSize(r.Args.VmSize, r.GetCloud())
		if err != nil {
			return nil, err
		}
	}

	numZones := 3
	if len(r.Args.AvailabilityZone) > 0 {
		numZones = len(r.Args.AvailabilityZone)
	}

	var zones []string
	for _, zone := range r.Args.AvailabilityZone {
		availabilityZone, err := common.GetAvailabilityZone(r.KubernetesCluster.GetLocation(), int(zone), r.GetCloud())
		if err != nil {
			return nil, err
		}
		zones = append(zones, availabilityZone)
	}

	var tags []string
	tags = append(tags, GcpVirtualNetwork{r.Subnet.VirtualNetwork}.getVnTag())
	tags = append(tags, GcpSubnet{r.Subnet}.getNetworkTags()...)

	nodePool := &kubernetes_node_pool.GoogleContainerNodePool{
		GcpResource:      common.NewGcpResource(r.ResourceId, r.Args.Name, r.KubernetesCluster.Args.GetGcpOverride().GetProject()),
		Cluster:          clusterId,
		NodeLocations:    zones,
		InitialNodeCount: int(r.Args.StartingNodeCount) / numZones,
		Autoscaling: []kubernetes_node_pool.GoogleContainerNodePoolAutoScaling{{
			MinNodeCount: int(r.Args.MinNodeCount) / numZones,
			MaxNodeCount: int(r.Args.MaxNodeCount) / numZones,
		}},
		NodeConfig: []kubernetes_node_pool.GoogleContainerNodeConfig{{
			DiskSizeGb:  int(r.Args.DiskSizeGb),
			Labels:      r.Args.Labels,
			MachineType: size,
			Tags:        tags,
			// Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
			ServiceAccount: fmt.Sprintf("%s.%s.email", output.GetResourceName(iam.GoogleServiceAccount{}), r.KubernetesCluster.ResourceId),
			OAuthScopes:    []string{"https://www.googleapis.com/auth/cloud-platform"},
		}},
	}

	return []output.TfBlock{nodePool}, nil

}

func (r GcpKubernetesNodePool) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_node_pool.GoogleContainerNodePool{}), nil
}
