package resources

import (
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources/common"
	"multy-go/validate"
)

type CommonResourceParams struct {
	ResourceId      string
	ResourceGroupId string
	Location        string            `hcl:"location,optional"`
	Clouds          []string          `hcl:"clouds,optional"`
	RgVars          map[string]string `hcl:"rg_vars,optional"`
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
