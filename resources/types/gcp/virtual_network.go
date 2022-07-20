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
)

type GcpVirtualNetwork struct {
	*types.VirtualNetwork
}

func InitVirtualNetwork(vn *types.VirtualNetwork) resources.ResourceTranslator[*resourcespb.VirtualNetworkResource] {
	return GcpVirtualNetwork{vn}
}

func (r GcpVirtualNetwork) FromState(state *output.TfState) (*resourcespb.VirtualNetworkResource, error) {
	out := &resourcespb.VirtualNetworkResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.GetCloud(),
		},
		Name:        r.Args.Name,
		CidrBlock:   r.Args.CidrBlock,
		GcpOverride: r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	stateResource, err := output.GetParsedById[virtual_network.GoogleComputeNetwork](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	out.GcpOutputs = &resourcespb.VirtualNetworkGcpOutputs{
		ComputeNetworkId: stateResource.SelfLink,
	}

	if stateResource, exists, err := output.MaybeGetParsedById[network_security_group.GoogleComputeFirewall](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.GcpOutputs.DefaultComputeFirewallId = stateResource.SelfLink
	}

	return out, nil
}

func (r GcpVirtualNetwork) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var result []output.TfBlock

	// gcp allows egress traffic by default: https://cloud.google.com/vpc/docs/firewalls
	// so we add default deny rule to egress
	ruleName := fmt.Sprintf("%s-%s", r.Args.Name, "default-deny-egress")
	result = append(result, &network_security_group.GoogleComputeFirewall{
		GcpResource:       common.NewGcpResource(r.ResourceId, ruleName, r.Args.GetGcpOverride().GetProject()),
		Network:           fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
		DestinationRanges: []string{"0.0.0.0/0"},
		Direction:         resourcespb.Direction_EGRESS.String(),
		DenyRules:         []network_security_group.GoogleComputeFirewallRule{{Protocol: "all"}},
		TargetTags:        []string{r.getVnTag()},
		Priority:          65535,
	})

	result = append(result, &virtual_network.GoogleComputeNetwork{
		GcpResource:                 common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
		RoutingMode:                 "REGIONAL",
		Description:                 "Managed by Multy",
		AutoCreateSubnetworks:       false,
		DeleteDefaultRoutesOnCreate: true,
	})

	return result, nil
}

func (r GcpVirtualNetwork) getVnTag() string {
	return fmt.Sprintf("vn-%s", r.Args.Name)
}

func (r GcpVirtualNetwork) GetMainResourceName() (string, error) {
	return virtual_network.GcpResourceName, nil
}
