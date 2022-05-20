package common

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/output"
)

const (
	AWS   = commonpb.CloudProvider_AWS
	AZURE = commonpb.CloudProvider_AZURE
)

var LOCATION = map[commonpb.Location]map[commonpb.CloudProvider]string{
	commonpb.Location_EU_WEST_1: {
		AWS:   "eu-west-1",
		AZURE: "northeurope",
	},
	commonpb.Location_EU_WEST_2: {
		AWS:   "eu-west-2",
		AZURE: "uksouth",
	},
	commonpb.Location_EU_WEST_3: {
		AWS:   "eu-west-3",
		AZURE: "francecentral",
	},

	commonpb.Location_US_EAST_1: {
		AWS:   "us-east-1",
		AZURE: "eastus",
	},
	commonpb.Location_US_EAST_2: {
		AWS:   "us-east-2",
		AZURE: "eastus2",
	},
	commonpb.Location_US_WEST_1: {
		AWS:   "us-west-1",
		AZURE: "westus2",
	},
	commonpb.Location_US_WEST_2: {
		AWS:   "us-west-2",
		AZURE: "westus3",
	},
	commonpb.Location_EU_NORTH_1: {
		AWS:   "eu-north-1",
		AZURE: "swedencentral",
	},
}

var AVAILABILITY_ZONES = map[commonpb.Location]map[commonpb.CloudProvider][]string{
	commonpb.Location_EU_WEST_1: {
		AWS:   []string{"eu-west-1a", "eu-west-1b", "eu-west-1c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_EU_WEST_2: {
		AWS:   []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_EU_WEST_3: {
		AWS:   []string{"eu-west-3a", "eu-west-3b", "eu-west-3c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_US_EAST_1: {
		AWS:   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_US_EAST_2: {
		AWS:   []string{"us-east-2a", "us-east-2b", "us-east-2c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_US_WEST_1: {
		AWS:   []string{"us-west-1a", "us-west-1b"},
		AZURE: []string{"1", "2"},
	},
	commonpb.Location_US_WEST_2: {
		AWS:   []string{"us-west-1a", "us-west-1b", "us-west-1c"},
		AZURE: []string{"1", "2", "3"},
	},
	commonpb.Location_EU_NORTH_1: {
		AWS:   []string{"eu-north-1a", "eu-north-1b", "eu-north-1c"},
		AZURE: []string{"1", "2", "3"},
	},
}

// eu-west-2 "ami-0fc15d50d39e4503c"
// amzn2-ami-hvm-2.0.20211103.0-x86_64-gp2
//https://cloud-images.ubuntu.com/locator/ec2/
var AMIMAP = map[string]string{
	"eu-west-1":  "ami-09d4a659cdd8677be",
	"eu-west-2":  "ami-0fc15d50d39e4503c",
	"eu-west-3":  "ami-0fc15d50d39e4503c",
	"us-east-1":  "ami-04ad2567c9e3d7893",
	"us-east-2":  "ami-04ad2567c9e3d7893",
	"us-west-1":  "ami-04ad2567c9e3d7893",
	"us-west-2":  "ami-04ad2567c9e3d7893",
	"eu-north-1": "ami-04ad2567c9e3d7893",
}

var AwsAmiOwners = map[resourcespb.ImageReference_OperatingSystemDistribution]string{
	resourcespb.ImageReference_UBUNTU:  "099720109477",
	resourcespb.ImageReference_CENT_OS: "125523088429",
	resourcespb.ImageReference_DEBIAN:  "136693071363",
}

func GetAvailabilityZone(location commonpb.Location, az int, cloud commonpb.CloudProvider) (string, error) {
	if AVAILABILITY_ZONES[location] == nil {
		return "", fmt.Errorf("invalid location: %s", location)
	}
	azArray := AVAILABILITY_ZONES[location][cloud]
	if az == 0 {
		return "", nil
	}
	if az <= len(azArray) {
		return azArray[az-1], nil
	}
	return "", fmt.Errorf("invalid az value: %d", az)
}

func GetCloudLocation(location commonpb.Location, provider commonpb.CloudProvider) (string, error) {
	if _, ok := LOCATION[location]; !ok {
		return "", fmt.Errorf("location %s is not defined", location)
	}
	if _, ok := LOCATION[location][provider]; !ok {
		return "", fmt.Errorf("location %s is not defined for cloud %s", location, provider)
	}
	return LOCATION[location][provider], nil
}

type AzResource struct {
	output.TerraformResource `hcl:",squash"`
	ResourceGroupName        string `hcl:"resource_group_name,expr" hcle:"omitempty"`
	Name                     string `hcl:"name" hcle:"omitempty"`
	Location                 string `hcl:"location" hcle:"omitempty"`
}

func NewAwsResource(resourceId string, name string) *AwsResource {
	return &AwsResource{
		TerraformResource: output.TerraformResource{ResourceId: resourceId},
		Tags:              map[string]string{"Name": name}}
}

func NewAwsResourceWithDeps(resourceId string, name string, deps []string) *AwsResource {
	return &AwsResource{
		TerraformResource: output.TerraformResource{ResourceId: resourceId, DependsOn: deps},
		Tags:              map[string]string{"Name": name}}
}

func NewAwsResourceWithIdOnly(resourceId string) *AwsResource {
	return &AwsResource{
		TerraformResource: output.TerraformResource{ResourceId: resourceId}}
}

func NewAzResource(resourceId string, name string, rgName string, location string) *AzResource {
	return &AzResource{
		TerraformResource: output.TerraformResource{ResourceId: resourceId},
		Name:              name,
		ResourceGroupName: rgName,
		Location:          location,
	}
}

func (r *AwsResource) SetName(name string) {
	//r.ResourceName = name
	r.TerraformResource.ResourceName = name
}

func (r *AzResource) SetName(name string) {
	r.TerraformResource.ResourceName = name
}

type AwsResource struct {
	output.TerraformResource `hcl:",squash"`
	Tags                     map[string]string `hcl:"tags" hcle:"omitempty"`
}
