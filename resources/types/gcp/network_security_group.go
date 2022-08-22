package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

type GcpNetworkSecurityGroup struct {
	*types.NetworkSecurityGroup
}

func InitNetworkSecurityGroup(r *types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource] {
	return GcpNetworkSecurityGroup{r}
}

func (r GcpNetworkSecurityGroup) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.NetworkSecurityGroupResource, error) {
	out := &resourcespb.NetworkSecurityGroupResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:             r.Args.Name,
		VirtualNetworkId: r.Args.VirtualNetworkId,
		Rules:            r.Args.Rules,
		GcpOverride:      r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	out.GcpOutputs = &resourcespb.NetworkSecurityGroupGcpOutputs{}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.GoogleComputeFirewall](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.GcpOutputs.ComputeFirewallId = append(out.GcpOutputs.ComputeFirewallId, stateResource.SelfLink)
	} else {
		statuses["gcp_default_firewall_rule"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	out.Rules = nil
	for i, rule := range r.NetworkSecurityGroup.Args.Rules {
		var firewalls []*network_security_group.GoogleComputeFirewall
		for _, firewallId := range r.getRuleIds(i, rule) {
			if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.GoogleComputeFirewall](state, firewallId); exists {
				if err != nil {
					return nil, err
				}
				firewalls = append(firewalls, stateResource)
				out.GcpOutputs.ComputeFirewallId = append(out.GcpOutputs.ComputeFirewallId, stateResource.SelfLink)
			}
		}

		if len(firewalls) == 0 {
			statuses[fmt.Sprintf("gcp_firewall_rule_%d", i)] = commonpb.ResourceStatus_NEEDS_CREATE
			continue
		}

		firewall := firewalls[0]
		direction := resourcespb.Direction_BOTH_DIRECTIONS
		if len(firewalls) == 2 {
			if !firewall.EqualsExceptDirection(firewalls[1]) {
				statuses[fmt.Sprintf("gcp_firewall_rule_%d", i)] = commonpb.ResourceStatus_NEEDS_UPDATE
			}
		} else {
			direction = resourcespb.Direction(resourcespb.Direction_value[firewall.Direction])
		}

		if len(firewall.AllowRules) != 1 || len(firewall.DenyRules) != 0 {
			statuses[fmt.Sprintf("gcp_firewall_rule_%d", i)] = commonpb.ResourceStatus_NEEDS_UPDATE
			continue
		}

		portRange, err := translateToPortRange(firewall.AllowRules[0].Ports)
		if err != nil {
			statuses[fmt.Sprintf("gcp_firewall_rule_%d", i)] = commonpb.ResourceStatus_NEEDS_UPDATE
			continue
		}

		if !slices.Equal(firewall.TargetTags, r.getNsgTag()) {
			statuses[fmt.Sprintf("gcp_firewall_rule_%d", i)] = commonpb.ResourceStatus_NEEDS_UPDATE
		}

		out.Rules = append(out.Rules, &resourcespb.NetworkSecurityRule{
			Protocol:  firewall.AllowRules[0].Protocol,
			Priority:  int64(firewall.Priority),
			PortRange: portRange,
			CidrBlock: firewall.GetCidrBlock(),
			Direction: direction,
		})
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r GcpNetworkSecurityGroup) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	return r.getFirewalls(), nil
}

func (r GcpNetworkSecurityGroup) getFirewalls() []output.TfBlock {
	var firewalls []output.TfBlock
	// gcp allows egress traffic by default: https://cloud.google.com/vpc/docs/firewalls
	// so we add default deny rule to egress
	ruleName := r.defaultRuleName()
	firewalls = append(firewalls, &network_security_group.GoogleComputeFirewall{
		GcpResource:       common.NewGcpResource(r.ResourceId, ruleName, r.Args.GetGcpOverride().GetProject()),
		Network:           fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
		DestinationRanges: []string{"0.0.0.0/0"},
		Direction:         resourcespb.Direction_EGRESS.String(),
		DenyRules:         []network_security_group.GoogleComputeFirewallRule{{Protocol: "all"}},
		TargetTags:        r.getNsgTag(),
		Priority:          65535,
	})
	for i, rule := range r.NetworkSecurityGroup.Args.Rules {
		if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
			firewalls = append(firewalls, r.buildRule(i, rule, resourcespb.Direction_INGRESS))
			firewalls = append(firewalls, r.buildRule(i, rule, resourcespb.Direction_EGRESS))
		} else {
			firewalls = append(firewalls, r.buildRule(i, rule, rule.Direction))
		}
	}
	return firewalls
}

