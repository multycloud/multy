package kubernetes_node_pool

import "github.com/multycloud/multy/resources/common"

type AzureKubernetesNodePool struct {
	*common.AzResource     `hcl:",squash" default:"name=azurerm_kubernetes_cluster_node_pool"`
	ClusterId              string            `hcl:"kubernetes_cluster_id,expr"  hcle:"omitempty" json:"kubernetes_cluster_id"`
	Name                   string            `hcl:"name,optional"  hcle:"omitempty" json:"name"`
	NodeCount              int               `hcl:"node_count" json:"node_count"`
	MaxSize                int               `hcl:"max_count" json:"max_count"`
	MinSize                int               `hcl:"min_count" json:"min_count"`
	Labels                 map[string]string `hcl:"node_labels" hcle:"omitempty" json:"node_labels"`
	EnableAutoScaling      bool              `hcl:"enable_auto_scaling" json:"enable_auto_scaling"`
	VmSize                 string            `hcl:"vm_size" json:"vm_size"`
	VirtualNetworkSubnetId string            `hcl:"vnet_subnet_id,expr" json:"vnet_subnet_id"`
	Zones                  []string          `hcl:"zones" hcle:"omitempty" json:"zones"`
	OsDiskSizeGb           int               `hcl:"os_disk_size_gb" hcle:"omitempty" json:"os_disk_size_gb" `
}
