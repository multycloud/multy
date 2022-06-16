package network_security_group

import "github.com/multycloud/multy/resources/common"

type GoogleComputeFirewall struct {
	*common.GcpResource `hcl:",squash" default:"name=google_compute_firewall"`

	Description string `hcl:"description" hcle:"omitempty"`
	Network     string `hcl:"network,expr"`

	Direction         string   `hcl:"direction"`
	SourceRanges      []string `hcl:"source_ranges" hcle:"omitempty"`
	DestinationRanges []string `hcl:"destination_ranges" hcle:"omitempty"`
	Priority          int      `hcl:"priority"`

	AllowRules []GoogleComputeFirewallRule `hcl:"allow,blocks"`
	DenyRules  []GoogleComputeFirewallRule `hcl:"deny,blocks"`

	TargetTags []string `hcl:"target_tags" hcle:"omitempty"`
}

type GoogleComputeFirewallRule struct {
	Protocol string   `hcl:"protocol"`
	Ports    []string `hcl:"ports" hcle:"omitempty"`
}
