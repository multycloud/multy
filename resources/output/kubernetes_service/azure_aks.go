package kubernetes_service

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
)

type AzureIdentity struct {
	PrincipalId string `hcl:"principal_id,expr"  hcle:"omitempty"`
	TenantId    string `hcl:"tenant_id,expr"  hcle:"omitempty"`
	Type        string `hcl:"type" hcle:"omitempty"`
}

type AzureEksCluster struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_kubernetes_cluster"`
	DefaultNodePool    *kubernetes_node_pool.AzureKubernetesNodePool `hcl:"default_node_pool" json:"-"`
	DnsPrefix          string                                        `hcl:"dns_prefix" json:"dns_prefix"`
	Identity           []AzureIdentity                               `hcl:"identity,blocks" json:"identity"`
	NetworkProfile     []NetworkProfile                              `hcl:"network_profile,blocks" json:"network_profile"`

	// outputs
	KubeConfigRaw      string                                          `json:"kube_config_raw" hcle:"omitempty" json:"kube_config_raw"`
	KubeConfig         []AzureKubeConfig                               `json:"kube_config" hcle:"omitempty" json:"kube_config"`
	DefaultNodePoolOut []*kubernetes_node_pool.AzureKubernetesNodePool `json:"default_node_pool"`
}

type NetworkProfile struct {
	NetworkPlugin    string `hcl:"network_plugin" json:"network_plugin"`
	DnsServiceIp     string `hcl:"dns_service_ip" json:"dns_service_ip"`
	DockerBridgeCidr string `hcl:"docker_bridge_cidr" json:"docker_bridge_cidr"`
	ServiceCidr      string `hcl:"service_cidr" json:"service_cidr"`
}

type AzureKubeConfig struct {
	Host                 string `json:"host" hcle:"omitempty"`
	ClusterCaCertificate string `json:"cluster_ca_certificate" hcle:"omitempty"`
}

type AzureUserAssignedIdentity struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_user_assigned_identity"`
}

func (r AzureUserAssignedIdentity) GetIdentity() AzureIdentity {
	return AzureIdentity{
		PrincipalId: fmt.Sprintf("azurerm_user_assigned_identity.%s.principal_id", r.ResourceId),
		TenantId:    fmt.Sprintf("azurerm_user_assigned_identity.%s.tenant_id", r.ResourceId),
		Type:        "UserAssigned",
	}
}
