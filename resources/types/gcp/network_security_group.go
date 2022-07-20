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
	"strconv"
)

type GcpNetworkSecurityGroup struct {
	*types.NetworkSecurityGroup
}

func InitNetworkSecurityGroup(r *types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource] {
	return GcpNetworkSecurityGroup{r}
}

func (r GcpNetworkSecurityGroup) FromState(state *output.TfState) (*resourcespb.NetworkSecurityGroupResource, error) {
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

	for _, firewall := range r.getFirewalls() {
		stateResource, err := output.GetParsedById[network_security_group.GoogleComputeFirewall](state, firewall.GetResourceId())
		if err != nil {
			return nil, err
		}
		out.GcpOutputs.ComputeFirewallId = append(out.GcpOutputs.ComputeFirewallId, stateResource.SelfLink)
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
	ruleName := fmt.Sprintf("%s-%s", r.Args.Name, "default-deny-egress")
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

func (r GcpNetworkSecurityGroup) buildRule(i int, rule *resourcespb.NetworkSecurityRule, direction resourcespb.Direction) *network_security_group.GoogleComputeFirewall {
	aliases := map[resourcespb.Direction]string{resourcespb.Direction_INGRESS: "i", resourcespb.Direction_EGRESS: "e"}
	suffix := fmt.Sprintf("%s-%d", aliases[direction], i)
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

func (r GcpNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("no main resource for gcp firewalls")
}
