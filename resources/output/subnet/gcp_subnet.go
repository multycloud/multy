package subnet

import "github.com/multycloud/multy/resources/common"

type GoogleComputeSubnetwork struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_compute_subnetwork"`
	IpCidrRange         string `hcl:"ip_cidr_range" json:"ip_cidr_range"`
	Network             string `hcl:"network,expr" json:"network"`
	Description         string `hcl:"description" hcle:"omitempty" json:"description"`
}
