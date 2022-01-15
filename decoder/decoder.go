package decoder

import (
	"fmt"
	"multy-go/hclutil"
	"multy-go/mhcl"
	"multy-go/parser"
	"multy-go/resources"
	"multy-go/resources/common"
	rg "multy-go/resources/resource_group"
	"multy-go/resources/types"
	"multy-go/validate"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

type DecodedGlobalConfig struct {
	Location      string
	Clouds        []common.CloudProvider
	DefaultRgName hcl.Expression
}

type DecodedResources struct {
	Resources    map[string]resources.CloudSpecificResource
	Providers    map[string]types.Provider
	GlobalConfig DecodedGlobalConfig
}

func Decode(config parser.ParsedConfig) *DecodedResources {
	vars := map[string]cty.Value{}
	for _, v := range config.Variables {
		vars[v.Name] = v.Value
	}
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{"var": cty.ObjectVal(vars)},
		Functions: nil,
	}

	globalConfig := decodeGlobalConfig(config, ctx)

	multyResources := topologicalSort(config.MultyResources)
	r := map[string]resources.CloudSpecificResource{}

	defaultRgId := getDefaultResourceGroupId(hclutil.GetOptionalAttributeAsExpr(config.GlobalConfig, "default_resource_group_name"))
	globalConfig.DefaultRgName = defaultRgId

	resourceDecoder := ResourceDecoder{globalConfig: globalConfig}
	cloudSpecificCtx := InitCloudSpecificContext(ctx)
	mhclProcessor := mhcl.MHCLProcessor{ResourceRefs: r}

	for _, resource := range multyResources {
		decodedResources, outCtx := resourceDecoder.Decode(*resource, cloudSpecificCtx, mhclProcessor)
		cloudSpecificCtx.AddCtx(outCtx)
		for _, decodedResource := range decodedResources {
			uniqueId := decodedResource.GetResourceId()
			if duplicate, ok := r[uniqueId]; ok {
				if !duplicate.ImplicitlyCreated || !decodedResource.ImplicitlyCreated {
					// TODO: allow specifying multiple ranges
					validate.LogFatalWithSourceRange(resource.DefinitionRange, "duplicate resources found with id %s ", uniqueId)
				}
			}
			r[uniqueId] = decodedResource
		}
	}

	return &DecodedResources{
		Resources:    r,
		GlobalConfig: globalConfig,
	}
}

func topologicalSort(allResources []*parser.MultyResource) []*parser.MultyResource {
	positions := map[*parser.MultyResource]int{}
	position := 0
	for _, resource := range allResources {
		position = topologicalSortHelper(positions, position, resource, nil)
	}

	result := make([]*parser.MultyResource, len(allResources))
	for resource, i := range positions {
		result[i] = resource
	}
	return result
}

func topologicalSortHelper(positions map[*parser.MultyResource]int, pos int, r *parser.MultyResource, currentStack []parser.MultyResourceDependency) int {
	if _, ok := positions[r]; ok {
		return pos
	}
	for i, otherResource := range currentStack {
		if otherResource.From == r {
			var errorMessages []string
			for j := i; j < len(currentStack); j++ {
				errorMessages = append(errorMessages, fmt.Sprintf("[%s]%s=>%s", currentStack[j].SourceRange, currentStack[j].From.ID, currentStack[j].To.ID))
			}
			validate.LogFatalWithSourceRange(currentStack[len(currentStack)-1].SourceRange, "found cycle while resolving dependencies for resource %s.\n%s", r.ID, strings.Join(errorMessages, "\n"))
		}
	}

	for _, dep := range r.Dependencies {
		currentStack = append(currentStack, dep)
		pos = topologicalSortHelper(positions, pos, dep.To, currentStack)
		currentStack = currentStack[:len(currentStack)-1]
	}
	positions[r] = pos
	return pos + 1

}

func decodeGlobalConfig(config parser.ParsedConfig, ctx *hcl.EvalContext) DecodedGlobalConfig {
	type globalConfigHcl struct {
		Location string   `hcl:"location"`
		Clouds   []string `hcl:"clouds"`
		HclBody  hcl.Body `hcl:",remain"`
	}

	var c globalConfigHcl

	diags := gohcl.DecodeBody(config.GlobalConfig, ctx, &c)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Unable to decode global configuration.")
	}

	globalConfig := DecodedGlobalConfig{
		Location: c.Location,
		Clouds:   castToCloudProvider(c.Clouds),
	}
	return globalConfig
}

func castToCloudProvider(c []string) []common.CloudProvider {
	var res []common.CloudProvider
	for _, cloud := range c {
		res = append(res, common.CloudProvider(cloud))
	}
	return res
}

func getDefaultResourceGroupId(specifiedDefault *hcl.Expression) hcl.Expression {
	if specifiedDefault != nil {
		return *specifiedDefault
	}
	return rg.GetDefaultResourceGroupId()
}

func resolveAttributes(attrs hcl.Attributes, ctx *hcl.EvalContext) (map[string]cty.Value, hcl.Diagnostics) {
	res := map[string]cty.Value{}
	for key, val := range attrs {
		resolved, diags := val.Expr.Value(ctx)
		if diags != nil {
			return res, diags
		}
		res[key] = resolved
	}
	return res, nil
}
