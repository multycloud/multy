package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"

	"github.com/zclconf/go-cty/cty"
)

// teamplte
type template struct {
	ResourceId      string
	ResourceGroupId string
	Name            string `hcl:"name"`
	Optional        bool   `hcl:"optional,optional""`
}

func (r *template) Translate(cloud common.CloudProvider, ctx resources.MultyContext) any {
	if cloud == common.AWS {
		return []any{}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		fmt.Sprintln(rgName)
		return []any{}
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
