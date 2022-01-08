package parser

import (
	"multy-go/resources/common"
	"multy-go/validate"
	"multy-go/variables"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Parser struct {
	CliVars variables.CommandLineVariables
}

type MultyResource struct {
	Type            string   `hcl:"type,label"`
	ID              string   `hcl:"id,label"`
	HCLBody         hcl.Body `hcl:",remain"`
	Dependencies    []MultyResourceDependency
	DefinitionRange hcl.Range
}

type MultyResourceDependency struct {
	From        *MultyResource
	To          *MultyResource
	SourceRange hcl.Range
}

type MultyConfig struct {
	Location                 string
	DefaultResourceGroupName *hcl.Expression
	Clouds                   []string
}

type ParsedConfig struct {
	Variables      []ParsedVariable
	MultyResources []*MultyResource
	GlobalConfig   hcl.Body
}

type ParsedVariable struct {
	Name  string
	Value cty.Value
}

func (p *Parser) Parse(filepath string) ParsedConfig {
	type multyConfig struct {
		HCLBody hcl.Body `hcl:",remain"`
	}
	type config struct {
		Variables      []variables.Variable `hcl:"variable,block"`
		MultyResources []*MultyResource     `hcl:"multy,block"`
		MultyConfig    multyConfig          `hcl:"config,block"`
	}
	var c config
	parser := hclparse.NewParser()
	f, _ := parser.ParseHCLFile(filepath)
	diags := gohcl.DecodeBody(f.Body, nil, &c)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Failed to load configuration.")
	}

	resourcesById := map[string]*MultyResource{}
	for _, r := range c.MultyResources {
		resourcesById[r.ID] = r
	}

	rootSchema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "multy", LabelNames: []string{"type", "id"}},
		},
	}

	content, _, diags := f.Body.PartialContent(rootSchema)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Failed to load configuration.")
	}

	for _, block := range content.Blocks {
		resourcesById[block.Labels[1]].DefinitionRange = block.DefRange
	}

	for _, r := range c.MultyResources {
		r.Dependencies = findDependencies(r, resourcesById)
	}

	return ParsedConfig{
		Variables:      convertVars(c.Variables, p.CliVars),
		MultyResources: c.MultyResources,
		GlobalConfig:   c.MultyConfig.HCLBody,
	}

}

func findDependencies(resource *MultyResource, resourcesById map[string]*MultyResource) []MultyResourceDependency {
	var result []MultyResourceDependency
	attrs, diags := resource.HCLBody.JustAttributes()
	if diags != nil {
		validate.LogFatalWithDiags(diags, "unable to get attributes %s")
	}
	for _, attr := range attrs {
		for _, v := range attr.Expr.Variables() {
			if _, ok := common.AsCloudProvider(v.RootName()); ok {
				if traverseAttr, ok := v[1].(hcl.TraverseAttr); ok {
					if _, ok := resourcesById[traverseAttr.Name]; ok {
						result = append(result, MultyResourceDependency{
							From:        resource,
							To:          resourcesById[traverseAttr.Name],
							SourceRange: v.SourceRange(),
						})
					}
				} else {
					validate.LogFatalWithSourceRange(v.SourceRange(), "expected attr lookup when referencing a cloud name (aws, az, ...)")
				}
			} else {
				if _, ok := resourcesById[v.RootName()]; ok {
					result = append(result, MultyResourceDependency{
						From:        resource,
						To:          resourcesById[v.RootName()],
						SourceRange: v.SourceRange(),
					})
				}
			}
		}
	}
	return result
}

func convertVars(allVars []variables.Variable, cliVars variables.CommandLineVariables) []ParsedVariable {
	cliVarsByName := map[string]string{}
	for _, cliVar := range cliVars {
		cliVarsByName[cliVar.Name] = cliVar.Value
	}
	var result []ParsedVariable
	for _, v := range allVars {
		var val cty.Value
		if cliVar, ok := cliVarsByName[v.Name]; ok {
			val = v.GetValue(&cliVar)
		}
		val = v.GetValue(nil)
		result = append(result, ParsedVariable{
			Name:  v.Name,
			Value: val,
		})
	}
	return result
}
