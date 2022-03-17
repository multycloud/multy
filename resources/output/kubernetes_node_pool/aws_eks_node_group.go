package kubernetes_node_pool

import "github.com/multycloud/multy/resources/common"

type ScalingConfig struct {
	DesiredSize int `hcl:"desired_size"`
	MaxSize     int `hcl:"max_size"`
	MinSize     int `hcl:"min_size"`
}

type UpdateConfig struct {
	MaxUnavailable int `hcl:"max_unavailable"`
}

type AwsKubernetesNodeGroup struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_eks_node_group"`
	ClusterName         string            `hcl:"cluster_name"`
	NodeGroupName       string            `hcl:"node_group_name"`
	NodeRoleArn         string            `hcl:"node_role_arn,expr"`
	SubnetIds           []string          `hcl:"subnet_ids"`
	ScalingConfig       ScalingConfig     `hcl:"scaling_config"`
	UpdateConfig        UpdateConfig      `hcl:"update_config" hcle:"omitempty"`
	Labels              map[string]string `hcl:"labels" hcle:"omitempty"`
	InstanceTypes       []string          `hcl:"instance_types"`
}
