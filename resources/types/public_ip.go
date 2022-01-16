package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/public_ip"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

/*
Notes
AWS: NIC association done on public_ip
Azure: NIC association done on NIC creation
*/

type PublicIp struct {
	*resources.CommonResourceParams
	Name               string `hcl:"name"`
	NetworkInterfaceId string `hcl:"network_interface_id,optional"`
}

func (r *PublicIp) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {

	if cloud == common.AWS {
		return []interface{}{
			public_ip.AwsElasticIp{
				AwsResource: common.AwsResource{
					ResourceName: public_ip.AwsResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
					Tags:         map[string]string{"Name": r.Name},
				},
				NetworkInterfaceId: r.NetworkInterfaceId,
				//Vpc:        true,
			},
		}
	} else if cloud == common.AZURE {
		return []interface{}{
			public_ip.AzurePublicIp{
				AzResource: common.AzResource{
					ResourceName:      public_ip.AzureResourceName,
					ResourceId:        r.GetTfResourceId(cloud),
					ResourceGroupName: rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					Name:              r.Name,
					Location:          ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				},
				AllocationMethod: "Static",
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *PublicIp) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: public_ip.AwsResourceName, common.AZURE: public_ip.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func (r *PublicIp) Validate(ctx resources.MultyContext) {
	//if r.NetworkInterfaceId != "" && r.InstanceId != "" {
	//	r.LogFatal(r.ResourceId, "instance_id", "cannot set both network_interface_id and instance_id")
	//}
	return
}

func (r *PublicIp) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return public_ip.AwsResourceName
	case common.AZURE:
		return public_ip.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
