package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/resources/output/network_interface_security_group_association"
	"github.com/multycloud/multy/validate"
)

var networkInterfaceSecurityGroupAssociationMetadata = resources.ResourceMetadata[*resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, *NetworkInterfaceSecurityGroupAssociation, *resourcespb.NetworkInterfaceSecurityGroupAssociationResource]{
	CreateFunc:        CreateNetworkInterfaceSecurityGroupAssociation,
	UpdateFunc:        UpdateNetworkInterfaceSecurityGroupAssociation,
	ReadFromStateFunc: NetworkInterfaceSecurityGroupAssociationFromState,
	ExportFunc: func(r *NetworkInterfaceSecurityGroupAssociation, _ *resources.Resources) (*resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewNetworkInterfaceSecurityGroupAssociation,
	AbbreviatedName: "nic",
}

type NetworkInterfaceSecurityGroupAssociation struct {
	resources.ChildResourceWithId[*NetworkInterface, *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs]

	NetworkInterface     *NetworkInterface
	NetworkSecurityGroup *NetworkSecurityGroup
}

func (r *NetworkInterfaceSecurityGroupAssociation) GetMetadata() resources.ResourceMetadataInterface {
	return &networkInterfaceSecurityGroupAssociationMetadata
}

func CreateNetworkInterfaceSecurityGroupAssociation(resourceId string, args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) (*NetworkInterfaceSecurityGroupAssociation, error) {
	return NewNetworkInterfaceSecurityGroupAssociation(resourceId, args, others)
}

func UpdateNetworkInterfaceSecurityGroupAssociation(resource *NetworkInterfaceSecurityGroupAssociation, vn *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func NetworkInterfaceSecurityGroupAssociationFromState(resource *NetworkInterfaceSecurityGroupAssociation, _ *output.TfState) (*resourcespb.NetworkInterfaceSecurityGroupAssociationResource, error) {
	return &resourcespb.NetworkInterfaceSecurityGroupAssociationResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  resource.ResourceId,
			NeedsUpdate: false,
		},
		NetworkInterfaceId: resource.Args.NetworkInterfaceId,
		SecurityGroupId:    resource.Args.SecurityGroupId,
	}, nil
}

func NewNetworkInterfaceSecurityGroupAssociation(resourceId string, args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) (*NetworkInterfaceSecurityGroupAssociation, error) {
	nicNsgAssociation := &NetworkInterfaceSecurityGroupAssociation{
		ChildResourceWithId: resources.ChildResourceWithId[*NetworkInterface, *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	nic, err := resources.Get[*NetworkInterface](resourceId, others, args.NetworkInterfaceId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{nicNsgAssociation.NewValidationError(err, "network_interface_id")})
	}
	nicNsgAssociation.Parent = nic
	nicNsgAssociation.NetworkInterface = nic

	nsg, err := resources.Get[*NetworkSecurityGroup](resourceId, others, args.SecurityGroupId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{nicNsgAssociation.NewValidationError(err, "security_group_id")})
	}
	nicNsgAssociation.NetworkSecurityGroup = nsg
	return nicNsgAssociation, nil
}

func (r *NetworkInterfaceSecurityGroupAssociation) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	nic, err := resources.GetMainOutputId(r.NetworkInterface)
	if err != nil {
		return nil, err
	}
	nsg, err := resources.GetMainOutputId(r.NetworkSecurityGroup)
	if err != nil {
		return nil, err
	}

	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{
			network_interface_security_group_association.AwsNetworkInterfaceSecurityGroupAssociation{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.NetworkInterface.ResourceId},
				},
				NetworkInterfaceId:     nic,
				NetworkSecurityGroupId: nsg,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return []output.TfBlock{network_interface_security_group_association.AzureNetworkInterfaceSecurityGroupAssociation{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.NetworkInterface.ResourceId},
			},
			NetworkInterfaceId: nic,
			SecurityGroupId:    nsg,
		}}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *NetworkInterfaceSecurityGroupAssociation) GetId(cloud commonpb.CloudProvider) string {
	types := map[commonpb.CloudProvider]string{common.AWS: network_interface.AwsResourceName, common.AZURE: network_interface.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.ResourceId)
}

func (r *NetworkInterfaceSecurityGroupAssociation) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	return nil
}

func (r *NetworkInterfaceSecurityGroupAssociation) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return network_interface_security_group_association.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return network_interface_security_group_association.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
