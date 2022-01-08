package decoder

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/validate"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type CloudSpecificContext struct {
	cloudDependentCtx *hcl.EvalContext
}

var cloudFunctions = map[common.CloudProvider]function.Function{
	common.AWS:   getCloudSpecificValFunction(common.AWS),
	common.AZURE: getCloudSpecificValFunction(common.AZURE),
}

func (c CloudSpecificContext) GetContext(cloud common.CloudProvider) *hcl.EvalContext {
	vars := map[string]cty.Value{}
	if v, ok := c.cloudDependentCtx.Variables[string(cloud)]; ok {
		vars = v.AsValueMap()
	}
	result := c.cloudDependentCtx.NewChild()
	result.Variables = vars
	result.Functions = map[string]function.Function{
		"cloud_specific_value": cloudFunctions[cloud],
	}
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

func InitCloudSpecificContext(ctx *hcl.EvalContext) CloudSpecificContext {
	if ctx == nil {
		ctx = &hcl.EvalContext{
			Variables: map[string]cty.Value{},
			Functions: map[string]function.Function{},
		}
	}
	return CloudSpecificContext{cloudDependentCtx: ctx}
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

func getCloudSpecificValFunction(cloud common.CloudProvider) function.Function {
	return function.New(&function.Spec{
		Params: []function.Parameter{{
			Name:             "values",
			Type:             getCloudSpecificValArgType(),
			AllowNull:        false,
			AllowUnknown:     false,
			AllowDynamicType: false,
			AllowMarked:      false,
		}},
		Type: func(args []cty.Value) (cty.Type, error) {
			valueMap := args[0].AsValueMap()
			returnType := cty.NilType
			for key, val := range valueMap {
				if val.IsNull() {
					continue
				}

				if returnType == cty.NilType {
					returnType = val.Type()
				}

				if errs := val.Type().TestConformance(returnType); errs != nil {
					return returnType, fmt.Errorf("value type for %s does not comform with other values: %s", key, errs[0].Error())
				}
			}
			return returnType, nil
		},
		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
			allValues := args[0].AsValueMap()
			if val, ok := allValues[string(cloud)]; ok && !val.IsNull() {
				return val, nil
			} else if val, ok := allValues["default"]; ok && !val.IsNull() {
				return val, nil
			} else {
				return cty.NilVal, fmt.Errorf("no value for %s", cloud)
			}
		},
	})
}

func getCloudSpecificValArgType() cty.Type {
	allTypes := map[string]cty.Type{"default": cty.DynamicPseudoType}
	optionalAttrs := []string{"default"}
	for _, cloud := range common.GetAllCloudProviders() {
		allTypes[string(cloud)] = cty.DynamicPseudoType
		optionalAttrs = append(optionalAttrs, string(cloud))
	}
	return cty.ObjectWithOptionalAttrs(allTypes, optionalAttrs)
}
