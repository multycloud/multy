package common

import (
	"fmt"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/validate"
	"reflect"
	"strings"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	AZURE               = "azure"
)

func GetAllCloudProviders() []CloudProvider {
	return []CloudProvider{AWS, AZURE}
}

func AsCloudProvider(s string) (CloudProvider, bool) {
	for _, cloud := range GetAllCloudProviders() {
		if string(cloud) == s {
			return cloud, true
		}
	}
	return "", false
}

const (
	IRELAND = "ireland"
	UK      = "uk"
	USEAST  = "us-east"
)

var LOCATION = map[string]map[CloudProvider]string{
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

var AVAILABILITY_ZONES = map[string]map[CloudProvider][]string{
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

var RESOURCETYPES = map[string]string{
	"virtual_network":        "vn",
	"subnet":                 "vn",
	"route_table":            "vn",
	"database":               "db",
	"virtual_machine":        "vm",
	"network_security_group": "nsg",
	"object_storage":         "st",
	"object_storage_object":  "st",
	"public_ip":              "pip",
	"network_interface":      "nic",
	"vault":                  "kv",
	"vault_secret":           "kv",
	"vault_access_policy":    "kv",
	"lambda":                 "fun",
	"kubernetes_service":     "ks",
	"kubernetes_node_pool":   "ks",
}

// eu-west-2 "ami-0fc15d50d39e4503c"
// amzn2-ami-hvm-2.0.20211103.0-x86_64-gp2
//https://cloud-images.ubuntu.com/locator/ec2/
var AMIMAP = map[string]string{
	"eu-west-1": "ami-09d4a659cdd8677be",
	"us-east-1": "ami-04ad2567c9e3d7893",
}

func GetResourceTypeAbbreviation(t string) string {
	if val, ok := RESOURCETYPES[t]; ok {
		return val
	}
	validate.LogInternalError("no resource abbreviation for: %s", t)
	return ""
}

func GetAvailabilityZone(location string, az int, cloud CloudProvider) string {
	azArray := AVAILABILITY_ZONES[location][cloud]
	if az == 0 {
		return ""
	}
	if az <= len(azArray) {
		return azArray[az-1]
	}
	validate.LogInternalError("invalid az value: %d", az)
	return ""

}

func GetCloudLocation(location string, provider CloudProvider) (string, error) {
	if _, ok := LOCATION[location]; !ok {
		return "", fmt.Errorf("location %s is not defined", location)
	}
	if _, ok := LOCATION[location][provider]; !ok {
		return "", fmt.Errorf("location %s is not defined for cloud %s", location, provider)
	}
	return LOCATION[location][provider], nil
}

func ValidateVmSize(s string) bool {
	_, ok := VMSIZE[s]
	return ok
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

func GetResourceName(r any) string {
	t := reflect.TypeOf(r)
	tagValue, ok := t.Field(0).Tag.Lookup("default")
	if !ok {
		validate.LogInternalError("no default resource name found")
	}
	tagValues := strings.Split(tagValue, ",")
	for _, v := range tagValues {
		keyVal := strings.SplitN(v, "=", 2)
		if keyVal[0] == "name" {
			return keyVal[1]
		}
	}
	validate.LogInternalError("no default resource name found")
	return ""
}

type AwsResource struct {
	output.TerraformResource `hcl:",squash"`
	Tags                     map[string]string `hcl:"tags" hcle:"omitempty"`
}
