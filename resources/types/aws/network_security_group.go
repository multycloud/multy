package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/types"
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

func (r AwsNetworkSecurityGroup) FromState(_ *output.TfState) (*resourcespb.NetworkSecurityGroupResource, error) {
	return &resourcespb.NetworkSecurityGroupResource{
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
	}, nil
}

func (r AwsNetworkSecurityGroup) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	awsRules := translateAwsNsgRules(r.Args.Rules)

	allowVpcTraffic := network_security_group.AwsSecurityGroupRule{
		Protocol:   "-1",
		FromPort:   0,
		ToPort:     0,
		CidrBlocks: []string{r.VirtualNetwork.Args.CidrBlock},
	}

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
			Ingress:     awsRules["ingress"],
			Egress:      awsRules["egress"],
		},
	}, nil
}
func translateAwsNsgRules(rules []*resourcespb.NetworkSecurityRule) map[string][]network_security_group.AwsSecurityGroupRule {
	awsRules := map[string][]network_security_group.AwsSecurityGroupRule{}

	for _, rule := range rules {
		awsFromPort := int(rule.PortRange.From)
		awsToPort := int(rule.PortRange.To)

		awsProtocol := rule.Protocol
		if rule.Protocol == "*" {
			awsProtocol = "-1"
			awsFromPort = 0
			awsToPort = 0
		}

		if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
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
		} else if rule.Direction == resourcespb.Direction_EGRESS {
			awsRules[EGRESS] = append(
				awsRules[EGRESS], network_security_group.AwsSecurityGroupRule{
					Protocol:   awsProtocol,
					FromPort:   awsFromPort,
					ToPort:     awsToPort,
					CidrBlocks: []string{rule.CidrBlock},
				},
			)
		} else if rule.Direction == resourcespb.Direction_INGRESS {
			awsRules[INGRESS] = append(
				awsRules[INGRESS], network_security_group.AwsSecurityGroupRule{
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

func (r AwsNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return network_security_group.AwsSecurityGroupResourceName, nil
}
