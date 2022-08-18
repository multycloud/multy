package kubernetes_service

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
)

type GoogleContainerCluster struct {
	*common.GcpResource   `hcl:",squash"  default:"name=google_container_cluster"`
	RemoveDefaultNodePool bool                                             `hcl:"remove_default_node_pool" json:"remove_default_node_pool"`
	InitialNodeCount      int                                              `hcl:"initial_node_count" json:"initial_node_count"`
	Subnetwork            string                                           `hcl:"subnetwork,expr" json:"subnetwork"`
	Network               string                                           `hcl:"network,expr" json:"network"`
	IpAllocationPolicy    []GoogleContainerClusterIpAllocationPolicy       `hcl:"ip_allocation_policy,blocks" json:"ip_allocation_policy"`
	Location              string                                           `hcl:"location" json:"location"`
	NodeConfig            []kubernetes_node_pool.GoogleContainerNodeConfig `hcl:"node_config,blocks" json:"node_config"`

	// outputs
	Endpoint   string                       `json:"endpoint" hcle:"omitempty" json:"endpoint"`
	MasterAuth []GoogleContainerClusterAuth `json:"master_auth" hcle:"omitempty" json:"master_auth"`
}

type GoogleContainerClusterAuth struct {
	ClusterCaCertificate string `json:"cluster_ca_certificate"`
	ClientCertificate    string `json:"client_certificate"`
	ClientKey            string `json:"client_key"`
}

type GoogleContainerClusterIpAllocationPolicy struct {
	//ClusterIpv4CidrBlock  string `hcl:"cluster_ipv_4_cidr_block"`
	ServicesIpv4CidrBlock string `hcl:"services_ipv4_cidr_block" json:"services_ipv4_cidr_block"`
}
