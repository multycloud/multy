package virtual_network

import "github.com/multycloud/multy/resources/common"

const GcpResourceName = "google_compute_network"

type GoogleComputeNetwork struct {
	*common.GcpResource         `hcl:",squash"  default:"name=google_compute_network"`
	RoutingMode                 string `hcl:"routing_mode"` // REGIONAL
	Description                 string `hcl:"description"`
	AutoCreateSubnetworks       bool   `hcl:"auto_create_subnetworks"`
	DeleteDefaultRoutesOnCreate bool   `hcl:"delete_default_routes_on_create"`
}
