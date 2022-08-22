package aws_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
)

type ResourceGroup struct {
	*types.ResourceGroup
}

func InitResourceGroup(rg *types.ResourceGroup) resources.ResourceTranslator[*resourcespb.ResourceGroupResource] {
	return ResourceGroup{rg}
}

func (rg ResourceGroup) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	return nil, nil
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
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", rg.GetCloud())
}
