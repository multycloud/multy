package rg

import (
	"fmt"
	"multy-go/hclutil"
	"multy-go/resources"
	"multy-go/resources/common"
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
	common.AzResource `hcl:",squash"`
	Location          string `hcl:"location"`
}

const AzureResourceName = "azurerm_resource_group"

func (rg *Type) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	if cloud == common.AZURE {
		return []interface{}{ResourceGroup{
			AzResource: common.AzResource{
				ResourceName: AzureResourceName,
				ResourceId:   rg.ResourceId,
				Name:         rg.Name,
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
