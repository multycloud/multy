package kubernetes_service

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
)

type GoogleContainerCluster struct {
	*common.GcpResource   `hcl:",squash"  default:"name=google_container_cluster"`
	RemoveDefaultNodePool bool                                           `hcl:"remove_default_node_pool"`
	InitialNodeCount      int                                            `hcl:"initial_node_count"`
	Subnetwork            string                                         `hcl:"subnetwork,expr"`
	Network               string                                         `hcl:"network,expr"`
	IpAllocationPolicy    GoogleContainerClusterIpAllocationPolicy       `hcl:"ip_allocation_policy"`
	Location              string                                         `hcl:"location"`
	NodeConfig            kubernetes_node_pool.GoogleContainerNodeConfig `hcl:"node_config"`

	// outputs
	Endpoint   string                       `json:"endpoint" hcle:"omitempty"`
	MasterAuth []GoogleContainerClusterAuth `json:"master_auth" hcle:"omitempty"`
}

type GoogleContainerClusterAuth struct {
	ClusterCaCertificate string `json:"cluster_ca_certificate"`
	ClientCertificate    string `json:"client_certificate"`
	ClientKey            string `json:"client_key"`
}

type GoogleContainerClusterIpAllocationPolicy struct {
	//ClusterIpv4CidrBlock  string `hcl:"cluster_ipv_4_cidr_block"`
	ServicesIpv4CidrBlock string `hcl:"services_ipv4_cidr_block"`
}
