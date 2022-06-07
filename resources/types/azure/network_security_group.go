package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/types"
	"strconv"
	"strings"
)

type AzureNetworkSecurityGroup struct {
	*types.NetworkSecurityGroup
}

func InitNetworkSecurityGroup(r *types.NetworkSecurityGroup) resources.ResourceTranslator[*resourcespb.NetworkSecurityGroupResource] {
	return AzureNetworkSecurityGroup{r}
}

func (r AzureNetworkSecurityGroup) FromState(_ *output.TfState) (*resourcespb.NetworkSecurityGroupResource, error) {
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
	m := map[resourcespb.Direction]string{
		resourcespb.Direction_INGRESS: "Inbound",
		resourcespb.Direction_EGRESS:  "Outbound",
	}

	var rls []network_security_group.AzureRule

	for _, rule := range rules {
		protocol := strings.Title(strings.ToLower(rule.Protocol))
		if rule.Direction == resourcespb.Direction_BOTH_DIRECTIONS {
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: "*",
					Direction:                m[resourcespb.Direction_INGRESS],
				},
			)
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: "*",
					Direction:                m[resourcespb.Direction_EGRESS],
				},
			)
		} else {
			rls = append(
				rls, network_security_group.AzureRule{
					Name:                     strconv.Itoa(len(rls)),
					Protocol:                 protocol,
					Priority:                 int(rule.Priority),
					Access:                   "Allow",
					SourcePortRange:          "*",
					SourceAddressPrefix:      "*",
					DestinationPortRange:     translatePortRange(rule.PortRange),
					DestinationAddressPrefix: "*",
					Direction:                m[rule.Direction],
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

func (r AzureNetworkSecurityGroup) GetMainResourceName() (string, error) {
	return network_security_group.AzureNetworkSecurityGroupResourceName, nil
}
