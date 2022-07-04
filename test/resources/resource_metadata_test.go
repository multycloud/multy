package resources_test

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/resources/types/metadata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

func TestResourceMetadataImport(t *testing.T) {
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			ResourceGroupId: "test-rg1",
			Location:        commonpb.Location_EU_WEST_1,
			CloudProvider:   commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}
	c := &configpb.Config{
		UserId: "test-user",
		Resources: []*configpb.Resource{
			newResource(t, "resource1", args),
		},
	}

	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	assert.Equal(t, "test-user", c.GetUserId())
	assert.Equal(t, 1, len(config.Resources.GetAll()))
	virtualNetwork, ok := config.Resources.ResourceMap["resource1"].(*types.VirtualNetwork)
	if !ok {
		t.Fatalf("resource should of type virtual network, but is %t", config.Resources.ResourceMap["resource1"])
	}
	assert.Equal(t, "resource1", virtualNetwork.GetResourceId())
	assert.Equal(t, args.String(), virtualNetwork.Args.String())
}

func TestResourceMetadataImport_returnsAffectedResources(t *testing.T) {
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			ResourceGroupId: "test-rg1",
			Location:        commonpb.Location_EU_WEST_1,
			CloudProvider:   commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}
	group := &configpb.DeployedResourceGroup{
		GroupId:          "group1",
		DeployedResource: []string{"aws_vpc.resource1", "aws_vpc.resource2"},
	}
	r1 := newResource(t, "resource1", args)
	r2 := newResource(t, "resource2", args)
	r1.DeployedResourceGroup = group
	r2.DeployedResourceGroup = group
	c := &configpb.Config{
		UserId: "test-user",
		Resources: []*configpb.Resource{
			r1, r2,
		},
	}

	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	assert.Equal(t, "test-user", c.GetUserId())
	assert.Contains(t, config.GetAffectedResources("resource1"), "aws_vpc.resource1")
	assert.Contains(t, config.GetAffectedResources("resource2"), "aws_vpc.resource1")
	assert.Contains(t, config.GetAffectedResources("resource1"), "aws_vpc.resource2")
	assert.Contains(t, config.GetAffectedResources("resource2"), "aws_vpc.resource2")
}

func TestResourceMetadataCreate_createsResourceGroup(t *testing.T) {
	c := &configpb.Config{
		UserId: "test-user",
	}
	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_NORTH_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}

	r, err := config.CreateResource(proto.Clone(args))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}

	resourceId := r.GetResourceId()
	assert.Equal(t, 2, len(config.Resources.GetAll()))

	// assert virtual network, except resource group id
	virtualNetwork, ok := config.Resources.ResourceMap[resourceId].(*types.VirtualNetwork)
	if !ok {
		t.Fatalf("resource should of type virtual network, but is %t", config.Resources.ResourceMap[resourceId])
	}
	assert.Equal(t, resourceId, virtualNetwork.GetResourceId())
	actualArgsExceptRgId := proto.Clone(virtualNetwork.Args).(*resourcespb.VirtualNetworkArgs)
	actualArgsExceptRgId.CommonParameters.ResourceGroupId = ""
	assert.Equal(t, args.String(), actualArgsExceptRgId.String())

	// assert resource group
	rgId := virtualNetwork.Args.CommonParameters.ResourceGroupId
	rg, ok := config.Resources.ResourceMap[rgId].(*types.ResourceGroup)
	if !ok {
		t.Fatalf("resource should of type resource group, but is %t", config.Resources.ResourceMap[rgId])
	}
	assert.Equal(t, rgId, rg.Args.Name)
	assert.Equal(t, args.CommonParameters.Location, rg.Args.CommonParameters.Location)
	assert.Equal(t, args.CommonParameters.CloudProvider, rg.Args.CommonParameters.CloudProvider)
}

func TestResourceMetadataDelete_deletesResourceGroup(t *testing.T) {
	c := &configpb.Config{
		UserId: "test-user",
	}
	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_NORTH_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}

	r, err := config.CreateResource(proto.Clone(args))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}

	resourceId := r.GetResourceId()
	// assert we have 2 resources: virtual network and resource group
	assert.Equal(t, 2, len(config.Resources.GetAll()))

	_, err = config.DeleteResource(resourceId)
	if err != nil {
		t.Fatalf("unable to delete resource: %s", err)
	}

	exportedConfig, err := config.ExportConfig()
	if err != nil {
		t.Fatalf("unable to export resources: %s", err)
	}

	assert.Equal(t, 0, len(exportedConfig.Resources))
}

func TestResourceMetadataExport(t *testing.T) {
	c := &configpb.Config{
		UserId: "test-user",
	}
	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_NORTH_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}

	_, err = config.CreateResource(proto.Clone(args))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}

	numResources := len(config.Resources.GetAll())
	assert.Greater(t, numResources, 0)

	config.UpdateMultyResourceGroups()
	exportedConfig, err := config.ExportConfig()
	if err != nil {
		t.Fatalf("unable to export resources: %s", err)
	}

	assert.Equal(t, len(exportedConfig.Resources), numResources)
}

