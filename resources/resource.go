package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Resources struct {
	ResourceMap  map[string]Resource
	dependencies map[string][]string
	resources    []Resource
}

func NewResources() *Resources {
	return &Resources{
		ResourceMap:  map[string]Resource{},
		dependencies: map[string][]string{},
	}
}

func (r *Resources) Add(resource Resource) error {
	if _, ok := r.ResourceMap[resource.GetResourceId()]; ok {
		return fmt.Errorf("attempted to add a resource with an already existing id (%s), this should never happen", resource.GetResourceId())
	}
	r.ResourceMap[resource.GetResourceId()] = resource
	r.resources = append(r.resources, resource)
	return nil
}

func (r *Resources) Delete(resourceId string) {
	for i, resource := range r.resources {
		if resource.GetResourceId() == resourceId {
			r.resources = slices.Delete(r.resources, i, i+1)
			break
		}
	}
	delete(r.ResourceMap, resourceId)
}

func (r *Resources) GetAll() []Resource {
	return r.resources
}

// Get finds the resource with the given id and adds a dependency between dependentResourceId and id.
func Get[T Resource](dependentResourceId string, resources *Resources, id string) (T, error) {
	// TODO return better error on empty ID
	item, exists, err := GetOptional[T](dependentResourceId, resources, id)
	if err != nil {
		return item, err
	}
	if !exists {
		return item, errors.ResourceNotFound(id)
	}

	return item, nil
}

func GetOptional[T Resource](dependentResourceId string, resources *Resources, id string) (T, bool, error) {
	if r, ok := resources.ResourceMap[id]; ok {
		if _, okType := r.(T); !okType {
			return *new(T), true, fmt.Errorf("resource with id %s is of the wrong type", id)
		}
		resources.AddDependency(dependentResourceId, id)
		return r.(T), true, nil
	}

	return *new(T), false, nil
}

func (r *Resources) AddDependency(dependentResourceId string, id string) {
	r.dependencies[dependentResourceId] = append(r.dependencies[dependentResourceId], id)
}

type MultyResourceGroup struct {
	GroupId   string
	Resources []Resource
}

func generateUniqueGroupId(existingGroups []string) (groupId string) {
	for groupId = common.RandomString(4); slices.Contains(existingGroups, groupId); groupId = common.RandomString(4) {
	}

	return
}

func (r *Resources) GetMultyResourceGroups(existingGroupsByResource map[string]string) map[Resource]*MultyResourceGroup {
	groups := map[Resource]*MultyResourceGroup{}
	// creates 1 group per resource
	for _, resource := range util.GetSortedMapValues(r.ResourceMap) {
		if _, ok := groups[resource]; !ok {
			var groupId string
			if existingGroupId, ok := existingGroupsByResource[resource.GetResourceId()]; ok {
				groupId = existingGroupId
			} else {
				groupId = generateUniqueGroupId(maps.Values(existingGroupsByResource))
			}
			groups[resource] = &MultyResourceGroup{
				GroupId:   groupId,
				Resources: []Resource{resource},
			}
		}
	}
	// merge all groups
	for _, resource := range util.GetSortedMapValues(r.ResourceMap) {
		for _, dep := range r.dependencies[resource.GetResourceId()] {
			// prefer to keep existing groups
			if _, hasGroup := existingGroupsByResource[resource.GetResourceId()]; hasGroup {
				mergeGroups(groups, resource, r.ResourceMap[dep])
			} else {
				mergeGroups(groups, r.ResourceMap[dep], resource)
			}
		}
	}

	return groups
}

func mergeGroups(all map[Resource]*MultyResourceGroup, res1 Resource, res2 Resource) {
	group := all[res1]
	if group2, ok := all[res2]; !ok {
		group.Resources = append(group.Resources, res2)
		all[res2] = group
	} else if group != group2 {
		for _, group2Resource := range group2.Resources {
			group.Resources = append(group.Resources, group2Resource)
			all[group2Resource] = group
		}
	}
}

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

	GetCommonArgs() any

	GetMetadata() ResourceMetadataInterface
}

func (c *CloudSpecificResource) GetMainOutputId() (string, error) {
	return GetMainOutputId(c.Resource)
}

func GetMainOutputId(r Resource) (string, error) {
	name, err := r.GetMainResourceName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s.id", name, r.GetResourceId()), nil
}

func GetMainOutputRef(r Resource) (string, error) {
	name, err := r.GetMainResourceName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", name, r.GetResourceId()), nil
}
