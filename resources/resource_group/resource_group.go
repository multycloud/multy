package rg

import (
	"fmt"
	commonpb "github.com/multycloud/multy/api/proto/commonpb"
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
	Cloud      commonpb.CloudProvider
}

type ResourceGroup struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_resource_group"`
	Location           string `hcl:"location"`
}

const AzureResourceName = "azurerm_resource_group"

func (rg *Type) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if rg.GetCloud() == common.AZURE {
		return []output.TfBlock{ResourceGroup{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: rg.ResourceId},
				Name:              rg.Name,
			},
			Location: rg.GetCloudSpecificLocation(),
		}}, nil
	} else if rg.GetCloud() == common.AWS {
		return nil, nil
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", rg.GetCloud())
	return nil, nil
}

func (rg *Type) GetOutputValues(cloud commonpb.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{}
}

func GetResourceGroupName(name string) string {
	return fmt.Sprintf("azurerm_resource_group.%s.name", name)
}
func GetDefaultResourceGroupIdString(resourceType string) string {
	return fmt.Sprintf("%s-rg", resourceType)
}

func (rg *Type) GetResourceId() string {
	if rg.Cloud != commonpb.CloudProvider_AZURE {
		// this should never be used, as the translate will not return anything
		return "_"
	}
	return rg.ResourceId
}

func (rg *Type) GetCloudSpecificLocation() string {
	if result, err := common.GetCloudLocationPb(rg.Location, rg.GetCloud()); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}

func (rg *Type) Validate(ctx resources.MultyContext) []validate.ValidationError {
	return nil
}

func (rg *Type) GetMainResourceName() (string, error) {
	switch rg.GetCloud() {
	case common.AWS:
		return "", nil
	case common.AZURE:
		return "AzureResourceName", nil
	default:
		validate.LogInternalError("unknown cloud %s", rg.GetCloud())
	}
	return "", nil
}

func (rg *Type) GetDependencies(ctx resources.MultyContext) []resources.CloudSpecificResource {
	return nil
}

func (rg *Type) GetCloud() commonpb.CloudProvider {
	return rg.Cloud
}
