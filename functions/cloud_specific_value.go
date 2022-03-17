package functions

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"multy/resources/common"
)

func getCloudSpecificValFunction(cloud common.CloudProvider, throwOnMissingVal bool) function.Function {
	return function.New(
		&function.Spec{
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
						return returnType, fmt.Errorf(
							"value type for %s does not comform with other values: %s", key, errs[0].Error(),
						)
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
					if throwOnMissingVal {
						return cty.NilVal, fmt.Errorf("no value for %s", cloud)
					} else {
						return cty.NullVal(retType), nil
					}
				}
			},
		},
	)
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

func GetCloudSpecificValueFunction(cloud common.CloudProvider, throwOnError bool) function.Function {
	return getCloudSpecificValFunction(cloud, throwOnError)
}
