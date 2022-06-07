package types

import (
	"fmt"
	commonpb "github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
	"regexp"
)

func (r *ResourceGroup) Create(_ string, args *resourcespb.ResourceGroupArgs, _ *resources.Resources) error {
	return ImportResourceGroup(r, args)
}

func (r *ResourceGroup) Update(args *resourcespb.ResourceGroupArgs, _ *resources.Resources) error {
	return fmt.Errorf("updates to resource groups are not supported")
}

func (r *ResourceGroup) Import(_ string, args *resourcespb.ResourceGroupArgs, _ *resources.Resources) error {
	return ImportResourceGroup(r, args)
}

func (r *ResourceGroup) Export(others *resources.Resources) (*resourcespb.ResourceGroupArgs, bool, error) {
	if len(r.GetAllDependentResources(others)) == 0 {
		return nil, false, nil
	}
	return &resourcespb.ResourceGroupArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      r.Args.CommonParameters.Location,
			CloudProvider: r.Args.CommonParameters.CloudProvider,
		},
	}, true, nil
}

// https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
type ResourceGroup struct {
	resources.ResourceWithId[*resourcespb.ResourceGroupArgs]

	ImplictlyCreated bool
}

func ImportResourceGroup(rg *ResourceGroup, args *resourcespb.ResourceGroupArgs) error {
	rg.ResourceWithId = resources.NewResource(args.Name, args)
	return nil
}

type AzureResourceGroup struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_resource_group"`
	Location           string `hcl:"location"`
}

const AzureResourceName = "azurerm_resource_group"

func NewRgFromParent(resourceType string, parentResourceGroupId string, r *resources.Resources, location commonpb.Location, cloud commonpb.CloudProvider) (string, error) {
	var rgId string
	if rg, exists, err := resources.GetOptional[*ResourceGroup]("", r, parentResourceGroupId); exists && err == nil {
		if matches := regexp.MustCompile("\\w+-(\\w+)-rg").FindStringSubmatch(rg.ResourceId); len(matches) >= 2 {
			rgId = getDefaultResourceGroupIdString(resourceType, matches[1])
			if !rgNameExists(r, rgId) {
				return rgId, r.Add(NewResourceGroup(rgId, location, cloud))
			}
			return rgId, nil
		}
	}
	return NewRg(resourceType, r, location, cloud)
}

func NewRg(resourceType string, r *resources.Resources, location commonpb.Location, cloud commonpb.CloudProvider) (string, error) {
	var rgName string
	for ok := false; !ok; ok = !rgNameExists(r, rgName) {
		rgName = getDefaultResourceGroupIdString(resourceType, common.RandomString(4))
	}
	return rgName, r.Add(NewResourceGroup(rgName, location, cloud))
}

func rgNameExists(r *resources.Resources, rgName string) bool {
	for _, resource := range r.GetAll() {
		if rg, ok := resource.(*ResourceGroup); ok {
			if rg.ResourceId == rgName {
				return true
			}
		}
	}
	return false
}

func NewResourceGroup(name string, location commonpb.Location, cloud commonpb.CloudProvider) *ResourceGroup {
	return &ResourceGroup{
		ResourceWithId: resources.ResourceWithId[*resourcespb.ResourceGroupArgs]{
			ResourceId: name,
			Args: &resourcespb.ResourceGroupArgs{
				CommonParameters: &commonpb.ResourceCommonArgs{
					Location:      location,
					CloudProvider: cloud,
				},
				Name: name,
			},
		},
		ImplictlyCreated: true,
	}
}

func GetResourceGroupName(name string) string {
	return fmt.Sprintf("azurerm_resource_group.%s.name", name)
}
func getDefaultResourceGroupIdString(resourceType string, groupId string) string {
	return fmt.Sprintf("%s-%s-rg", resourceType, groupId)
}

func (rg *ResourceGroup) Validate(ctx resources.MultyContext) []validate.ValidationError {
	return nil
}

func (r *ResourceGroup) GetAllDependentResources(others *resources.Resources) (res []string) {
	for _, other := range others.GetAll() {
		if wrg, ok := other.(interface{ GetResourceGroupId() string }); ok {
			if wrg.GetResourceGroupId() == r.GetResourceId() {
				res = append(res, other.GetResourceId())
			}
		}
	}

	return res
}
