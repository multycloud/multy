package kubernetes_service

import (
	"github.com/multycloud/multy/resources/common"
)

type AwsEksCluster struct {
	*common.AwsResource     `hcl:",squash" default:"name=aws_eks_cluster"`
	RoleArn                 string                  `hcl:"role_arn,expr"`
	VpcConfig               VpcConfig               `hcl:"vpc_config"`
	KubernetesNetworkConfig KubernetesNetworkConfig `hcl:"kubernetes_network_config"`
	Name                    string                  `hcl:"name"`

	// outputs
	Endpoint             string                 `json:"endpoint" hcle:"omitempty"`
	CertificateAuthority []CertificateAuthority `json:"certificate_authority"  hcle:"omitempty"`
	Arn                  string                 `json:"arn" hcle:"omitempty"`
}

type VpcConfig struct {
	SubnetIds             []string `hcl:"subnet_ids,expr"`
	EndpointPrivateAccess bool     `hcl:"endpoint_private_access"`
}

type KubernetesNetworkConfig struct {
	ServiceIpv4Cidr string `hcl:"service_ipv4_cidr"`
}

type CertificateAuthority struct {
	Data string `json:"data"  hcle:"omitempty"`
}
