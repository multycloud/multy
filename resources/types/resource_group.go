package types

import (
	"fmt"
	commonpb "github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/validate"
	"regexp"
)

var resourceGroupMetadata = resources.ResourceMetadata[*resourcespb.ResourceGroupArgs, *ResourceGroup, *resourcespb.ResourceGroupResource]{
	CreateFunc:        CreateResourceGroup,
	UpdateFunc:        UpdateResourceGroup,
	ReadFromStateFunc: ResourceGroupFromState,
	ExportFunc: func(r *ResourceGroup, others *resources.Resources) (*resourcespb.ResourceGroupArgs, bool, error) {
		if len(getAllDependentResources(r, others)) == 0 {
			return nil, false, nil
		}

		return &resourcespb.ResourceGroupArgs{
			CommonParameters: &commonpb.ResourceCommonArgs{
				Location:      r.Location,
				CloudProvider: r.Cloud,
			},
			Name: r.Name,
		}, true, nil
	},
	ImportFunc:      ImportResourceGroup,
	AbbreviatedName: "rg",
}

// https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
type ResourceGroup struct {
	ResourceId string
	Name       string `hcl:"name"`
	App        string `hcl:"app"`
	Location   commonpb.Location
	Cloud      commonpb.CloudProvider

	ImplictlyCreated bool
}

func (r *ResourceGroup) GetMetadata() resources.ResourceMetadataInterface {
	return &resourceGroupMetadata
}

func CreateResourceGroup(resourceId string, args *resourcespb.ResourceGroupArgs, others *resources.Resources) (*ResourceGroup, error) {
	return ImportResourceGroup(resourceId, args, others)
}

func UpdateResourceGroup(resource *ResourceGroup, vn *resourcespb.ResourceGroupArgs, others *resources.Resources) error {
	return fmt.Errorf("can't update resource group")
}

func ResourceGroupFromState(resource *ResourceGroup, _ *output.TfState) (*resourcespb.ResourceGroupResource, error) {
	return &resourcespb.ResourceGroupResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:    resource.ResourceId,
			Location:      resource.Location,
			CloudProvider: resource.Cloud,
			NeedsUpdate:   false,
		},
		Name: resource.Name,
	}, nil
}

func ImportResourceGroup(_ string, args *resourcespb.ResourceGroupArgs, _ *resources.Resources) (*ResourceGroup, error) {
	return &ResourceGroup{
		ResourceId:       args.Name,
		Name:             args.Name,
		Location:         args.GetCommonParameters().GetLocation(),
		Cloud:            args.GetCommonParameters().GetCloudProvider(),
		ImplictlyCreated: false,
	}, nil
}

type AzureResourceGroup struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_resource_group"`
	Location           string `hcl:"location"`
}

const AzureResourceName = "azurerm_resource_group"

func NewRgFromParent(resourceType string, parentResourceGroupId string, r *resources.Resources, location commonpb.Location, cloud commonpb.CloudProvider) string {
	var rgId string
	if rg, exists, err := resources.GetOptional[*ResourceGroup]("", r, parentResourceGroupId); exists && err == nil {
		if matches := regexp.MustCompile("\\w+-(\\w+)-rg").FindStringSubmatch(rg.Name); len(matches) >= 2 {
			rgId = getDefaultResourceGroupIdString(resourceType, matches[1])
			if !rgNameExists(r, rgId) {
				r.Add(NewResourceGroup(rgId, location, cloud))
			}
			return rgId
		}
	}
	return NewRg(resourceType, r, location, cloud)
}

func NewRg(resourceType string, r *resources.Resources, location commonpb.Location, cloud commonpb.CloudProvider) string {
	var rgName string
	for ok := false; !ok; ok = !rgNameExists(r, rgName) {
		rgName = getDefaultResourceGroupIdString(resourceType, common.RandomString(4))
	}
	r.Add(NewResourceGroup(rgName, location, cloud))
	return rgName
}

func rgNameExists(r *resources.Resources, rgName string) bool {
	for _, resource := range r.GetAll() {
		if rg, ok := resource.(*ResourceGroup); ok {
			if rg.Name == rgName {
				return true
			}
		}
	}
	return false
}

func NewResourceGroup(name string, location commonpb.Location, cloud commonpb.CloudProvider) *ResourceGroup {
	return &ResourceGroup{
		ResourceId: name,
		Name:       name,
		Location:   location,
		Cloud:      cloud,

		ImplictlyCreated: true,
	}
}

func (rg *ResourceGroup) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	allDeps := getAllDependentResources(rg, ctx.Resources)
	if len(allDeps) == 0 {
		// if no resources are in this group, just don't output anything
		return nil, nil
	}
	for _, r := range allDeps {
		ctx.Resources.AddDependency(rg.ResourceId, r)
	}

	if rg.GetCloud() == common.AZURE {
		return []output.TfBlock{AzureResourceGroup{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: rg.ResourceId, ResourceName: AzureResourceName},
				Name:              rg.Name,
			},
			Location: rg.GetCloudSpecificLocation(),
		}}, nil
	} else if rg.GetCloud() == common.AWS {
		return nil, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", rg.GetCloud())
}

func GetResourceGroupName(name string) string {
	return fmt.Sprintf("azurerm_resource_group.%s.name", name)
}
func getDefaultResourceGroupIdString(resourceType string, groupId string) string {
	return fmt.Sprintf("%s-%s-rg", resourceType, groupId)
}

func (rg *ResourceGroup) GetResourceId() string {
	return rg.ResourceId
}

func (rg *ResourceGroup) GetCloudSpecificLocation() string {
	if result, err := common.GetCloudLocation(rg.Location, rg.GetCloud()); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}

func (rg *ResourceGroup) Validate(ctx resources.MultyContext) []validate.ValidationError {
	if _, err := common.GetCloudLocation(rg.Location, rg.GetCloud()); err != nil {
		return []validate.ValidationError{
			{
				ErrorMessage: err.Error(),
				ResourceId:   rg.ResourceId,
				FieldName:    "location",
			},
		}
	}
	return nil
}

func (rg *ResourceGroup) GetMainResourceName() (string, error) {
	switch rg.GetCloud() {
	case common.AWS:
		return "", nil
	case common.AZURE:
		return "AzureResourceName", nil
	default:
		return "", fmt.Errorf("unknown cloud %s", rg.GetCloud())
	}
}

func (rg *ResourceGroup) GetCloud() commonpb.CloudProvider {
	return rg.Cloud
}

func (rg *ResourceGroup) GetCommonArgs() any {
	return nil
}

func getAllDependentResources(r *ResourceGroup, others *resources.Resources) (res []string) {
	for _, other := range others.GetAll() {
		if wrg, ok := other.(interface{ GetResourceGroupId() string }); ok {
			if wrg.GetResourceGroupId() == r.GetResourceId() {
				res = append(res, other.GetResourceId())
			}
		}
	}

	return res
}
