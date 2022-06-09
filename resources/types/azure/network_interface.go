package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/resources/types"
)

type AzureNetworkInterface struct {
	*types.NetworkInterface
}

func InitNetworkInterface(r *types.NetworkInterface) resources.ResourceTranslator[*resourcespb.NetworkInterfaceResource] {
	return AzureNetworkInterface{r}
}

func (r AzureNetworkInterface) FromState(_ *output.TfState) (*resourcespb.NetworkInterfaceResource, error) {
	return &resourcespb.NetworkInterfaceResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:       r.Args.Name,
		SubnetId:   r.Args.SubnetId,
		PublicIpId: r.Args.PublicIpId,
	}, nil
}

func (r AzureNetworkInterface) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var pIpId string
	subnetId, err := resources.GetMainOutputId(AzureSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}
	if r.PublicIp != nil {
		pIpId, err = resources.GetMainOutputId(AzurePublicIp{r.PublicIp})
		if err != nil {
			return nil, err
		}
	}
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
	if pIpId != "" {
		nic.IpConfigurations = []network_interface.AzureIpConfiguration{{
			Name:                       fmt.Sprintf("external-%s", r.Args.Name),
			PrivateIpAddressAllocation: "Dynamic",
			PublicIpAddressId:          pIpId,
			SubnetId:                   subnetId,
			Primary:                    true,
		}}
	}
	return []output.TfBlock{nic}, nil
}

func (r AzureNetworkInterface) GetMainResourceName() (string, error) {
	return network_interface.AzureResourceName, nil
}
