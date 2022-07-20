package azure_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/types"
)

type AzurePublicIp struct {
	*types.PublicIp
}

func InitPublicIp(r *types.PublicIp) resources.ResourceTranslator[*resourcespb.PublicIpResource] {
	return AzurePublicIp{r}
}

func (r AzurePublicIp) FromState(state *output.TfState) (*resourcespb.PublicIpResource, error) {
	out := &resourcespb.PublicIpResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Ip:          "dryrun",
		GcpOverride: r.Args.GcpOverride,
	}
	if flags.DryRun {
		return out, nil
	}

	stateResource, err := output.GetParsedById[public_ip.AzurePublicIp](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	out.Ip = stateResource.IpAddress
	out.AzureOutputs = &resourcespb.PublicIpAzureOutputs{
		PublicIpId: stateResource.ResourceId,
	}

	return out, nil
}

func (r AzurePublicIp) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		public_ip.AzurePublicIp{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.GetCommonParameters().ResourceGroupId),
				r.GetCloudSpecificLocation(),
			),
			AllocationMethod: "Static",
			Sku:              "Standard",
		},
	}, nil
}

func (r AzurePublicIp) GetMainResourceName() (string, error) {
	return public_ip.AzureResourceName, nil
}
