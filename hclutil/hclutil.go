package hclutil

import (
	"fmt"
	"multy-go/validate"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// TODO: figure out a way to do this programatically instead
func StringToHclExpression(s string) (hcl.Expression, error) {
	f, err := hclparse.NewParser().ParseHCL([]byte(fmt.Sprintf("e = \"%s\"", s)), "test")
	if err != nil {
		return nil, err
	}

	var expr map[string]hcl.Expression

	err = gohcl.DecodeBody(f.Body, nil, &expr)

	if err != nil {
		return nil, err
	}
	return expr["e"], nil
}

func StringToUnquotedHclExpression(s string) (hcl.Expression, error) {
	f, err := hclparse.NewParser().ParseHCL([]byte(fmt.Sprintf("e = %s", s)), "test")
	if err != nil {
		return nil, err
	}

	var expr map[string]hcl.Expression

	err = gohcl.DecodeBody(f.Body, nil, &expr)

	if err != nil {
		return nil, err
	}
	return expr["e"], nil
}

func GetOptionalAttributeAsExpr(hclBody hcl.Body, key string) *hcl.Expression {
	attrs, diags := hclBody.JustAttributes()
	if diags != nil {
		validate.LogFatalWithDiags(diags, "error while decoding hcl body")
	}
	if result, ok := attrs[key]; ok {
		return &result.Expr
	}
	return nil
}

func PartialDecode(hclBody hcl.Body, ctx *hcl.EvalContext, val any) hcl.Body {
	tValue := reflect.ValueOf(val)
	if tValue.Kind() != reflect.Ptr {
		validate.LogInternalError("target value must be a pointer, not: %s", tValue.Type().String())
	}
	t := reflect.TypeOf(val).Elem()

	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		if _, ok := t.Field(i).Tag.Lookup("hcl"); ok {
			fields = append(fields, t.Field(i))
		}
	}
	fields = append(
		fields, reflect.StructField{
			Name: "HclBody",
			Type: reflect.TypeOf((*hcl.Body)(nil)).Elem(),
			Tag:  `hcl:",remain"`,
		},
	)
	newT := reflect.StructOf(fields)
	newTInstance := reflect.New(newT)
	diags := gohcl.DecodeBody(hclBody, ctx, newTInstance.Interface())
	if diags != nil {
		validate.LogInternalError("error while trying to partial decode: %s", diags.Errs())
	}

	dest := tValue.Elem()
	for _, field := range fields {
		if field.Name != "HclBody" {
			dest.FieldByName(field.Name).Set(newTInstance.Elem().FieldByName(field.Name))
		}
	}

	return newTInstance.Elem().FieldByName("HclBody").Interface().(hcl.Body)
}

func IsNullExpr(expression hcl.Expression) bool {
	v, diags := expression.Value(nil)
	return v.IsNull() && diags == nil
}
