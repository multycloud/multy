package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/network_interface"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

type NetworkInterface struct {
	*resources.CommonResourceParams
	Name     string `hcl:"name"`
	SubnetId string `hcl:"subnet_id,optional"`
}

func (r *NetworkInterface) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	var subnetId string
	if s, err := ctx.GetResource(r.SubnetId); err != nil {
		r.LogFatal(r.ResourceId, "subnet_id", err.Error())
	} else {
		subnetId = s.Resource.(*Subnet).GetId(cloud)
	}

	if cloud == common.AWS {
		return []interface{}{
			network_interface.AwsNetworkInterface{
				AwsResource: common.AwsResource{
					ResourceName: network_interface.AwsResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
					Tags:         map[string]string{"Name": r.Name},
				},
				SubnetId: subnetId,
			},
		}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		nic := network_interface.AzureNetworkInterface{
			AzResource: common.AzResource{
				ResourceName:      network_interface.AzureResourceName,
				ResourceId:        r.GetTfResourceId(cloud),
				ResourceGroupName: rgName,
				Name:              r.Name,
				Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			},
			IpConfigurations: []network_interface.AzureIpConfiguration{{
				Name:                       "internal", // this name shouldn't be vm.name
				PrivateIpAddressAllocation: "Dynamic",
				SubnetId:                   subnetId,
				Primary:                    true,
			}},
		}
		if publicIpReference := getPublicIpReferences(ctx, *r, subnetId); len(publicIpReference) != 0 {
			nic.IpConfigurations = publicIpReference
		}
		return []interface{}{nic}
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *NetworkInterface) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: network_interface.AwsResourceName, common.AZURE: network_interface.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func getPublicIpReferences(ctx resources.MultyContext, nic NetworkInterface, subnetId string) []network_interface.AzureIpConfiguration {
	var ipConfigurations []network_interface.AzureIpConfiguration
	for _, resource := range ctx.Resources {
		switch resource.Resource.(type) {
		case *PublicIp:
			r := resource.Resource.(*PublicIp)
			if resources.GetCloudSpecificResourceId(r, common.AZURE) == r.NetworkInterfaceId {
				ipConfigurations = append(ipConfigurations, network_interface.AzureIpConfiguration{
					Name:                       fmt.Sprintf("external-%s", r.Name),
					PrivateIpAddressAllocation: "Dynamic",
					PublicIpAddressId:          r.GetId(common.AZURE),
					SubnetId:                   subnetId,
					Primary:                    true,
				})
			}
		}
	}
	return ipConfigurations
}

func (r *NetworkInterface) Validate(ctx resources.MultyContext) {
	return
}
