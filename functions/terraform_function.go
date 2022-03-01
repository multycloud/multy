package functions

import (
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"strings"
)

func getTerraformFunction(functionName string) function.Function {
	return function.New(
		&function.Spec{
			VarParam: &function.Parameter{
				Name:             "values",
				Type:             cty.DynamicPseudoType,
				AllowNull:        false,
				AllowUnknown:     false,
				AllowDynamicType: false,
				AllowMarked:      false,
			},
			Type: func(args []cty.Value) (cty.Type, error) {
				return cty.DynamicPseudoType, nil
			},
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				var argStrings []string
				for _, arg := range args {
					s, err := hclencoder.ValueToString(arg)
					if err != nil {
						return cty.NilVal, err
					}
					argStrings = append(argStrings, s)
				}
				return cty.StringVal(fmt.Sprintf("${%s(%s)}", functionName, strings.Join(argStrings, ","))), nil
			},
		},
	)
}
