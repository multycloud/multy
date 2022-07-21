package aws_resources

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

type AwsKubernetesNodePool struct {
	*types.KubernetesNodePool
}

func InitKubernetesNodePool(r *types.KubernetesNodePool) resources.ResourceTranslator[*resourcespb.KubernetesNodePoolResource] {
	return AwsKubernetesNodePool{r}
}

func (r AwsKubernetesNodePool) FromState(state *output.TfState) (*resourcespb.KubernetesNodePoolResource, error) {
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

	stateResource, err := output.GetParsedById[kubernetes_node_pool.AwsKubernetesNodeGroup](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	out.AwsOutputs = &resourcespb.KubernetesNodePoolAwsOutputs{
		EksNodePoolId: stateResource.Arn,
	}

	if stateResource, exists, err := output.MaybeGetParsedById[iam.AwsIamRole](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.IamRoleArn = stateResource.Arn
	}

	return out, nil

}

func (r AwsKubernetesNodePool) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	var subnetIds []string
	if len(r.Args.AvailabilityZone) == 0 {
		subnetIds = AwsSubnet{r.Subnet}.GetSubnetIds()
	} else {
		for _, zone := range r.Args.AvailabilityZone {
			id, err := AwsSubnet{r.Subnet}.GetSubnetId(zone)
			if err != nil {
				return nil, err
			}
			subnetIds = append(subnetIds, id)
		}
	}

	var instanceTypes []string
	if r.Args.AwsOverride.GetInstanceTypes() != nil {
		instanceTypes = r.Args.AwsOverride.GetInstanceTypes()
	} else {
		instanceTypes = []string{common.VMSIZE[r.Args.VmSize][r.GetCloud()]}
	}

	roleName := fmt.Sprintf("multy-k8nodepool-%s-%s-role", r.KubernetesCluster.Args.Name, r.Args.Name)
	role := iam.AwsIamRole{
		AwsResource:      common.NewAwsResource(r.ResourceId, roleName),
		Name:             roleName,
		AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
	}
	clusterId, err := resources.GetMainOutputId(AwsKubernetesCluster{r.KubernetesCluster})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		&role,
		iam.AwsIamRolePolicyAttachment{
			AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSWorkerNodePolicy")),
			Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		},
		iam.AwsIamRolePolicyAttachment{
			AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKS_CNI_Policy")),
			Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
		},
		iam.AwsIamRolePolicyAttachment{
			AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEC2ContainerRegistryReadOnly")),
			Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
			PolicyArn:   "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
		},
		&kubernetes_node_pool.AwsKubernetesNodeGroup{
			AwsResource:   common.NewAwsResourceWithIdOnly(r.ResourceId),
			ClusterName:   clusterId,
			NodeGroupName: r.Args.Name,
			NodeRoleArn:   fmt.Sprintf("aws_iam_role.%s.arn", r.ResourceId),
			SubnetIds:     subnetIds,
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

func (r AwsKubernetesNodePool) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_node_pool.AwsKubernetesNodeGroup{}), nil
}
