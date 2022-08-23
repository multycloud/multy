package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
)

type ResourceGroup struct {
	*types.ResourceGroup
}

func InitResourceGroup(rg *types.ResourceGroup) resources.ResourceTranslator[*resourcespb.ResourceGroupResource] {
	return ResourceGroup{rg}
}

type AzureResourceGroup struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_resource_group"`
	Location           string `hcl:"location"`
}

func (rg ResourceGroup) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	allDeps := rg.GetAllDependentResources(ctx.Resources)
	if len(allDeps) == 0 {
		// if no resources are in this group, just don't output anything
		return nil, nil
	}
	for _, r := range allDeps {
		ctx.Resources.AddDependency(rg.ResourceId, r)
	}

	return []output.TfBlock{AzureResourceGroup{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: rg.ResourceId},
			Name:              rg.Args.Name,
		},
		Location: rg.GetCloudSpecificLocation(),
	}}, nil

}

func (rg ResourceGroup) FromState(_ *output.TfState, _ *output.TfPlan) (*resourcespb.ResourceGroupResource, error) {
	return &resourcespb.ResourceGroupResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:    rg.ResourceId,
			Location:      rg.Args.CommonParameters.Location,
			CloudProvider: rg.Args.CommonParameters.CloudProvider,
		},
	}, nil
}

func (rg ResourceGroup) GetMainResourceName() (string, error) {
	return "azurerm_resource_group", nil
}

func GetResourceGroupName(name string) string {
	return fmt.Sprintf("azurerm_resource_group.%s.name", name)
}
