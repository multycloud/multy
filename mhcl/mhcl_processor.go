package mhcl

import (
	"github.com/hashicorp/hcl/v2"
	"multy-go/resources"
	"multy-go/validate"
	"reflect"
	"strings"
)

// format: mhcl:"{commnad},hcl_tags"
// example: mhcl:"ref=subnet,optional"
const tagName = "mhcl"

type MHCLProcessor struct {
	// These will be used to resolve mhcl refs
	ResourceRefs map[string]resources.CloudSpecificResource
}

type MHclTag struct {
	// Reference to a multy resource
	Ref      string
	Optional bool
}

func (p *MHCLProcessor) Process(body hcl.Body, r interface{}, ctx *hcl.EvalContext) hcl.Body {
	// Find tag
	tValue := reflect.ValueOf(r)
	if tValue.Kind() != reflect.Ptr {
		validate.LogInternalError("target value must be a pointer, not: %s", tValue.Type().String())
	}
	t := reflect.TypeOf(r).Elem()

	// Create Hcl schema
	schema := hcl.BodySchema{}
	refFields := map[string]reflect.StructField{}

	for i := 0; i < t.NumField(); i++ {
		if tagValue, ok := t.Field(i).Tag.Lookup("mhcl"); ok {

			mhcltag := parseMHclTag(tagValue)
			if _, ok := refFields[mhcltag.Ref]; ok {
				validate.LogInternalError("duplicate mhcl tag for ref %s", mhcltag.Ref)
			}
			refFields[mhcltag.Ref] = t.Field(i)
			schema.Attributes = append(schema.Attributes, hcl.AttributeSchema{
				Name:     mhcltag.Ref,
				Required: !mhcltag.Optional,
			})
		}
	}

	// Partial decode using gohcl
	content, remaining, diags := body.PartialContent(&schema)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "error while processing multy reference")
	}

	// Set field of r by getting the reference from ResourceRefs
	for _, attr := range content.Attributes {
		val, diags := attr.Expr.Value(ctx)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "error while evaluating attribute")
		}
		if !val.Type().HasAttribute("id") {
			validate.LogFatalWithSourceRange(attr.Range, "expected a multy resource, but got %s", val.Type().GoString())
		}
		refId := val.AsValueMap()["id"].AsString()
		if resource, ok := p.ResourceRefs[refId]; ok {
			actualType := reflect.TypeOf(resource.Resource).Elem()
			expectedType := refFields[attr.Name].Type.Elem()
			if !actualType.AssignableTo(expectedType) {
				validate.LogFatalWithSourceRange(attr.Range, "expected resource of type %s but got %s", expectedType.Name(), actualType.Name())
			}
			tValue.Elem().FieldByName(refFields[attr.Name].Name).Set(reflect.ValueOf(resource.Resource))
		} else {
			validate.LogInternalError("unknown resource %s", refId)
		}
	}

	return remaining

}

func parseMHclTag(value string) MHclTag {
	result := MHclTag{Optional: false}
	values := strings.SplitN(value, ",", 2)
	mainValue := values[0]
	if strings.HasPrefix(mainValue, "ref=") {
		result.Ref = strings.TrimPrefix(mainValue, "ref=")
	} else {
		validate.LogInternalError("unknown tag value for mhcl: %s", mainValue)
	}

	if len(values) >= 2 {
		switch values[1] {
		case "optional":
			result.Optional = true
		default:
			validate.LogInternalError("unsupported option for mhcl: %s", values[1])
		}
	}

	return result
}