func (r GcpNetworkSecurityGroup) defaultRuleName() string {
	return fmt.Sprintf("%s-%s", r.Args.Name, "default-deny-egress")
}

func (r GcpNetworkSecurityGroup) getRuleIds(i int, rule *resourcespb.NetworkSecurityRule) []string {
	var directions []resourcespb.Direction
	if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
		directions = append(directions, resourcespb.Direction_INGRESS, resourcespb.Direction_EGRESS)
	} else {
		directions = append(directions, rule.Direction)
	}

	var result []string
	for _, direction := range directions {
		suffix := getRuleSuffix(i, direction)
		ruleId := fmt.Sprintf("%s-%s", r.ResourceId, suffix)
		result = append(result, ruleId)
	}
	return result
}

func getRuleSuffix(i int, direction resourcespb.Direction) string {
	aliases := map[resourcespb.Direction]string{resourcespb.Direction_INGRESS: "i", resourcespb.Direction_EGRESS: "e"}
	suffix := fmt.Sprintf("%s-%d", aliases[direction], i)
	return suffix
}

func (r GcpNetworkSecurityGroup) buildRule(i int, rule *resourcespb.NetworkSecurityRule, direction resourcespb.Direction) *network_security_group.GoogleComputeFirewall {
	suffix := getRuleSuffix(i, direction)
	ruleId := fmt.Sprintf("%s-%s", r.ResourceId, suffix)
	ruleName := fmt.Sprintf("%s-%s", r.Args.Name, suffix)
	firewall := &network_security_group.GoogleComputeFirewall{
		GcpResource: common.NewGcpResource(ruleId, ruleName, r.Args.GetGcpOverride().GetProject()),
		Direction:   direction.String(),
		Network:     fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
		AllowRules: []network_security_group.GoogleComputeFirewallRule{
			{
				Protocol: rule.Protocol,
				// TODO: group similar rules
				Ports: translatePortRange(rule.PortRange),
			},
		},
		Priority:   int(rule.Priority),
		TargetTags: r.getNsgTag(),
	}
	if direction == resourcespb.Direction_INGRESS {
		firewall.SourceRanges = []string{rule.CidrBlock}
	} else if direction == resourcespb.Direction_EGRESS {
		firewall.DestinationRanges = []string{rule.CidrBlock}
	}

	return firewall
}

func (r GcpNetworkSecurityGroup) getNsgTag() []string {
	return []string{fmt.Sprintf("nsg-%s", r.Args.Name)}
}

func translatePortRange(ports *resourcespb.PortRange) []string {
	from := ports.From
	to := ports.To
	if ports.To == 0 {
		to = 65535
	}
	if from == to {
		return []string{strconv.Itoa(int(from))}
	}

	return []string{fmt.Sprintf("%d-%d", from, to)}
}

func translateToPortRange(ranges []string) (*resourcespb.PortRange, error) {
	if len(ranges) != 1 {
		return nil, fmt.Errorf("invalid multy port range")
	}
	r := strings.SplitN(ranges[0], "-", 2)
	if len(r) == 0 {
		return nil, fmt.Errorf("empty port range")
	}
	fromPort, err := strconv.Atoi(r[0])
	if err != nil {
		return nil, err
	}
	if len(r) == 1 {
		return &resourcespb.PortRange{
			From: int32(fromPort),
			To:   int32(fromPort),
		}, nil
	}

	toPort, err := strconv.Atoi(r[1])
	if err != nil {
		return nil, err
	}
	if toPort == 65535 {
		toPort = 0
	}

	return &resourcespb.PortRange{
		From: int32(fromPort),
		To:   int32(toPort),
	}, nil
}

func (r GcpNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main resource for gcp firewalls")
}
