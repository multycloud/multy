package resources

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
)

type CommonResourceParams struct {
	ResourceId      string
	ResourceGroupId string
	Location        string            `hcl:"location,optional"`
	Clouds          []string          `hcl:"clouds,optional"`
	RgVars          map[string]string `hcl:"rg_vars,optional"`
	DependsOn       []string
	*validate.ResourceValidationInfo
}

func (c *CommonResourceParams) GetResourceId() string {
	return c.ResourceId
}

func (c *CommonResourceParams) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{}
}

func (c *CommonResourceParams) GetTfResourceId(cloud common.CloudProvider) string {
	return getTfResourceId(c.ResourceId, cloud)
}

func (c *CommonResourceParams) GetLocation(cloud common.CloudProvider, ctx MultyContext) string {
	return ctx.GetLocationFromCommonParams(c, cloud)
}

func (c *CommonResourceParams) GetDependencies(ctx MultyContext) []CloudSpecificResource {
	var result []CloudSpecificResource
	for _, r := range ctx.Resources {
		// we add all clouds because id specified by user is cloud-agnostic
		if util.Contains(c.DependsOn, r.Resource.GetResourceId()) {
			result = append(result, r)
		}
	}

	return util.SortResourcesById(
		result, func(t CloudSpecificResource) string {
			return t.GetResourceId()
		},
	)
}
