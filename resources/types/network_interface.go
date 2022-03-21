package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

type NetworkInterface struct {
	*resources.CommonResourceParams
	Name     string  `hcl:"name"`
	SubnetId *Subnet `mhcl:"ref=subnet_id,optional"`
}

func (r *NetworkInterface) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	subnetId := resources.GetMainOutputId(r.SubnetId, cloud)

	if cloud == common.AWS {
		return []output.TfBlock{
			network_interface.AwsNetworkInterface{
				AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
				SubnetId:    subnetId,
			},
		}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		nic := network_interface.AzureNetworkInterface{
			AzResource: common.NewAzResource(
				r.GetTfResourceId(cloud), r.Name, rgName,
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
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
		if publicIpReference := getPublicIpReferences(ctx, subnetId); len(publicIpReference) != 0 {
			nic.IpConfigurations = publicIpReference
		}
		return []output.TfBlock{nic}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *NetworkInterface) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: network_interface.AwsResourceName, common.AZURE: network_interface.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func getPublicIpReferences(ctx resources.MultyContext, subnetId string) []network_interface.AzureIpConfiguration {
	var ipConfigurations []network_interface.AzureIpConfiguration
	for _, resource := range resources.GetAllResources[*PublicIp](ctx) {
		if resource.NetworkInterfaceId != nil && resources.GetCloudSpecificResourceId(resource, common.AZURE) == resource.NetworkInterfaceId.ResourceId {
			ipConfigurations = append(
				ipConfigurations, network_interface.AzureIpConfiguration{
					Name:                       fmt.Sprintf("external-%s", resource.Name),
					PrivateIpAddressAllocation: "Dynamic",
					PublicIpAddressId:          resource.GetId(common.AZURE),
					SubnetId:                   subnetId,
					Primary:                    true,
				},
			)
		}
	}
	return ipConfigurations
}

func (r *NetworkInterface) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *NetworkInterface) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return network_interface.AwsResourceName
	case common.AZURE:
		return network_interface.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
