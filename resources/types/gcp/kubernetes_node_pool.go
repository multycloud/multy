package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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
		Labels:            r.Args.Labels,
		GcpOverride:       r.Args.GcpOverride,
		AzureOverride:     r.Args.AzureOverride,
	}, nil
}

func (r GcpKubernetesNodePool) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	subnetId, err := resources.GetMainOutputId(GcpSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}

	var instanceTypes []string
	if r.Args.GcpOverride.GetInstanceTypes() != nil {
		instanceTypes = r.Args.GcpOverride.GetInstanceTypes()
	} else {
		instanceTypes = []string{common.VMSIZE[r.Args.VmSize][r.GetCloud()]}
	}

	roleName := fmt.Sprintf("multy-k8nodepool-%s-%s-role", r.KubernetesCluster.Args.Name, r.Args.Name)
	role := iam.GcpIamRole{
		GcpResource:      common.NewGcpResource(r.ResourceId, roleName),
		Name:             roleName,
		AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonGcp.com"),
	}
	clusterId, err := resources.GetMainOutputId(GcpKubernetesCluster{r.KubernetesCluster})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		&role,
		iam.GcpIamRolePolicyAttachment{
			GcpResource: common.NewGcpResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSWorkerNodePolicy")),
			Role:        fmt.Sprintf("Gcp_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:Gcp:iam::Gcp:policy/AmazonEKSWorkerNodePolicy",
		},
		iam.GcpIamRolePolicyAttachment{
			GcpResource: common.NewGcpResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKS_CNI_Policy")),
			Role:        fmt.Sprintf("Gcp_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:Gcp:iam::Gcp:policy/AmazonEKS_CNI_Policy",
		},
		iam.GcpIamRolePolicyAttachment{
			GcpResource: common.NewGcpResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEC2ContainerRegistryReadOnly")),
			Role:        fmt.Sprintf("Gcp_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:Gcp:iam::Gcp:policy/AmazonEC2ContainerRegistryReadOnly",
		},
		&kubernetes_node_pool.GcpKubernetesNodeGroup{
			GcpResource:   common.NewGcpResourceWithIdOnly(r.ResourceId),
			ClusterName:   clusterId,
			NodeGroupName: r.Args.Name,
			NodeRoleArn:   fmt.Sprintf("Gcp_iam_role.%s.arn", r.ResourceId),
			SubnetIds:     []string{subnetId},
			ScalingConfig: kubernetes_node_pool.ScalingConfig{
				DesiredSize: int(r.Args.StartingNodeCount),
				MaxSize:     int(r.Args.MaxNodeCount),
				MinSize:     int(r.Args.MinNodeCount),
			},
			Labels:        r.Args.Labels,
			InstanceTypes: instanceTypes,
		},
	}, nil

}

func (r GcpKubernetesNodePool) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_node_pool.GcpKubernetesNodeGroup{}), nil
}
