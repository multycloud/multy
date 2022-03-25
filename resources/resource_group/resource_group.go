package rg

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/hclutil"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/validate"

	"github.com/zclconf/go-cty/cty"
)

const (
	PREFIX       = "Prefix"
	ENVIRONMENT  = "Env"
	APP          = "App"
	RESOURCETYPE = "ResourceType"
	LOCATION     = "Location"
	SUFFIX       = "Suffix"
)

// https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations

type Type struct {
	ResourceId string
	Name       string `hcl:"name"`
	Location   string `hcl:"location"`
	App        string `hcl:"app"`
}

type ResourceGroup struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_resource_group"`
	Location           string `hcl:"location"`
}

const AzureResourceName = "azurerm_resource_group"

func (rg *Type) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	if cloud == common.AZURE {
		return []output.TfBlock{ResourceGroup{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: rg.ResourceId},
				Name:              rg.Name,
			},
			Location: ctx.GetLocation(rg.Location, cloud),
		}}, nil
	} else if cloud == common.AWS {
		return nil, nil
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil, nil
}

func (rg *Type) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{}
}

func GetResourceGroupName(name string, cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("azurerm_resource_group.%s.name", name)

	}
	validate.LogInternalError("cloud %s is not supported for resource groups ", cloud)
	return ""
}

func GetDefaultResourceGroupId() hcl.Expression {
	name, err := hclutil.StringToHclExpression("${resource_type}-rg")
	if err != nil {
		validate.LogInternalError("error setting default rg name: %s", err.Error())
	}
	return name
}

func DefaultResourceGroup(id string) *Type {
	return &Type{
		ResourceId: id,
		Name:       id,
	}
}

func (rg *Type) GetResourceId() string {
	return rg.ResourceId
}

func (rg *Type) GetLocation(cloud common.CloudProvider, ctx resources.MultyContext) string {
	return ctx.GetLocation(rg.Location, cloud)
}

func (rg *Type) Validate(ctx resources.MultyContext, cloud common.CloudProvider) []validate.ValidationError {
	return nil
}

func (rg *Type) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return "", nil
	case common.AZURE:
		return "AzureResourceName", nil
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return "", nil
}

func (rg *Type) GetDependencies(ctx resources.MultyContext) []resources.CloudSpecificResource {
	return nil
}
