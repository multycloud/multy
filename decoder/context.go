package decoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/functions"
	"github.com/multycloud/multy/parser"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type CloudSpecificContext struct {
	cloudDependentCtx *hcl.EvalContext
	variables         []parser.ParsedVariable
}

func (c CloudSpecificContext) GetContext(cloud common.CloudProvider) *hcl.EvalContext {
	vars := map[string]cty.Value{}
	if v, ok := c.cloudDependentCtx.Variables[string(cloud)]; ok {
		vars = v.AsValueMap()
	}
	result := c.cloudDependentCtx.NewChild()
	result.Variables = vars
	result.Functions = functions.GetAllFunctions(cloud)
	return result
}

func (c CloudSpecificContext) GetCloudAgnosticContext() *hcl.EvalContext {
	result := c.cloudDependentCtx.NewChild()
	result.Variables = map[string]cty.Value{}
	result.Functions = functions.GetAllCloudAgnosticFunctions()
	return result
}

func (c CloudSpecificContext) AddCloudDependentVar(key string, v cty.Value, cloud common.CloudProvider) {
	c.AddCloudDependentVars(map[string]cty.Value{key: v}, cloud)
}

func (c CloudSpecificContext) AddCloudDependentVars(vals map[string]cty.Value, cloud common.CloudProvider) {
	vars := map[string]cty.Value{}
	if _, ok := c.cloudDependentCtx.Variables[string(cloud)]; ok {
		vars = c.cloudDependentCtx.Variables[string(cloud)].AsValueMap()
	}
	for key, val := range vals {
		vars[key] = val
	}
	c.cloudDependentCtx.Variables[string(cloud)] = cty.ObjectVal(vars)
}

func (c CloudSpecificContext) AddVar(key string, v cty.Value) {
	c.cloudDependentCtx.Variables[key] = v
}

func InitCloudSpecificContext(ctx *hcl.EvalContext, variables []parser.ParsedVariable) CloudSpecificContext {
	if ctx == nil {
		ctx = &hcl.EvalContext{
			Variables: map[string]cty.Value{},
			Functions: map[string]function.Function{},
		}
	}
	result := CloudSpecificContext{cloudDependentCtx: ctx, variables: variables}
	for _, c := range common.GetAllCloudProviders() {
		vars := map[string]cty.Value{}
		for _, v := range variables {
			vars[v.Name] = v.Value(c)
		}
		result.AddCloudDependentVar("var", cty.ObjectVal(vars), c)
	}
	return result
}

func (c CloudSpecificContext) AddCtx(otherC CloudSpecificContext) {
	for key, value := range otherC.cloudDependentCtx.Variables {
		if _, ok := c.cloudDependentCtx.Variables[key]; !ok {
			c.cloudDependentCtx.Variables[key] = value
			continue
		}
		if cloud, ok := common.AsCloudProvider(key); ok {
			c.AddCloudDependentVars(value.AsValueMap(), cloud)
		} else {
			validate.LogInternalError("attempted to override context var for key %s", key)
		}
	}
	for key, value := range otherC.cloudDependentCtx.Functions {
		if _, ok := c.cloudDependentCtx.Functions[key]; ok {
			validate.LogInternalError("attempted to override context function for key %s", key)
		}
		c.cloudDependentCtx.Functions[key] = value
	}
}