func TestResourceMetadataUpdateMultyResourceGroups(t *testing.T) {
	c := &configpb.Config{
		UserId: "test-user",
	}
	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_NORTH_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}

	vn1, err := config.CreateResource(proto.Clone(args))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}
	vn2, err := config.CreateResource(proto.Clone(args))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}

	subnetArgs1 := &resourcespb.SubnetArgs{
		Name:             "test-vn",
		CidrBlock:        "10.0.0.0/16",
		VirtualNetworkId: vn1.GetResourceId(),
	}
	subnetArgs2 := &resourcespb.SubnetArgs{
		Name:             "test-vn",
		CidrBlock:        "10.0.0.0/16",
		VirtualNetworkId: vn2.GetResourceId(),
	}
	subnet1, err := config.CreateResource(proto.Clone(subnetArgs1))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}
	subnet2, err := config.CreateResource(proto.Clone(subnetArgs2))
	if err != nil {
		t.Fatalf("unable to create resource: %s", err)
	}

	config.UpdateMultyResourceGroups()

	config.UpdateDeployedResourceList(map[string][]string{})
	affectedResourcesSubnet1 := config.GetAffectedResources(subnet1.GetResourceId())
	affectedResourcesVn1 := config.GetAffectedResources(vn1.GetResourceId())
	affectedResourcesSubnet2 := config.GetAffectedResources(subnet2.GetResourceId())
	affectedResourcesVn2 := config.GetAffectedResources(vn2.GetResourceId())

	// asserts that vns and subnets are in the same group
	assert.Equal(t, affectedResourcesSubnet1, affectedResourcesVn1)
	assert.Equal(t, affectedResourcesSubnet2, affectedResourcesVn2)
	// asserts that the two groups are not merged together
	for _, affectedResource := range affectedResourcesSubnet1 {
		assert.NotContains(t, affectedResourcesSubnet2, affectedResource)
	}
	for _, affectedResource := range affectedResourcesSubnet2 {
		assert.NotContains(t, affectedResourcesSubnet1, affectedResource)
	}
}

func TestResourceMetadataUpdate(t *testing.T) {
	args := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			ResourceGroupId: "test-rg1",
			Location:        commonpb.Location_EU_WEST_1,
			CloudProvider:   commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	}
	c := &configpb.Config{
		UserId: "test-user",
		Resources: []*configpb.Resource{
			newResource(t, "resource1", args),
		},
	}

	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	updatedVn, err := config.UpdateResource("resource1", &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			ResourceGroupId: "test-rg1",
			Location:        commonpb.Location_EU_WEST_1,
			CloudProvider:   commonpb.CloudProvider_AZURE,
		},
		Name:      "changed-name",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Fatalf("unable to update resource: %s", err)
	}

	assert.Equal(t, "changed-name", updatedVn.(*types.VirtualNetwork).Args.Name)
	assert.Equal(t, 1, len(config.Resources.GetAll()))
	virtualNetwork, ok := config.Resources.ResourceMap["resource1"].(*types.VirtualNetwork)
	if !ok {
		t.Fatalf("resource should of type virtual network, but is %t", config.Resources.ResourceMap["resource1"])
	}
	assert.Equal(t, "changed-name", virtualNetwork.Args.Name)
}

func TestResourceMetadataCreate_shouldNotCreateDuplicateRg(t *testing.T) {
	vnArgs := &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			ResourceGroupId: "vn-test-rg",
			Location:        commonpb.Location_EU_WEST_1,
			CloudProvider:   commonpb.CloudProvider_AZURE,
		},
		Name: "test-vn",
	}
	subnetArgs := &resourcespb.SubnetArgs{
		Name:             "vn-test-rg",
		CidrBlock:        "10.0.0.0/16",
		VirtualNetworkId: "resource1",
	}
	rgArgs := &resourcespb.ResourceGroupArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name: "vn-test-rg",
	}
	c := &configpb.Config{
		UserId: "test-user",
		Resources: []*configpb.Resource{
			newResource(t, "resource1", vnArgs),
			newResource(t, "resource2", subnetArgs),
			newResource(t, "vn-test-rg", rgArgs),
		},
	}

	config, err := resources.LoadConfig(c, metadata.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	db1, err := config.CreateResource(&resourcespb.DatabaseArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:     "db1",
		SubnetId: "resource2",
	})
	if err != nil {
		t.Fatalf("unable to create db1: %s", err)
	}

	db2, err := config.CreateResource(&resourcespb.DatabaseArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:     "db2",
		SubnetId: "resource2",
	})
	if err != nil {
		t.Fatalf("unable to create db2: %s", err)
	}

	assert.Equal(t, len(config.Resources.ResourceMap), len(config.Resources.GetAll()))
	assert.Equal(t, db1.(*types.Database).GetResourceGroupId(), db2.(*types.Database).GetResourceGroupId())
}

func newResource(t *testing.T, resourceId string, msg proto.Message) *configpb.Resource {
	a, err := anypb.New(msg)
	if err != nil {
		t.Fatalf("unable to marshal resource: %s", err)
	}

	return &configpb.Resource{ResourceId: resourceId, ResourceArgs: &configpb.ResourceArgs{ResourceArgs: a}}
}
