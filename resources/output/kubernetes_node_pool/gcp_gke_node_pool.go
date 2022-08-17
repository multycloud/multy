package kubernetes_node_pool

import "github.com/multycloud/multy/resources/common"

type GoogleContainerNodePool struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_container_node_pool"`
	Cluster             string                               `hcl:"cluster,expr" json:"cluster"` //expr
	InitialNodeCount    int                                  `hcl:"initial_node_count" json:"initial_node_count"`
	NodeLocations       []string                             `hcl:"node_locations" hcle:"omitempty" json:"node_locations"`
	Autoscaling         []GoogleContainerNodePoolAutoScaling `hcl:"autoscaling,blocks" json:"autoscaling"`
	NodeConfig          []GoogleContainerNodeConfig          `hcl:"node_config,blocks" json:"node_config"`
	NetworkConfig       []GoogleContainerNetworkConfig       `hcl:"network_config,blocks" hcle:"omitempty" json:"network_config"`
}

type GoogleContainerNodePoolAutoScaling struct {
	MinNodeCount int `hcl:"min_node_count" json:"min_node_count"`
	MaxNodeCount int `hcl:"max_node_count" json:"max_node_count"`
}

type GoogleContainerNodeConfig struct {
	DiskSizeGb  int               `hcl:"disk_size_gb" hcle:"omitempty" json:"disk_size_gb"`
	DiskType    string            `hcl:"disk_type" hcle:"omitempty" json:"disk_type"`
	ImageType   string            `hcl:"image_type" hcle:"omitempty" json:"image_type"`
	Labels      map[string]string `hcl:"labels" hcle:"omitempty" json:"labels"`
	MachineType string            `hcl:"machine_type" json:"machine_type"`
	Metadata    map[string]string `hcl:"metadata" hcle:"omitempty" json:"metadata"`
	Tags        []string          `hcl:"tags" hcle:"omitempty" json:"tags"`

	ServiceAccount string   `hcl:"service_account,expr" json:"service_account"`
	OAuthScopes    []string `hcl:"oauth_scopes" json:"o_auth_scopes"`
}

type GoogleContainerNetworkConfig struct {
	CreatePodRange   bool   `hcl:"create_pod_range" json:"create_pod_range"`
	PodIpv4CidrBlock string `hcl:"pod_ipv4_cidr_block" json:"pod_ipv_4_cidr_block"`
}
