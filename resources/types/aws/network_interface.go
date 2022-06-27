package aws_resources

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

type AwsNetworkInterface struct {
	*types.NetworkInterface
}

func InitNetworkInterface(r *types.NetworkInterface) resources.ResourceTranslator[*resourcespb.NetworkInterfaceResource] {
	return AwsNetworkInterface{r}
}

func (r AwsNetworkInterface) FromState(_ *output.TfState) (*resourcespb.NetworkInterfaceResource, error) {
	return &resourcespb.NetworkInterfaceResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:             r.Args.Name,
		SubnetId:         r.Args.SubnetId,
		PublicIpId:       r.Args.PublicIpId,
		AvailabilityZone: r.Args.AvailabilityZone,
	}, nil
}

func (r AwsNetworkInterface) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var pIpId string
	subnetId, err := AwsSubnet{r.Subnet}.GetSubnetId(r.Args.AvailabilityZone)
	if err != nil {
		return nil, err
	}
	if r.PublicIp != nil {
		pIpId, err = resources.GetMainOutputId(AwsPublicIp{r.PublicIp})
		if err != nil {
			return nil, err
		}
	}

	var res []output.TfBlock
	nic := network_interface.AwsNetworkInterface{
		AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
		SubnetId:    subnetId,
	}
	if pIpId != "" {
		res = append(res, network_interface.AwsEipAssociation{
			AwsResource: &common.AwsResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			},
			AllocationId:       pIpId,
			NetworkInterfaceId: fmt.Sprintf("%s.%s.id", output.GetResourceName(nic), nic.ResourceId),
		})
	}

	res = append(res, nic)

	return res, nil
}

func (r AwsNetworkInterface) GetMainResourceName() (string, error) {
	return network_interface.AwsResourceName, nil
}
