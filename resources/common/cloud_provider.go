package common

import (
	"fmt"
	"multy-go/validate"
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
}

var AVAILABILITY_ZONES = map[string]map[CloudProvider][]string{
	UK: {
		AWS:   []string{"eu-west-2a", "eu-west-2b", "eu-west2c"},
		AZURE: []string{"1", "2", "3"},
	},
	IRELAND: {
		AWS:   []string{"eu-west-1a", "eu-west-1b", "eu-west1c"},
		AZURE: []string{"1", "2", "3"},
	},
}

type Size string

const (
	MICRO = "nano"
)

var DBSIZE = map[string]map[CloudProvider]string{
	MICRO: {
		AWS:   "db.t2.micro",
		AZURE: "GP_Gen5_2",
	},
}

var VMSIZE = map[string]map[CloudProvider]string{
	MICRO: {
		AWS:   "t2.nano",
		AZURE: "Standard_B1ls",
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
	"lambda":                 "fun",
}

func GetResourceTypeAbbreviation(t string) string {
	if val, ok := RESOURCETYPES[t]; ok {
		return val
	}
	validate.LogInternalError("unexpected resource type name: %s", t)
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

type AwsResource struct {
	ResourceName string            `hcl:",key"`
	ResourceId   string            `hcl:",key"`
	Tags         map[string]string `hcl:"tags" hcle:"omitempty"`
}

type AzResource struct {
	ResourceName      string `hcl:",key"`
	ResourceId        string `hcl:",key"`
	ResourceGroupName string `hcl:"resource_group_name,expr" hcle:"omitempty"`
	Name              string `hcl:"name" hcle:"omitempty"`
	Location          string `hcl:"location" hcle:"omitempty"`
}

func NewAwsResource(resourceName string, resourceId string, name string) AwsResource {
	return AwsResource{ResourceName: resourceName, ResourceId: resourceId, Tags: map[string]string{"Name": name}}
}

func NewAzResource(resourceName string, resourceId string, name string, rgName string, location string) AzResource {
	return AzResource{
		ResourceName:      resourceName,
		ResourceId:        resourceId,
		Name:              name,
		ResourceGroupName: rgName,
		Location:          location,
	}
}
