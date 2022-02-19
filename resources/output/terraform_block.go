package output

import (
	"fmt"
	"multy-go/validate"
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

func WrapWithBlockType(block TfBlock) any {
	if block.GetBlockType() == "resource" {
		return ResourceWrapper{R: block}
	}
	if block.GetBlockType() == "data" {
		return DataSourceWrapper{R: block}
	}
	validate.LogInternalError("unknown block type %T", block)
	return nil
}
