package resources

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/validate"

	"github.com/zclconf/go-cty/cty"
)

type CommonResourceParams struct {
	ResourceId      string
	ResourceGroupId string
	Location        string            `hcl:"location,optional"`
	Clouds          []string          `hcl:"clouds,optional"`
	RgVars          map[string]string `hcl:"rg_vars,optional"`
	*validate.ResourceValidationInfo
}

type CommonResourceOutputs struct {
	ResourceId      string
	ResourceGroupId string
	Location        string
	Clouds          []string
	RgVars          map[string]string
	Name            string
}

type MultyContext struct {
	Resources map[string]CloudSpecificResource
	Location  string
}

func (ctx *MultyContext) GetResource(id string) (*CloudSpecificResource, error) {
	if r, ok := ctx.Resources[id]; ok {
		return &r, nil
	}
	return nil, fmt.Errorf("resource %s not found", id)
}

func (ctx *MultyContext) GetLocationFromCommonParams(commonParams *CommonResourceParams, cloud common.CloudProvider) string {
	location := ctx.Location
	if commonParams.Location != "" {
		location = commonParams.Location
	}

	if result, err := common.GetCloudLocation(location, cloud); err != nil {
		if location == commonParams.Location {
			commonParams.LogFatal(commonParams.ResourceId, "location", err.Error())
		} else {
			// TODO: throw a user error when validating global config
			validate.LogInternalError(err.Error())
		}
		return ""
	} else {
		return result
	}
}

func (ctx *MultyContext) GetLocation(specifiedLocation string, cloud common.CloudProvider) string {
	location := ctx.Location
	if specifiedLocation != "" {
		location = specifiedLocation
	}

	if result, err := common.GetCloudLocation(location, cloud); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}

func (c *CommonResourceParams) GetResourceId() string {
	return c.ResourceId
}

func (c *CommonResourceParams) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	return map[string]cty.Value{}
}

func (c *CommonResourceParams) GetTfResourceId(cloud common.CloudProvider) string {
	return fmt.Sprintf("%s_%s", c.GetResourceId(), cloud)
}

func (c *CommonResourceParams) GetLocation(cloud common.CloudProvider, ctx MultyContext) string {
	return ctx.GetLocationFromCommonParams(c, cloud)
}

func GetCloudSpecificResourceId(r Resource, cloud common.CloudProvider) string {
	return fmt.Sprintf("%s.%s", cloud, r.GetResourceId())
}

type CloudSpecificResource struct {
	Cloud             common.CloudProvider
	Resource          Resource
	ImplicitlyCreated bool
}

func (c *CloudSpecificResource) GetResourceId() string {
	return GetCloudSpecificResourceId(c.Resource, c.Cloud)
}

func (c *CloudSpecificResource) GetLocation(ctx MultyContext) string {
	return c.Resource.GetLocation(c.Cloud, ctx)
}

func (c *CloudSpecificResource) Translate(ctx MultyContext) []interface{} {
	return c.Resource.Translate(c.Cloud, ctx)
}

type Resource interface {
	Translate(cloud common.CloudProvider, ctx MultyContext) []interface{}
	// GetOutputValues returns values that should be passed around when parsing the remainder of the config file.
	GetOutputValues(cloud common.CloudProvider) map[string]cty.Value

	GetResourceId() string

	GetLocation(cloud common.CloudProvider, ctx MultyContext) string

	Validate(ctx MultyContext)
}
