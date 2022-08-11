package azure_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface_security_group_association"
	"github.com/multycloud/multy/resources/types"
)

type AzureNetworkInterfaceSecurityGroupAssociation struct {
	*types.NetworkInterfaceSecurityGroupAssociation
}

func InitNetworkInterfaceSecurityGroupAssociation(vn *types.NetworkInterfaceSecurityGroupAssociation) resources.ResourceTranslator[*resourcespb.NetworkInterfaceSecurityGroupAssociationResource] {
	return AzureNetworkInterfaceSecurityGroupAssociation{vn}
}

func (r AzureNetworkInterfaceSecurityGroupAssociation) FromState(state *output.TfState) (*resourcespb.NetworkInterfaceSecurityGroupAssociationResource, error) {
	out := &resourcespb.NetworkInterfaceSecurityGroupAssociationResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		NetworkInterfaceId: r.Args.NetworkInterfaceId,
		SecurityGroupId:    r.Args.SecurityGroupId,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if _, exists, _ := output.MaybeGetParsedById[network_interface_security_group_association.AzureNetworkInterfaceSecurityGroupAssociation](state, r.ResourceId); !exists {
		statuses["azure_network_interface_security_group_association"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}

	return out, nil
}

func (r AzureNetworkInterfaceSecurityGroupAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	nicId, err := resources.GetMainOutputId(AzureNetworkInterface{r.NetworkInterface})
	if err != nil {
		return nil, err
	}
	nsgId, err := resources.GetMainOutputId(AzureNetworkSecurityGroup{r.NetworkSecurityGroup})
	if err != nil {
		return nil, err
	}
	return []output.TfBlock{network_interface_security_group_association.AzureNetworkInterfaceSecurityGroupAssociation{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
		},
		NetworkInterfaceId: nicId,
		SecurityGroupId:    nsgId,
	}}, nil
}

func (r AzureNetworkInterfaceSecurityGroupAssociation) GetMainResourceName() (string, error) {
	return network_interface_security_group_association.AzureResourceName, nil
}
