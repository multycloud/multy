package route_table

import (
	"github.com/multycloud/multy/resources/common"
)

type GoogleComputeRoute struct {
	*common.GcpResource `hcl:",squash" default:"name=google_compute_route"`
	DestRange           string   `hcl:"dest_range" json:"dest_range" hcle:"omitempty"`
	Network             string   `hcl:"network,expr"`
	Priority            int      `hcl:"priority"`
	Tags                []string `hcl:"tags"`
	NextHopGateway      string   `hcl:"next_hop_gateway" json:"next_hop_gateway"`
}
