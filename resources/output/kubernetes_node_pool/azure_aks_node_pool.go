package kubernetes_node_pool

import "multy/resources/common"

type AzureKubernetesNodePool struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_kubernetes_cluster_node_pool"`
	ClusterId          string            `hcl:"kubernetes_cluster_id"  hcle:"omitempty"`
	Name               string            `hcl:"name,optional"  hcle:"omitempty"`
	NodeCount          int               `hcl:"node_count"`
	MaxSize            int               `hcl:"max_count"`
	MinSize            int               `hcl:"min_count"`
	Labels             map[string]string `hcl:"node_labels" hcle:"omitempty"`
	EnableAutoScaling  bool              `hcl:"enable_auto_scaling"`
	VmSize             string            `hcl:"vm_size"`
}
