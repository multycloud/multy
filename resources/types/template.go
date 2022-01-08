package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"

	"github.com/zclconf/go-cty/cty"
)

// teamplte
type template struct {
	ResourceId      string
	ResourceGroupId string
	Name            string `hcl:"name"`
	Optional        bool   `hcl:"optional,optional""`
}

func (r *template) Translate(cloud common.CloudProvider, ctx resources.MultyContext) interface{} {
	if cloud == common.AWS {
		return []interface{}{}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		fmt.Sprintln(rgName)
		return []interface{}{}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *template) GetContext(cloud common.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{}
}

func (r *template) getResourceId(cloud common.CloudProvider) string {
	return fmt.Sprintf("%s_%s", r.ResourceId, cloud)
}
