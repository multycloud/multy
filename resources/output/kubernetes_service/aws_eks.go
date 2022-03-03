package kubernetes_service

import (
	"multy-go/resources/common"
)

type AwsEksCluster struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_eks_cluster"`
	RoleArn             string    `hcl:"role_arn,expr"`
	VpcConfig           VpcConfig `hcl:"vpc_config"`
	Name                string    `hcl:"name"`
}

type VpcConfig struct {
	SubnetIds []string `hcl:"subnet_ids"`
}
