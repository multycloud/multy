package rg

import (
	"fmt"
	"multy-go/hclutil"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/validate"

	"github.com/hashicorp/hcl/v2"
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

func (rg *Type) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AZURE {
		return []output.TfBlock{ResourceGroup{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: rg.ResourceId},
				Name:              rg.Name,
			},
			Location: ctx.GetLocation(rg.Location, cloud),
		}}
	} else if cloud == common.AWS {
		return nil
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
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

func (rg *Type) Validate(ctx resources.MultyContext) {
	return
}

func (rg *Type) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return ""
	case common.AZURE:
		return "AzureResourceName"
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}

func (rg *Type) GetDependencies(ctx resources.MultyContext) []resources.CloudSpecificResource {
	return nil
}
