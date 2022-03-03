package variables

import (
	"bytes"
	"fmt"
	"multy-go/functions"
	"multy-go/resources/common"
	"multy-go/validate"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

type Variable struct {
	Name    string         `hcl:"name,label"`
	Default hcl.Expression `hcl:"default"`
	Type    hcl.Expression `hcl:"type"`
}

func (v Variable) GetValueFromCli(cliVar *string) cty.Value {
	if cliVar == nil {
		validate.LogInternalError("nil cli var for var %s", v.Name)
	}

	varType, diags := decodeVariableType(v.Type)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "unable to decode variable type of variable '%s'", v.Name)
	}

	val, err := json.Unmarshal([]byte(*cliVar), varType)
	if err != nil {
		validate.LogInternalError("unable to parse %+v as type %s. err: %s", v, varType, err.Error())
	}

	return val
}

func (v Variable) GetDefaultValueFunction() func(common.CloudProvider) cty.Value {
	return func(cloud common.CloudProvider) cty.Value {
		val, diags := v.Default.Value(
			&hcl.EvalContext{
				Variables: map[string]cty.Value{},
				Functions: functions.GetAllFunctionsForVarEvaluation(cloud)},
		)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "unable to decode variable default value of variable '%s'", v.Name)
		}
		return val
	}
}

func decodeVariableType(expr hcl.Expression) (cty.Type, hcl.Diagnostics) {

	// First we'll deal with some shorthand forms that the HCL-level type
	// expression parser doesn't include. These both emulate pre-0.12 behavior
	// of allowing a list or map of any element type as long as all of the
	// elements are consistent. This is the same as list(any) or map(any).
	switch hcl.ExprAsKeyword(expr) {
	case "list":
		return cty.List(cty.DynamicPseudoType), nil
	case "map":
		return cty.Map(cty.DynamicPseudoType), nil
	}

	ty, diags := typeexpr.TypeConstraint(expr)
	if diags.HasErrors() {
		return cty.DynamicPseudoType, diags
	}

	switch {
	case ty.IsPrimitiveType():
		// Primitive types use literal parsing.
		return ty, diags
	default:
		// Everything else uses HCL parsing
		return ty, diags
	}
}

type CommandLineVariable struct {
	Name  string
	Value string
}

type CommandLineVariables []CommandLineVariable

func (i *CommandLineVariables) String() string {
	var b bytes.Buffer
	for _, v := range *i {
		b.WriteString(fmt.Sprintf("%s=%+v", v.Name, v.Value))
	}
	return b.String()
}

func (i *CommandLineVariables) Set(value string) error {
	split := strings.SplitAfterN(value, "=", 2)
	varName := split[0][:len(split[0])-1]
	varValue := split[1]
	*i = append(
		*i, CommandLineVariable{
			Name:  varName,
			Value: varValue,
		},
	)
	return nil
}

func (i *CommandLineVariables) Type() string {
	return "string"
}
