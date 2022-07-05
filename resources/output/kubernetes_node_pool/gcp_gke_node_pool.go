package kubernetes_node_pool

import "github.com/multycloud/multy/resources/common"

type GoogleContainerNodePool struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_container_node_pool"`
	Cluster             string                             `hcl:"cluster,expr"` //expr
	InitialNodeCount    int                                `hcl:"initial_node_count"`
	NodeLocations       []string                           `hcl:"node_locations" hcle:"omitempty"`
	Autoscaling         GoogleContainerNodePoolAutoScaling `hcl:"autoscaling"`
	NodeConfig          GoogleContainerNodeConfig          `hcl:"node_config"`
	NetworkConfig       GoogleContainerNetworkConfig       `hcl:"network_config" hcle:"omitempty"`
}

type GoogleContainerNodePoolAutoScaling struct {
	MinNodeCount int `hcl:"min_node_count"`
	MaxNodeCount int `hcl:"max_node_count"`
}

type GoogleContainerNodeConfig struct {
	DiskSizeGb  int               `hcl:"disk_size_gb" hcle:"omitempty"`
	DiskType    string            `hcl:"disk_type" hcle:"omitempty"`
	ImageType   string            `hcl:"image_type" hcle:"omitempty"`
	Labels      map[string]string `hcl:"labels" hcle:"omitempty"`
	MachineType string            `hcl:"machine_type"`
	Metadata    map[string]string `hcl:"metadata" hcle:"omitempty"`
	Tags        []string          `hcl:"tags" hcle:"omitempty"`

	ServiceAccount string   `hcl:"service_account,expr"`
	OAuthScopes    []string `hcl:"oauth_scopes"`
}

type GoogleContainerNetworkConfig struct {
	CreatePodRange   bool   `hcl:"create_pod_range"`
	PodIpv4CidrBlock string `hcl:"pod_ipv4_cidr_block"`
}
