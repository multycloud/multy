package output

import (
	"fmt"
	"github.com/multycloud/multy/validate"
	"reflect"
	"strings"
)

// ResourceWrapper just to add a resource {} around when encoding into hcl
type ResourceWrapper struct {
	R any `hcl:"resource"`
}

type DataSourceWrapper struct {
	R any `hcl:"data"`
}

func (r DataSourceWrapper) GetR() any {
	return r.R
}

func (r ResourceWrapper) GetR() any {
	return r.R
}

type TfBlock interface {
	GetFullResourceRef() string
	GetBlockType() string
	AddDependency(string)
	GetResourceId() string
}

type TerraformResource struct {
	ResourceName string   `hcl:",key"`
	ResourceId   string   `hcl:",key"`
	DependsOn    []string `hcl:"depends_on,expr" hcle:"omitempty"`
}

func (t TerraformResource) GetFullResourceRef() string {
	return fmt.Sprintf("%s.%s", t.ResourceName, t.ResourceId)
}

func (t TerraformResource) GetBlockType() string {
	return "resource"
}

func (t *TerraformResource) AddDependency(dep string) {
	t.DependsOn = append(t.DependsOn, dep)
}

func (t *TerraformResource) SetName(name string) {
	t.ResourceName = name
}

func (t *TerraformResource) GetResourceId() string {
	return t.ResourceId
}

type TerraformDataSource struct {
	ResourceName string   `hcl:",key"`
	ResourceId   string   `hcl:",key"`
	DependsOn    []string `hcl:"depends_on,expr"  hcle:"omitempty"`
}

func (t TerraformDataSource) GetFullResourceRef() string {
	return fmt.Sprintf("data.%s.%s", t.ResourceName, t.ResourceId)
}

func (t TerraformDataSource) GetBlockType() string {
	return "data"
}

func (t *TerraformDataSource) AddDependency(dep string) {
	t.DependsOn = append(t.DependsOn, dep)
}

func (t *TerraformDataSource) SetName(name string) {
	t.ResourceName = name
}

func (t *TerraformDataSource) GetResourceId() string {
	return t.ResourceId
}

func WrapWithBlockType(block TfBlock) (any, error) {
	if block.GetBlockType() == "resource" {
		return ResourceWrapper{R: block}, nil
	}
	if block.GetBlockType() == "data" {
		return DataSourceWrapper{R: block}, nil
	}
	return nil, fmt.Errorf("unknown block type %T", block)
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
