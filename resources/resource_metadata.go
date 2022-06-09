package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/resources/output"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
	"sort"
)

type ResourceMetadatas map[proto.Message]ResourceMetadataInterface

type ResourceTranslator[OutT proto.Message] interface {
	Translate(MultyContext) ([]output.TfBlock, error)
	FromState(state *output.TfState) (OutT, error)
}

type ResourceExporter[ArgsT proto.Message] interface {
	Create(resourceId string, args ArgsT, others *Resources) error
	Update(args ArgsT, others *Resources) error
	Import(resourceId string, args ArgsT, others *Resources) error
	Export(others *Resources) (ArgsT, bool, error)

	Resource
}

type ResourceMetadata[ArgsT proto.Message, R ResourceExporter[ArgsT], OutT proto.Message] struct {
	Translators map[commonpb.CloudProvider]func(R) ResourceTranslator[OutT]

	AbbreviatedName string
	// Used for error messages
	ResourceType string
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Create(resourceId string, args proto.Message, resources *Resources) (Resource, error) {
	r := reflect.New(reflect.TypeOf(*new(R)).Elem()).Interface().(R)
	err := r.Create(resourceId, args.(ArgsT), resources)
	return r, err
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Update(resource Resource, args proto.Message, resources *Resources) error {
	r := resource.(ResourceExporter[ArgsT])
	return r.Update(args.(ArgsT), resources)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) ReadFromState(resource Resource, state *output.TfState) (proto.Message, error) {
	out, err := m.Translators[resource.GetCloud()](resource.(R)).FromState(state)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Export(resource Resource, resources *Resources) (proto.Message, bool, error) {
	r := resource.(ResourceExporter[ArgsT])
	return r.Export(resources)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Import(resourceId string, args proto.Message, resources *Resources) (Resource, error) {
	// TODO: do this without reflection
	r := reflect.New(reflect.TypeOf(*new(R)).Elem()).Interface().(R)
	err := r.Import(resourceId, args.(ArgsT), resources)
	return r, err
}

func (m *ResourceMetadata[ArgsT, R, OutT]) GetCloudSpecificResource(r Resource) (CloudSpecificResourceTranslator, error) {
	if _, ok := m.Translators[r.GetCloud()]; !ok {
		return nil, fmt.Errorf("resource of type %s not supported in cloud %s", m.ResourceType, r.GetCloud())
	}
	return m.Translators[r.GetCloud()](r.(R)).(CloudSpecificResourceTranslator), nil
}

func (m *ResourceMetadata[ArgsT, R, OutT]) GetAbbreviatedName() string {
	return m.AbbreviatedName
}

type ResourceMetadataInterface interface {
	Create(string, proto.Message, *Resources) (Resource, error)
	Update(Resource, proto.Message, *Resources) error
	ReadFromState(Resource, *output.TfState) (proto.Message, error)

	Export(Resource, *Resources) (proto.Message, bool, error)
	Import(string, proto.Message, *Resources) (Resource, error)
	GetCloudSpecificResource(r Resource) (CloudSpecificResourceTranslator, error)

	GetAbbreviatedName() string
}

type MultyConfig struct {
	Resources                  *Resources
	c                          *configpb.Config
	groupsByResourceId         map[Resource]*MultyResourceGroup
	affectedResourcesById      map[string][]string
	metadatas                  ResourceMetadatas
	affectedResourcesByGroupId map[string][]string
}

func (c *MultyConfig) GetUserId() string {
	return c.c.UserId
}

func LoadConfig(c *configpb.Config, metadatas ResourceMetadatas) (*MultyConfig, error) {
	multyc := &MultyConfig{
		c:         c,
		metadatas: metadatas,
	}
	res := NewResources()

	for _, r := range c.Resources {
		conv, err := multyc.metadatas.getConverter(r.ResourceArgs.ResourceArgs.MessageName())
		if err != nil {
			return multyc, err
		}
		err = addMultyResource(r, res, conv)
		if err != nil {
			return multyc, err
		}
	}

	multyc.Resources = res

	multyc.affectedResourcesById = map[string][]string{}
	multyc.affectedResourcesByGroupId = map[string][]string{}
	for _, r := range c.Resources {
		if r.GetDeployedResourceGroup() != nil {
			multyc.affectedResourcesById[r.ResourceId] = r.GetDeployedResourceGroup().GetDeployedResource()
			multyc.affectedResourcesByGroupId[r.GetDeployedResourceGroup().GetGroupId()] = r.GetDeployedResourceGroup().GetDeployedResource()
		}
	}

	return multyc, nil
}

func (c *MultyConfig) GetOriginalConfig(metadatas ResourceMetadatas) (*MultyConfig, error) {
	return LoadConfig(c.c, metadatas)
}

func addMultyResource(r *configpb.Resource, res *Resources, metadata ResourceMetadataInterface) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}

	translatedResource, err := metadata.Import(r.ResourceId, m, res)
	if err != nil {
		return err
	}

	return res.Add(translatedResource)
}

func (c *MultyConfig) CreateResource(args proto.Message) (Resource, error) {
	conv, err := c.metadatas.getConverter(proto.MessageName(args))
	if err != nil {
		return nil, err
	}
	c.c.ResourceCounter += 1
	resourceId := fmt.Sprintf("multy_%s_r%d", conv.GetAbbreviatedName(), c.c.ResourceCounter)
	r, err := conv.Create(resourceId, args, c.Resources)
	if err != nil {
		return nil, err
	}
	err = c.Resources.Add(r)
	return r, err
}

func (c *MultyConfig) UpdateResource(resourceId string, args proto.Message) (Resource, error) {
	conv, err := c.metadatas.getConverter(proto.MessageName(args))
	if err != nil {
		return nil, err
	}
	r, exists := c.Resources.ResourceMap[resourceId]
	if !exists {
		return nil, errors.ResourceNotFound(resourceId)
	}
	err = conv.Update(r, args, c.Resources)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *MultyConfig) DeleteResource(resourceId string) (Resource, error) {
	if r, exists := c.Resources.ResourceMap[resourceId]; exists {
		c.Resources.Delete(resourceId)
		return r, nil
	} else {
		return nil, errors.ResourceNotFound(resourceId)
	}

}

func (c *MultyConfig) UpdateMultyResourceGroups() {
	groupIdsByResourceIds := map[string]string{}
	for _, r := range c.c.Resources {
		if r.DeployedResourceGroup != nil {
			groupIdsByResourceIds[r.ResourceId] = r.DeployedResourceGroup.GroupId
		}
	}
	c.groupsByResourceId = c.Resources.GetMultyResourceGroups(groupIdsByResourceIds)
}

func (c *MultyConfig) UpdateDeployedResourceList(deployedResources map[Resource][]string) {
	affectedResources := map[string][]string{}

	for r, deployedResource := range deployedResources {
		groupId := c.groupsByResourceId[r]
		for _, d := range deployedResource {
			if !slices.Contains(affectedResources[groupId.GroupId], d) {
				affectedResources[groupId.GroupId] = append(affectedResources[groupId.GroupId], d)
			}
		}
		sort.Strings(affectedResources[groupId.GroupId])
	}
	c.affectedResourcesByGroupId = affectedResources

	affectedResourcesById := map[string][]string{}
	for _, r := range c.Resources.GetAll() {
		affectedResourcesById[r.GetResourceId()] = affectedResources[c.groupsByResourceId[r].GroupId]
	}
	c.affectedResourcesById = affectedResourcesById
}

func (c *MultyConfig) GetAffectedResources(resourceId string) []string {
	return c.affectedResourcesById[resourceId]
}

func (c *MultyConfig) ExportConfig() (*configpb.Config, error) {
	result := &configpb.Config{
		UserId:          c.GetUserId(),
		ResourceCounter: c.c.ResourceCounter,
	}

	groups := map[string]*configpb.DeployedResourceGroup{}
	for groupId, affectedResources := range c.affectedResourcesByGroupId {
		groups[groupId] = &configpb.DeployedResourceGroup{
			GroupId:          groupId,
			DeployedResource: affectedResources,
		}
	}

	for _, r := range c.Resources.GetAll() {
		m, err := r.GetMetadata(c.metadatas)
		if err != nil {
			return nil, err
		}
		msg, export, err := m.Export(r, c.Resources)
		if err != nil {
			return nil, err
		}
		if !export {
			continue
		}
		a, err := anypb.New(msg)
		if err != nil {
			return nil, err
		}
		resource := configpb.Resource{
			ResourceId:            r.GetResourceId(),
			ResourceArgs:          &configpb.ResourceArgs{ResourceArgs: a},
			DeployedResourceGroup: groups[c.groupsByResourceId[r].GroupId],
		}
		result.Resources = append(result.Resources, &resource)
	}

	return result, nil
}

func (m ResourceMetadatas) getConverter(name protoreflect.FullName) (ResourceMetadataInterface, error) {
	for messageType, conv := range m {
		if name == proto.MessageName(messageType) {
			return conv, nil
		}
	}
	return nil, fmt.Errorf("unknown resource type %s", name)
}
