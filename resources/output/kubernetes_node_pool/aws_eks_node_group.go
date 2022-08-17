package kubernetes_node_pool

import "github.com/multycloud/multy/resources/common"

type ScalingConfig struct {
	DesiredSize int `hcl:"desired_size" json:"desired_size"`
	MaxSize     int `hcl:"max_size" json:"max_size"`
	MinSize     int `hcl:"min_size" json:"min_size"`
}

type UpdateConfig struct {
	MaxUnavailable int `hcl:"max_unavailable"`
}

type AwsKubernetesNodeGroup struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_eks_node_group"`
	ClusterName         string            `hcl:"cluster_name,expr" json:"cluster_name"`
	NodeGroupName       string            `hcl:"node_group_name" json:"node_group_name"`
	NodeRoleArn         string            `hcl:"node_role_arn,expr" json:"node_role_arn"`
	SubnetIds           []string          `hcl:"subnet_ids,expr" json:"subnet_ids"`
	ScalingConfig       []ScalingConfig   `hcl:"scaling_config,blocks" json:"scaling_config"`
	UpdateConfig        []UpdateConfig    `hcl:"update_config,blocks" hcle:"omitempty" json:"update_config"`
	Labels              map[string]string `hcl:"labels" hcle:"omitempty" json:"labels"`
	InstanceTypes       []string          `hcl:"instance_types" json:"instance_types"`
	DiskSize            int               `hcl:"disk_size" hcle:"omitempty"  json:"disk_size"`

	Arn string `json:"arn" hcle:"omitempty" json:"arn"`
}
