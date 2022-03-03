package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/iam"
	"multy-go/resources/output/kuerbenetes_node_pool"
	"multy-go/validate"
)

type KubernetesServiceNodePool struct {
	*resources.CommonResourceParams
	Name        string `hcl:"name"`
	ClusterName string `hcl:"cluster_name"`
	// TODO: somehow set map_public_ip_on_launch on public subnet for aws subbnets
	SubnetIds []string `hcl:"subnet_ids"` // azure??
	// TODO: check if this is set
	StartingNodeCount int               `hcl:"starting_node_count,optional"`
	MaxNodeCount      int               `hcl:"max_node_count"`
	MinNodeCount      int               `hcl:"min_node_count"`
	Tags              map[string]string `hcl:"tags,optional"`
	VmSize            string            `hcl:"vm_size,optional"`
	DiskSizeGiB       int               `hcl:"disk_size_gib,optional"`
}

func (r *KubernetesServiceNodePool) Validate(ctx resources.MultyContext) []validate.ValidationError {
	return nil
}

func (r *KubernetesServiceNodePool) GetMainResourceName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return common.GetResourceName(kuerbenetes_node_pool.AwsKubernetesNodeGroup{})
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *KubernetesServiceNodePool) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {

	roleName := fmt.Sprintf("iam_for_k8nodepool_%s", r.Name)
	role := iam.AwsIamRole{
		AwsResource:      common.NewAwsResource(r.GetTfResourceId(cloud), roleName),
		Name:             fmt.Sprintf("iam_for_k8nodepool_%s", r.Name),
		AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
	}

	if cloud == common.AWS {
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
			&kuerbenetes_node_pool.AwsKubernetesNodeGroup{
				AwsResource:   common.NewAwsResourceWithIdOnly(r.ResourceId),
				ClusterName:   r.ClusterName,
				NodeGroupName: r.Name,
				NodeRoleArn:   fmt.Sprintf("aws_iam_role.%s.arn", r.GetTfResourceId(cloud)),
				SubnetIds:     r.SubnetIds,
				ScalingConfig: kuerbenetes_node_pool.ScalingConfig{
					DesiredSize: r.StartingNodeCount,
					MaxSize:     r.MaxNodeCount,
					MinSize:     r.MinNodeCount,
				},
				Tags: r.Tags,
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}
