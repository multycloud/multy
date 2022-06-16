package public_ip

import "github.com/multycloud/multy/resources/common"

type GoogleComputeAddress struct {
	*common.GcpResource `hcl:",squash" default:"name=google_compute_address"`
	NetworkTier         string `hcl:"network_tier"`

	// outputs
	Address string `json:"address"  hcle:"omitempty"`
}
