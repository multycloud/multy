package kuerbenetes_node_pool

import (
	"multy-go/resources/common"
)

type AwsEksCluster struct {
	common.AwsResource `hcl:",squash" default:"name=aws_eks_cluster"`
	RoleArn            string    `hcl:"role_arn"`
	VpcConfig          VpcConfig `hcl:"vpc_config"`
}

type VpcConfig struct {
	SubnetIds []string `hcl:"subnet_ids"`
}
