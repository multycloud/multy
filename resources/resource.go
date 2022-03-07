package resources

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/validate"
)

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

func (c *CloudSpecificResource) Translate(ctx MultyContext) []output.TfBlock {
	return c.Resource.Translate(c.Cloud, ctx)
}

type Resource interface {
	Translate(cloud common.CloudProvider, ctx MultyContext) []output.TfBlock
	// GetOutputValues returns values that should be passed around when parsing the remainder of the config file.
	GetOutputValues(cloud common.CloudProvider) map[string]cty.Value

	GetResourceId() string

	GetLocation(cloud common.CloudProvider, ctx MultyContext) string

	Validate(ctx MultyContext, cloud common.CloudProvider) []validate.ValidationError

	GetMainResourceName(cloud common.CloudProvider) string

	GetDependencies(ctx MultyContext) []CloudSpecificResource
}

func getTfResourceId(resourceId string, cloud common.CloudProvider) string {
	return fmt.Sprintf("%s_%s", resourceId, cloud)

}

func (c *CloudSpecificResource) GetMainOutputId() string {
	return GetMainOutputId(c.Resource, c.Cloud)
}

func GetMainOutputId(r Resource, cloud common.CloudProvider) string {
	return fmt.Sprintf("${%s.%s.id}", r.GetMainResourceName(cloud), getTfResourceId(r.GetResourceId(), cloud))
}
