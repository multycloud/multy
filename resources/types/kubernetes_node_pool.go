package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

type KubernetesServiceNodePool struct {
	*resources.CommonResourceParams
	Name              string             `hcl:"name"`
	ClusterId         *KubernetesService `mhcl:"ref=cluster_id"`
	IsDefaultPool     bool               `hcl:"is_default_pool,optional"`
	SubnetIds         []*Subnet          `mhcl:"ref=subnet_ids"` // azure??
	StartingNodeCount *int               `hcl:"starting_node_count,optional"`
	MaxNodeCount      int                `hcl:"max_node_count"`
	MinNodeCount      int                `hcl:"min_node_count"`
	Labels            map[string]string  `hcl:"labels,optional"`
	VmSize            string             `hcl:"vm_size"`
	DiskSizeGiB       int                `hcl:"disk_size_gib,optional"`
}

func (r *KubernetesServiceNodePool) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	if r.MinNodeCount < 1 {
		errs = append(errs, r.NewError("min_node_count", "node pool must have a min node count of at least 1"))
	}
	if r.MaxNodeCount < 1 {
		errs = append(errs, r.NewError("max_node_count", "node pool must have a max node count of at least 1"))
	}
	if r.MinNodeCount > r.MaxNodeCount {
		errs = append(errs, r.NewError("min_node_count", "min_node_count must be lower or equal to max_node_count"))
	}
	startingNodeCount := util.GetOrDefault(r.StartingNodeCount, r.MinNodeCount)
	if startingNodeCount < r.MinNodeCount || startingNodeCount > r.MaxNodeCount {
		errs = append(errs, r.NewError("starting_node_count", "starting_node_count must be between min and max node count"))
	}
	if err := common.CheckIfSizeIsValid(r.VmSize); err != nil {
		errs = append(errs, r.NewError("vm_size", err.Error()))
	}

	return errs
}

func (r *KubernetesServiceNodePool) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	if cloud == common.AWS {
		return output.GetResourceName(kubernetes_node_pool.AwsKubernetesNodeGroup{}), nil
	}
	if cloud == common.AZURE {
		return output.GetResourceName(kubernetes_node_pool.AzureKubernetesNodePool{}), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *KubernetesServiceNodePool) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds, err := util.MapSliceValuesErr(r.SubnetIds, func(v *Subnet) (string, error) {
		return resources.GetMainOutputId(v, cloud)
	})
	if err != nil {
		return nil, err
	}
	if cloud == common.AWS {
		roleName := fmt.Sprintf("iam_for_k8nodepool_%s", r.Name)
		role := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.GetTfResourceId(cloud), roleName),
			Name:             fmt.Sprintf("iam_for_k8nodepool_%s", r.Name),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
		}
		clusterId, err := resources.GetMainOutputId(r.ClusterId, cloud)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			&role,
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "AmazonEKSWorkerNodePolicy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.GetTfResourceId(cloud)),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "AmazonEKS_CNI_Policy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.GetTfResourceId(cloud)),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "AmazonEC2ContainerRegistryReadOnly")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.GetTfResourceId(cloud)),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			},
			&kubernetes_node_pool.AwsKubernetesNodeGroup{
				AwsResource:   common.NewAwsResourceWithIdOnly(r.GetTfResourceId(cloud)),
				ClusterName:   clusterId,
				NodeGroupName: r.Name,
				NodeRoleArn:   fmt.Sprintf("aws_iam_role.%s.arn", r.GetTfResourceId(cloud)),
				SubnetIds:     subnetIds,
				ScalingConfig: kubernetes_node_pool.ScalingConfig{
					DesiredSize: util.GetOrDefault(r.StartingNodeCount, r.MinNodeCount),
					MaxSize:     r.MaxNodeCount,
					MinSize:     r.MinNodeCount,
				},
				Labels:        r.Labels,
				InstanceTypes: []string{common.VMSIZE[r.VmSize][cloud]},
			},
		}, nil
	} else if cloud == common.AZURE {
		if r.IsDefaultPool {
			// this will be embedded in the cluster instead
			return nil, nil
		}

		pool, err := r.translateAzNodePool(ctx)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			pool,
		}, nil

	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *KubernetesServiceNodePool) translateAzNodePool(ctx resources.MultyContext) (*kubernetes_node_pool.AzureKubernetesNodePool, error) {
	clusterId, err := resources.GetMainOutputId(r.ClusterId, common.AZURE)
	if err != nil {
		return nil, err
	}
	return &kubernetes_node_pool.AzureKubernetesNodePool{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(common.AZURE)},
			Name:              r.Name,
			ResourceGroupName: rg.GetResourceGroupName(r.ResourceGroupId, common.AZURE),
		},
		ClusterId:         clusterId,
		NodeCount:         util.GetOrDefault(r.StartingNodeCount, r.MinNodeCount),
		MaxSize:           r.MaxNodeCount,
		MinSize:           r.MinNodeCount,
		Labels:            r.Labels,
		EnableAutoScaling: true,
		VmSize:            common.VMSIZE[r.VmSize][common.AZURE],
	}, nil
}

func (r *KubernetesServiceNodePool) GetLocation(cloud common.CloudProvider, ctx resources.MultyContext) string {
	return r.ClusterId.GetLocation(cloud, ctx)
}
