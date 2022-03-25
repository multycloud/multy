package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
	"strconv"
)

/*
Notes
NSG can only be applied to NIC (currently done in VM creation, to be changed to separate resource)
When NSG is applied, only rules specified are allowed.
AWS: VPC traffic is always added as an extra rule
*/

type NetworkSecurityGroup struct {
	*resources.CommonResourceParams
	Name           string          `hcl:"name"`
	VirtualNetwork *VirtualNetwork `mhcl:"ref=virtual_network"`
	Rules          []RuleType      `hcl:"rules,optional"`
}

type RuleType struct {
	Protocol  string `cty:"protocol"`
	Priority  int    `cty:"priority"`
	FromPort  string `cty:"from_port"`
	ToPort    string `cty:"to_port"`
	CidrBlock string `cty:"cidr_block"`
	Direction string `cty:"direction"`
}

const (
	INGRESS = "ingress"
	EGRESS  = "egress"
	BOTH    = "both"
	ALLOW   = "allow"
	DENY    = "deny"
)

func (r *NetworkSecurityGroup) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		awsRules := translateAwsNsgRules(r, r.Rules)

		allowVpcTraffic := network_security_group.AwsSecurityGroupRule{
			Protocol:   "-1",
			FromPort:   0,
			ToPort:     0,
			CidrBlocks: []string{r.VirtualNetwork.CidrBlock},
		}

		awsRules[INGRESS] = append(awsRules[INGRESS], allowVpcTraffic)
		awsRules[EGRESS] = append(awsRules[EGRESS], allowVpcTraffic)

		vnId, err := resources.GetMainOutputId(r.VirtualNetwork, cloud)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			network_security_group.AwsSecurityGroup{
				AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
				VpcId:       vnId,
				Ingress:     awsRules["ingress"],
				Egress:      awsRules["egress"],
			},
		}, nil
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			network_security_group.AzureNsg{
				AzResource: common.NewAzResource(
					r.GetTfResourceId(cloud), r.Name, rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				Rules: translateAzNsgRules(r.Rules),
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func translateAwsNsgRules(r *NetworkSecurityGroup, rules []RuleType) map[string][]network_security_group.AwsSecurityGroupRule {
	awsRules := map[string][]network_security_group.AwsSecurityGroupRule{}

	for _, rule := range rules {
		var awsFromPort, awsToPort int
		var awsProtocol string
		var err error

		if rule.FromPort == "*" {
			awsFromPort = 0
		} else {
			awsFromPort, err = strconv.Atoi(rule.FromPort)
			if err != nil {
				r.LogFatal(r.ResourceId, "rules", fmt.Sprintf("Invalid FromPort: %s", err.Error()))
			}
		}

		if rule.ToPort == "*" {
			awsToPort = 0
		} else {
			awsToPort, err = strconv.Atoi(rule.FromPort)
			if err != nil {
				r.LogFatal(r.ResourceId, "rules", fmt.Sprintf("Invalid ToPort: %s", err.Error()))
			}
		}

		awsProtocol = rule.Protocol
		if rule.Protocol == "*" {
			awsProtocol = "-1"
			awsFromPort = 0
			awsToPort = 0
		}

		if rule.Direction == BOTH {
			awsRules[INGRESS] = append(
				awsRules[INGRESS], network_security_group.AwsSecurityGroupRule{
					Protocol:   awsProtocol,
					FromPort:   awsFromPort,
					ToPort:     awsToPort,
					CidrBlocks: []string{rule.CidrBlock},
				},
			)
			awsRules[EGRESS] = append(
				awsRules[EGRESS], network_security_group.AwsSecurityGroupRule{
					Protocol:   awsProtocol,
					FromPort:   awsFromPort,
					ToPort:     awsToPort,
					CidrBlocks: []string{rule.CidrBlock},
				},
			)
		} else {
			awsRules[rule.Direction] = append(
				awsRules[rule.Direction], network_security_group.AwsSecurityGroupRule{
					Protocol:   awsProtocol,
					FromPort:   awsFromPort,
					ToPort:     awsToPort,
					CidrBlocks: []string{rule.CidrBlock},
				},
			)
		}
	}
	return awsRules
}

func translateAzNsgRules(rules []RuleType) []network_security_group.AzureRule {
	m := map[string]string{
		"ingress": "Inbound",
		"egress":  "Outbound",
	}

	var rls []network_security_group.AzureRule

	for _, rule := range rules {
		if rule.Direction == BOTH {
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 rule.Protocol,
					Priority:                 rule.Priority,
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     fmt.Sprintf("%s-%s", rule.ToPort, rule.FromPort),
					DestinationAddressPrefix: "*",
					Direction:                m[INGRESS],
				},
			)
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 rule.Protocol,
					Priority:                 rule.Priority,
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     fmt.Sprintf("%s-%s", rule.ToPort, rule.FromPort),
					DestinationAddressPrefix: "*",
					Direction:                m[EGRESS],
				},
			)
		} else {
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 rule.Protocol,
					Priority:                 rule.Priority,
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     fmt.Sprintf("%s-%s", rule.ToPort, rule.FromPort),
					DestinationAddressPrefix: "*",
					Direction:                m[rule.Direction],
				},
			)
		}
	}

	return rls
}

func validateRuleDirection(s string) bool {
	return s == INGRESS || s == EGRESS || s == BOTH
}

func validateRuleAction(s string) bool {
	return s == ALLOW || s == DENY
}

func validatePort(s string) bool {
	if i, err := strconv.Atoi(s); err != nil || i < -1 {
		return false
	}
	return true
}

func (r *NetworkSecurityGroup) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	for _, rule := range r.Rules {
		// TODO: get better source ranges
		if !validateRuleDirection(rule.Direction) {
			r.NewError("rules", fmt.Sprintf(
				"rule direction \"%s\" is not valid. direction must be \"%s\", \"%s\" or \"%s\"", rule.Direction,
				INGRESS, EGRESS, BOTH,
			))
		}
		if !validatePort(rule.ToPort) {
			errs = append(errs, r.NewError("rules", fmt.Sprintf("rule to_port \"%s\" is not valid", rule.ToPort)))
		}
		if !validatePort(rule.FromPort) {
			errs = append(errs, r.NewError("rules", fmt.Sprintf("rule from_port \"%s\" is not valid", rule.FromPort)))
		}
		// TODO validate CIDR
		//		validate protocol
	}
	// TODO validate location matches with VN location
	return errs
}

func (r *NetworkSecurityGroup) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return network_security_group.AwsSecurityGroupResourceName, nil
	case common.AZURE:
		return network_security_group.AzureNetworkSecurityGroupResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
