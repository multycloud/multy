package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/public_ip"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
)

/*
Notes
AWS: NIC association done on public_ip
Azure: NIC association done on NIC creation
*/

type PublicIp struct {
	*resources.CommonResourceParams
	Name               string            `hcl:"name"`
	NetworkInterfaceId *NetworkInterface `mhcl:"ref=network_interface_id,optional"`
}

func (r *PublicIp) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AWS {
		nid := ""
		if r.NetworkInterfaceId != nil {
			var err error
			nid, err = resources.GetMainOutputId(r.NetworkInterfaceId, cloud)
			if err != nil {
				return nil, err
			}
		}
		return []output.TfBlock{
			public_ip.AwsElasticIp{
				AwsResource:        common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
				NetworkInterfaceId: nid,
				//Vpc:        true,
			},
		}, nil
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			public_ip.AzurePublicIp{
				AzResource: common.NewAzResource(
					r.GetTfResourceId(cloud), r.Name, rg.GetResourceGroupName(r.ResourceGroupId, cloud),
					ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				AllocationMethod: "Static",
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *PublicIp) GetId(cloud common.CloudProvider) string {
	types := map[common.CloudProvider]string{common.AWS: public_ip.AwsResourceName, common.AZURE: public_ip.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.GetTfResourceId(cloud))
}

func (r *PublicIp) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	errs = append(errs, r.CommonResourceParams.Validate(ctx, cloud)...)
	//if r.NetworkInterfaceId != "" && r.InstanceId != "" {
	//	errs = append(errs, r.NewError(r.ResourceId, "instance_id", "cannot set both network_interface_id and instance_id"))
	//}
	return errs
}

func (r *PublicIp) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return public_ip.AwsResourceName, nil
	case common.AZURE:
		return public_ip.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
