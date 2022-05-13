package resources_test

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/types"
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

	config, err := resources.LoadConfig(c, types.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	assert.Equal(t, c.GetUserId(), "test-user")
	assert.Equal(t, len(config.Resources.GetAll()), 1)
	virtualNetwork, ok := config.Resources.ResourceMap["resource1"].(*types.VirtualNetwork)
	if !ok {
		t.Fatalf("resource should of type virtual network, but is %t", config.Resources.ResourceMap["resource1"])
	}
	assert.Equal(t, virtualNetwork.GetResourceId(), "resource1")
	assert.Equal(t, virtualNetwork.Args.String(), args.String())
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

	config, err := resources.LoadConfig(c, types.Metadatas)
	if err != nil {
		t.Fatalf("unable to load config: %s", err)
	}

	assert.Equal(t, c.GetUserId(), "test-user")
	assert.Contains(t, config.GetAffectedResources("resource1"), "aws_vpc.resource1")
	assert.Contains(t, config.GetAffectedResources("resource2"), "aws_vpc.resource1")
	assert.Contains(t, config.GetAffectedResources("resource1"), "aws_vpc.resource2")
	assert.Contains(t, config.GetAffectedResources("resource2"), "aws_vpc.resource2")
}

func TestResourceMetadataCreate_createsResourceGroup(t *testing.T) {
	c := &configpb.Config{
		UserId: "test-user",
	}
	config, err := resources.LoadConfig(c, types.Metadatas)
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
	assert.Equal(t, len(config.Resources.GetAll()), 2)

	// assert virtual network, except resource group id
	virtualNetwork, ok := config.Resources.ResourceMap[resourceId].(*types.VirtualNetwork)
	if !ok {
		t.Fatalf("resource should of type virtual network, but is %t", config.Resources.ResourceMap[resourceId])
	}
	assert.Equal(t, virtualNetwork.GetResourceId(), resourceId)
	actualArgsExceptRgId := proto.Clone(virtualNetwork.Args).(*resourcespb.VirtualNetworkArgs)
	actualArgsExceptRgId.CommonParameters.ResourceGroupId = ""
	assert.Equal(t, actualArgsExceptRgId.String(), args.String())

	// assert resource group
	rgId := virtualNetwork.Args.CommonParameters.ResourceGroupId
	rg, ok := config.Resources.ResourceMap[rgId].(*types.ResourceGroup)
	if !ok {
		t.Fatalf("resource should of type resource group, but is %t", config.Resources.ResourceMap[rgId])
	}
	assert.Equal(t, rg.Name, rgId)
	assert.Equal(t, rg.Location, args.CommonParameters.Location)
	assert.Equal(t, rg.Cloud, args.CommonParameters.CloudProvider)
}

func newResource(t *testing.T, resourceId string, msg proto.Message) *configpb.Resource {
	a, err := anypb.New(msg)
	if err != nil {
		t.Fatalf("unable to marshal resource: %s", err)
	}

	return &configpb.Resource{ResourceId: resourceId, ResourceArgs: &configpb.ResourceArgs{ResourceArgs: a}}
}
