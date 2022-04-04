package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/validate"
)

type Resources map[string]Resource

func GetCloudSpecificResourceId(r Resource, cloud commonpb.CloudProvider) string {
	return GetResourceIdForCloud(r.GetResourceId(), cloud)
}

func GetResourceIdForCloud(resourceId string, cloud commonpb.CloudProvider) string {
	return fmt.Sprintf("%s.%s", cloud, resourceId)
}

type CloudSpecificResource struct {
	Cloud             commonpb.CloudProvider
	Resource          Resource
	ImplicitlyCreated bool
}

func (c *CloudSpecificResource) GetResourceId() string {
	return GetCloudSpecificResourceId(c.Resource, c.Cloud)
}

func (c *CloudSpecificResource) GetLocation(ctx MultyContext) string {
	return c.Resource.GetCloudSpecificLocation()
}

func (c *CloudSpecificResource) Translate(ctx MultyContext) ([]output.TfBlock, error) {
	return c.Resource.Translate(ctx)
}

type Resource interface {
	Translate(ctx MultyContext) ([]output.TfBlock, error)

	GetResourceId() string

	GetCloudSpecificLocation() string

	Validate(ctx MultyContext) []validate.ValidationError

	GetMainResourceName() (string, error)

	GetCloud() commonpb.CloudProvider
}

func (c *CloudSpecificResource) GetMainOutputId() (string, error) {
	return GetMainOutputId(c.Resource)
}

func GetMainOutputId(r Resource) (string, error) {
	name, err := r.GetMainResourceName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("${%s.%s.id}", name, r.GetResourceId()), nil
}
