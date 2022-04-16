package common

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/resources/output"
)

const (
	AWS   = commonpb.CloudProvider_AWS
	AZURE = commonpb.CloudProvider_AZURE
)
const (
	IRELAND = "ireland"
	UK      = "uk"
	USEAST  = "us_east"
)

var LOCATION = map[string]map[commonpb.CloudProvider]string{
	UK: {
		AWS:   "eu-west-2",
		AZURE: "ukwest",
	},
	IRELAND: {
		AWS:   "eu-west-1",
		AZURE: "northeurope",
	},
	USEAST: {
		AWS:   "us-east-1",
		AZURE: "eastus",
	},
}

var AVAILABILITY_ZONES = map[string]map[commonpb.CloudProvider][]string{
	UK: {
		AWS:   []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"},
		AZURE: []string{"1", "2", "3"},
	},
	IRELAND: {
		AWS:   []string{"eu-west-1a", "eu-west-1b", "eu-west-1c"},
		AZURE: []string{"1", "2", "3"},
	},
	USEAST: {
		AWS:   []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		AZURE: []string{"1", "2", "3"},
	},
}

// eu-west-2 "ami-0fc15d50d39e4503c"
// amzn2-ami-hvm-2.0.20211103.0-x86_64-gp2
//https://cloud-images.ubuntu.com/locator/ec2/
var AMIMAP = map[string]string{
	"eu-west-1": "ami-09d4a659cdd8677be",
	"eu-west-2": "ami-0fc15d50d39e4503c",
	"us-east-1": "ami-04ad2567c9e3d7893",
}

func GetAvailabilityZone(location string, az int, cloud commonpb.CloudProvider) (string, error) {
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

func GetCloudLocation(location string, provider commonpb.CloudProvider) (string, error) {
	if _, ok := LOCATION[location]; !ok {
		return "", fmt.Errorf("location %s is not defined", location)
	}
	if _, ok := LOCATION[location][provider]; !ok {
		return "", fmt.Errorf("location %s is not defined for cloud %s", location, provider)
	}
	return LOCATION[location][provider], nil
}

func GetCloudLocationPb(location string, provider commonpb.CloudProvider) (string, error) {
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
