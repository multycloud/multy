package network_security_group

import (
	"github.com/multycloud/multy/resources/common"
	"golang.org/x/exp/slices"
)

type GoogleComputeFirewall struct {
	*common.GcpResource `hcl:",squash" default:"name=google_compute_firewall"`

	Description string `hcl:"description" hcle:"omitempty" json:"description"`
	Network     string `hcl:"network,expr" json:"network"`

	Direction         string   `hcl:"direction" json:"direction"`
	SourceRanges      []string `hcl:"source_ranges" hcle:"omitempty" json:"source_ranges"`
	DestinationRanges []string `hcl:"destination_ranges" hcle:"omitempty" json:"destination_ranges"`
	Priority          int      `hcl:"priority" json:"priority"`

	AllowRules []GoogleComputeFirewallRule `hcl:"allow,blocks" json:"allow"`
	DenyRules  []GoogleComputeFirewallRule `hcl:"deny,blocks" json:"deny"`

	TargetTags []string `hcl:"target_tags" hcle:"omitempty" json:"target_tags"`
}

type GoogleComputeFirewallRule struct {
	Protocol string   `hcl:"protocol" json:"protocol"`
	Ports    []string `hcl:"ports" hcle:"omitempty" json:"ports"`
}

func (f *GoogleComputeFirewall) EqualsExceptDirection(other *GoogleComputeFirewall) bool {
	return f.Direction != other.Direction && f.Priority == other.Priority &&
		f.GetCidrBlock() == other.GetCidrBlock() && slices.EqualFunc(f.AllowRules, other.AllowRules,
		func(rule1, rule2 GoogleComputeFirewallRule) bool { return rule1.Equal(rule2) }) &&
		slices.Equal(f.TargetTags, other.TargetTags)
}

func (r GoogleComputeFirewallRule) Equal(other GoogleComputeFirewallRule) bool {
	return r.Protocol == other.Protocol && slices.Compare(r.Ports, other.Ports) == 0
}

func (f *GoogleComputeFirewall) GetCidrBlock() string {
	if f.Direction == "INGRESS" {
		if len(f.SourceRanges) == 0 {
			return ""
		}
		return f.SourceRanges[0]
	} else {
		if len(f.DestinationRanges) == 0 {
			return ""
		}
		return f.DestinationRanges[0]
	}
}
