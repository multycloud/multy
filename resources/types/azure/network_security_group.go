package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/types"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"
	"sort"
	"strconv"
	"strings"
)

type AzureNetworkSecurityGroup struct {
	*types.NetworkSecurityGroup
}

func InitNetworkSecurityGroup(r *types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource] {
	return AzureNetworkSecurityGroup{r}
}

func (r AzureNetworkSecurityGroup) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.NetworkSecurityGroupResource, error) {
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

	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.AzureNsg](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		singleDirectionRules := map[string]*resourcespb.NetworkSecurityRule{}
		bothDirectionRules := map[string]*resourcespb.NetworkSecurityRule{}
		for _, rule := range stateResource.Rules {
			direction := convertToDirection(rule.Direction)
			cidrBlock := ""
			if direction == resourcespb.Direction_INGRESS {
				cidrBlock = rule.SourceAddressPrefix
			} else {
				cidrBlock = rule.DestinationAddressPrefix
			}
			if cidrBlock == "*" {
				cidrBlock = "0.0.0.0/0"
			}
			inRule := &resourcespb.NetworkSecurityRule{
				Protocol:  strings.ToLower(rule.Protocol),
				Priority:  int64(rule.Priority),
				PortRange: translateToPortRange(rule.DestinationPortRange),
				CidrBlock: cidrBlock,
				Direction: direction,
			}

			if strings.HasSuffix(rule.Name, "-"+rule.Direction) {
				if ruleNumber, err := strconv.Atoi(strings.TrimSuffix(rule.Name, "-"+rule.Direction)); err == nil {
					var matchingRuleName string
					if direction == resourcespb.Direction_INGRESS {
						matchingRuleName = getRuleNameForBidirectionalRule(ruleNumber, resourcespb.Direction_EGRESS)
					} else {
						matchingRuleName = getRuleNameForBidirectionalRule(ruleNumber, resourcespb.Direction_INGRESS)
					}

					if matchingRule, ok := singleDirectionRules[matchingRuleName]; ok {
						if matchingRule.Protocol == inRule.Protocol &&
							matchingRule.Priority == inRule.Priority &&
							proto.Equal(matchingRule.PortRange, inRule.PortRange) &&
							matchingRule.CidrBlock == inRule.CidrBlock &&
							matchingRule.Direction != inRule.Direction {
							matchingRule.Direction = resourcespb.Direction_BOTH_DIRECTIONS
							bothDirectionRules[strconv.Itoa(ruleNumber)] = matchingRule
							delete(singleDirectionRules, matchingRuleName)
							continue
						}
					}
				}
			}
			singleDirectionRules[rule.Name] = inRule
		}

		var keys []string
		keys = append(keys, maps.Keys(singleDirectionRules)...)
		keys = append(keys, maps.Keys(bothDirectionRules)...)
		sort.Strings(keys)
		out.Rules = nil
		for _, key := range keys {
			if rule, ok := singleDirectionRules[key]; ok {
				out.Rules = append(out.Rules, rule)
			}
			if rule, ok := bothDirectionRules[key]; ok {
				out.Rules = append(out.Rules, rule)
			}
		}
		out.AzureOutputs = &resourcespb.NetworkSecurityGroupAzureOutputs{NetworkSecurityGroupId: stateResource.ResourceId}
		output.AddToStatuses(statuses, "azure_network_security_group", output.MaybeGetPlannedChageById[network_security_group.AzureNsg](plan, r.ResourceId))
	} else {
		statuses["azure_network_security_group"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r AzureNetworkSecurityGroup) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		network_security_group.AzureNsg{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId),
				r.GetCloudSpecificLocation(),
			),
			Rules: translateAzNsgRules(r.Args.Rules),
		},
	}, nil
}

func translateAzNsgRules(rules []*resourcespb.NetworkSecurityRule) []network_security_group.AzureRule {
	var rls []network_security_group.AzureRule

	for i, rule := range rules {
		protocol := strings.Title(strings.ToLower(rule.Protocol))
		cidrBlock := rule.CidrBlock
		if cidrBlock == "0.0.0.0/0" {
			cidrBlock = "*"
		}
		if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     getRuleNameForBidirectionalRule(i, resourcespb.Direction_INGRESS),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      cidrBlock,
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: "*",
					Direction:                convertDirection(resourcespb.Direction_INGRESS),
				},
			)
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     getRuleNameForBidirectionalRule(i, resourcespb.Direction_EGRESS),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: cidrBlock,
					Direction:                convertDirection(resourcespb.Direction_EGRESS),
				},
			)
		} else {
			sourceAddress := "*"
			destinationAddress := "*"
			if rule.Direction == resourcespb.Direction_INGRESS {
				sourceAddress = cidrBlock
			} else {
				destinationAddress = cidrBlock
			}
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(i),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      sourceAddress,
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: destinationAddress,
					Direction:                convertDirection(rule.Direction),
				},
			)
		}
	}

	return rls
}

func translatePortRange(pr *resourcespb.PortRange) string {
	from := "*"
	if pr.From != 0 {
		from = strconv.Itoa(int(pr.From))
	}
	to := "*"
	if pr.To != 0 {
		to = strconv.Itoa(int(pr.To))
	}
	return fmt.Sprintf("%s-%s", from, to)
}

func translateToPortRange(pr string) *resourcespb.PortRange {
	ports := strings.SplitN(pr, "-", 2)
	if len(ports) == 0 {
		return nil
	}
	result := &resourcespb.PortRange{}
	if from, err := strconv.Atoi(ports[0]); err == nil {
		result.From = int32(from)
	}
	if len(ports) == 1 {
		return result
	}
	if to, err := strconv.Atoi(ports[1]); err == nil {
		result.To = int32(to)
	}

	return result
}

func convertDirection(d resourcespb.Direction) string {
	m := map[resourcespb.Direction]string{
		resourcespb.Direction_INGRESS: "Inbound",
		resourcespb.Direction_EGRESS:  "Outbound",
	}
	return m[d]
}

func convertToDirection(d string) resourcespb.Direction {
	m := map[string]resourcespb.Direction{
		"Inbound":  resourcespb.Direction_INGRESS,
		"Outbound": resourcespb.Direction_EGRESS,
	}
	return m[d]
}

func getRuleNameForBidirectionalRule(i int, d resourcespb.Direction) string {
	return fmt.Sprintf("%d-%s", i, convertDirection(d))
}

func (r AzureNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return network_security_group.AzureNetworkSecurityGroupResourceName, nil
}
