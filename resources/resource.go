package resources

import (
	"fmt"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
)

func GetCloudSpecificResourceId(r Resource, cloud common.CloudProvider) string {
	return GetResourceIdForCloud(r.GetResourceId(), cloud)
}

func GetResourceIdForCloud(resourceId string, cloud common.CloudProvider) string {
	return fmt.Sprintf("%s.%s", cloud, resourceId)
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

func (c *CloudSpecificResource) Translate(ctx MultyContext) ([]output.TfBlock, error) {
	return c.Resource.Translate(c.Cloud, ctx)
}

type Resource interface {
	Translate(cloud common.CloudProvider, ctx MultyContext) ([]output.TfBlock, error)
	// GetOutputValues returns values that should be passed around when parsing the remainder of the config file.
	GetOutputValues(cloud common.CloudProvider) map[string]cty.Value

	GetResourceId() string

	GetLocation(cloud common.CloudProvider, ctx MultyContext) string

	Validate(ctx MultyContext, cloud common.CloudProvider) []validate.ValidationError

	GetMainResourceName(cloud common.CloudProvider) (string, error)

	GetDependencies(ctx MultyContext) []CloudSpecificResource
}

func (c *CloudSpecificResource) GetMainOutputId() (string, error) {
	return GetMainOutputId(c.Resource, c.Cloud)
}

func GetMainOutputId(r Resource, cloud common.CloudProvider) (string, error) {
	name, err := r.GetMainResourceName(cloud)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("${%s.%s.id}", name, util.GetTfResourceId(r.GetResourceId(), string(cloud))), nil
}
