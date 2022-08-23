package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/protobuf/proto"
)

type AwsNetworkSecurityGroup struct {
	*types.NetworkSecurityGroup
}

func InitNetworkSecurityGroup(r *types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource] {
	return AwsNetworkSecurityGroup{r}
}

const (
	INGRESS = "ingress"
	EGRESS  = "egress"
	BOTH    = "both"
	ALLOW   = "allow"
	DENY    = "deny"
)

func (r AwsNetworkSecurityGroup) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.NetworkSecurityGroupResource, error) {
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

	statuses := map[string]commonpb.ResourceStatus_Status{}
	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.AwsSecurityGroup](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.Name = stateResource.AwsResource.Tags["Name"]
		out.AwsOutputs = &resourcespb.NetworkSecurityGroupAwsOutputs{SecurityGroupId: stateResource.ResourceId}
		out.Rules = nil
		output.AddToStatuses(statuses, "aws_network_security_group", output.MaybeGetPlannedChageById[network_security_group.AwsSecurityGroup](plan, r.ResourceId))

		// first attempt to get user-provided rules so that the order remains the same even if new rules were added
		// this is also needed because certain properties (priority for example), are lost after deployment
		for _, rule := range r.Args.Rules {
			awsRule := translateAwsRule(rule)
			if rule.Direction == resourcespb.Direction_INGRESS {
				if i := checkIfExists(awsRule, stateResource.Ingress); i >= 0 {
					stateResource.Ingress = append(stateResource.Ingress[:i], stateResource.Ingress[i+1:]...)
					out.Rules = append(out.Rules, rule)
				}
			} else if rule.Direction == resourcespb.Direction_EGRESS {
				if i := checkIfExists(awsRule, stateResource.Egress); i >= 0 {
					stateResource.Egress = append(stateResource.Egress[:i], stateResource.Egress[i+1:]...)
					out.Rules = append(out.Rules, rule)
				}
			} else if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
				i := checkIfExists(awsRule, stateResource.Ingress)
				j := checkIfExists(awsRule, stateResource.Egress)
				if i < 0 && j < 0 {
					continue
				}

				if i >= 0 {
					stateResource.Ingress = append(stateResource.Ingress[:i], stateResource.Ingress[i+1:]...)
				}
				if j >= 0 {
					stateResource.Egress = append(stateResource.Egress[:j], stateResource.Egress[j+1:]...)
				}
				clonedRule := proto.Clone(rule).(*resourcespb.NetworkSecurityRule)
				if i < 0 {
					clonedRule.Direction = resourcespb.Direction_EGRESS
				} else if j < 0 {
					clonedRule.Direction = resourcespb.Direction_EGRESS
				} else {
					clonedRule.Direction = resourcespb.Direction_BOTH_DIRECTIONS
				}

				out.Rules = append(out.Rules, clonedRule)
			}
		}

		allowVpcTrafficRule := r.getAllowVpcTrafficRule()
		for _, rule := range stateResource.Ingress {
			if rule.Equals(allowVpcTrafficRule) {
				continue
			}
			if len(rule.CidrBlocks) == 0 {
				statuses["azure_network_security_group"] = commonpb.ResourceStatus_NEEDS_UPDATE
				continue
			}
			out.Rules = append(out.Rules, translateFromAwsRule(rule))
		}
		for _, rule := range stateResource.Egress {
			if rule.Equals(allowVpcTrafficRule) {
				continue
			}
			if len(rule.CidrBlocks) == 0 {
				statuses["azure_network_security_group"] = commonpb.ResourceStatus_NEEDS_UPDATE
				continue
			}
			out.Rules = append(out.Rules, translateFromAwsRule(rule))
		}

	} else {
		statuses["aws_network_security_group"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func checkIfExists(rule network_security_group.AwsSecurityGroupRule, rules []network_security_group.AwsSecurityGroupRule) int {
	for i, r := range rules {
		if r.Equals(rule) {
			return i
		}
	}
	return -1
}

func translateFromAwsRule(rule network_security_group.AwsSecurityGroupRule) *resourcespb.NetworkSecurityRule {
	protocol := rule.Protocol
	if rule.Protocol == "-1" {
		protocol = "*"
	}
	return &resourcespb.NetworkSecurityRule{
		Protocol: protocol,
		PortRange: &resourcespb.PortRange{
			From: int32(rule.FromPort),
			To:   int32(rule.ToPort),
		},
		CidrBlock: rule.CidrBlocks[0],
		Direction: resourcespb.Direction_INGRESS,
	}
}

func (r AwsNetworkSecurityGroup) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	awsRules := translateAwsNsgRules(r.Args.Rules)

	allowVpcTraffic := r.getAllowVpcTrafficRule()

	awsRules[INGRESS] = append(awsRules[INGRESS], allowVpcTraffic)
	awsRules[EGRESS] = append(awsRules[EGRESS], allowVpcTraffic)

	vnId, err := resources.GetMainOutputId(AwsVirtualNetwork{r.VirtualNetwork})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{
		network_security_group.AwsSecurityGroup{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			VpcId:       vnId,
			Name:        r.Args.Name,
			Description: "Managed by Multy",
			Ingress:     awsRules[INGRESS],
			Egress:      awsRules[EGRESS],
		},
	}, nil
}

func (r AwsNetworkSecurityGroup) getAllowVpcTrafficRule() network_security_group.AwsSecurityGroupRule {
	return network_security_group.AwsSecurityGroupRule{
		Protocol:   "-1",
		FromPort:   0,
		ToPort:     0,
		CidrBlocks: []string{r.VirtualNetwork.Args.CidrBlock},
	}
}
func translateAwsNsgRules(rules []*resourcespb.NetworkSecurityRule) map[string][]network_security_group.AwsSecurityGroupRule {
	awsRules := map[string][]network_security_group.AwsSecurityGroupRule{}

	for _, rule := range rules {
		awsRule := translateAwsRule(rule)
		if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
			awsRules[INGRESS] = append(awsRules[INGRESS], awsRule)
			awsRules[EGRESS] = append(awsRules[EGRESS], awsRule)
		} else if rule.Direction == resourcespb.Direction_EGRESS {
			awsRules[EGRESS] = append(awsRules[EGRESS], awsRule)
		} else if rule.Direction == resourcespb.Direction_INGRESS {
			awsRules[INGRESS] = append(awsRules[INGRESS], awsRule)
		}
	}
	return awsRules
}

func translateAwsRule(rule *resourcespb.NetworkSecurityRule) network_security_group.AwsSecurityGroupRule {
	awsFromPort := int(rule.PortRange.From)
	awsToPort := int(rule.PortRange.To)

	awsProtocol := rule.Protocol
	if rule.Protocol == "*" {
		awsProtocol = "-1"
		awsFromPort = 0
		awsToPort = 0
	}

	return network_security_group.AwsSecurityGroupRule{
		Protocol:   awsProtocol,
		FromPort:   awsFromPort,
		ToPort:     awsToPort,
		CidrBlocks: []string{rule.CidrBlock},
	}
}

func (r AwsNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return network_security_group.AwsSecurityGroupResourceName, nil
}
