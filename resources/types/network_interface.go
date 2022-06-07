package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/validate"
)

var networkInterfaceMetadata = resources.ResourceMetadata[*resourcespb.NetworkInterfaceArgs, *NetworkInterface, *resourcespb.NetworkInterfaceResource]{
	CreateFunc:        CreateNetworkInterface,
	UpdateFunc:        UpdateNetworkInterface,
	ReadFromStateFunc: NetworkInterfaceFromState,
	ExportFunc: func(r *NetworkInterface, _ *resources.Resources) (*resourcespb.NetworkInterfaceArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewNetworkInterface,
	AbbreviatedName: "nic",
}

type NetworkInterface struct {
	resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]

	Subnet   *Subnet
	PublicIp *PublicIp
}

func (r *NetworkInterface) GetMetadata() resources.ResourceMetadataInterface {
	return &networkInterfaceMetadata
}

func CreateNetworkInterface(resourceId string, args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) (*NetworkInterface, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return nil, err
		}
		rgId, err := NewRgFromParent("nic", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return nil, err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewNetworkInterface(resourceId, args, others)
}

func UpdateNetworkInterface(resource *NetworkInterface, vn *resourcespb.NetworkInterfaceArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func NetworkInterfaceFromState(resource *NetworkInterface, _ *output.TfState) (*resourcespb.NetworkInterfaceResource, error) {
	return &resourcespb.NetworkInterfaceResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:       resource.Args.Name,
		SubnetId:   resource.Args.SubnetId,
		PublicIpId: resource.Args.PublicIpId,
	}, nil
}

func NewNetworkInterface(resourceId string, args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) (*NetworkInterface, error) {
	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return nil, err
	}
	pIp, _, err := resources.GetOptional[*PublicIp](resourceId, others, args.PublicIpId)
	if err != nil {
		return nil, err
	}
	return &NetworkInterface{
		ResourceWithId: resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		Subnet:   subnet,
		PublicIp: pIp,
	}, nil
}

func (r *NetworkInterface) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var pIp string
	subnetId, err := resources.GetMainOutputId(r.Subnet)
	if err != nil {
		return nil, err
	}
	if r.PublicIp != nil {
		pIp, err = resources.GetMainOutputId(r.PublicIp)
		if err != nil {
			return nil, err
		}
	}

	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var res []output.TfBlock
		nic := network_interface.AwsNetworkInterface{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			SubnetId:    subnetId,
		}
		if pIp != "" {
			res = append(res, network_interface.AwsEipAssociation{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.Subnet.ResourceId},
				},
				AllocationId:       pIp,
				NetworkInterfaceId: fmt.Sprintf("%s.%s.id", output.GetResourceName(nic), nic.ResourceId),
			})
		}

		res = append(res, nic)

		return res, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
		nic := network_interface.AzureNetworkInterface{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rgName,
				r.GetCloudSpecificLocation(),
			),
			// by default, virtual_machine will have a private ip
			IpConfigurations: []network_interface.AzureIpConfiguration{{
				Name:                       "internal", // this name shouldn't be vm.name
				PrivateIpAddressAllocation: "Dynamic",
				SubnetId:                   subnetId,
				Primary:                    true,
			}},
		}
		// associate a public ip configuration in case a public_ip resource references this network_interface
		if pIp != "" {
			nic.IpConfigurations = []network_interface.AzureIpConfiguration{{
				Name:                       fmt.Sprintf("external-%s", r.Args.Name),
				PrivateIpAddressAllocation: "Dynamic",
				PublicIpAddressId:          pIp,
				SubnetId:                   subnetId,
				Primary:                    true,
			}}
		}
		return []output.TfBlock{nic}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *NetworkInterface) GetId(cloud commonpb.CloudProvider) string {
	types := map[commonpb.CloudProvider]string{common.AWS: network_interface.AwsResourceName, common.AZURE: network_interface.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.ResourceId)
}

//func (r *NetworkInterface) getPublicIpReferences(ctx resources.MultyContext, subnetId string) []network_interface.AzureIpConfiguration {
//	var ipConfigurations []network_interface.AzureIpConfiguration
//	for _, resource := range resources.GetAllResourcesWithRef(ctx, func(i *PublicIp) *NetworkInterface { return i.NetworkInterface }, r) {
//		ipConfigurations = append(
//			ipConfigurations, network_interface.AzureIpConfiguration{
//				Name:                       fmt.Sprintf("external-%s", resource.Args.Name),
//				PrivateIpAddressAllocation: "Dynamic",
//				PublicIpAddressId:          resource.GetId(common.AZURE),
//				SubnetId:                   subnetId,
//				Primary:                    true,
//			},
//		)
//	}
//	return ipConfigurations
//}

func (r *NetworkInterface) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}

func (r *NetworkInterface) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return network_interface.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return network_interface.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
