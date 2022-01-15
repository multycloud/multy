package decoder

import (
	"fmt"
	"multy-go/hclutil"
	"multy-go/mhcl"
	"multy-go/parser"
	"multy-go/resources"
	"multy-go/resources/common"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

type ResourceDecoder struct {
	globalConfig DecodedGlobalConfig
}

func (d ResourceDecoder) Decode(resource parser.MultyResource, ctx CloudSpecificContext, mhclProcessor mhcl.MHCLProcessor) ([]resources.CloudSpecificResource, CloudSpecificContext) {
	attrs, diags := resource.HCLBody.JustAttributes()
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Unable to parse attributes of resource %s.", resource.ID)
	}
	clouds := d.globalConfig.Clouds
	if cloudsAttr, ok := attrs["clouds"]; ok {
		clouds = []common.CloudProvider{}
		value, diags := cloudsAttr.Expr.Value(ctx.cloudDependentCtx)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "Unable to parse cloud list.")
		}
		for _, cloudValue := range value.AsValueSlice() {
			if provider, ok := common.AsCloudProvider(cloudValue.AsString()); ok {
				clouds = append(clouds, provider)
			} else {
				validate.LogFatalWithSourceRange(cloudsAttr.Range, fmt.Sprintf("Unknown cloud provider %s.", cloudValue.AsString()))
			}
		}
	}
	resultCtx := InitCloudSpecificContext(nil)
	var result []resources.CloudSpecificResource
	for _, cloud := range clouds {
		cloudCtx := ctx.GetContext(cloud)
		rgId, shouldCreateRg := getRgId(d.globalConfig.DefaultRgName, attrs, cloudCtx, resource)
		r, diags := decode(resource, cloudCtx, rgId, mhclProcessor)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "Unable to decode resource %s.", resource.ID)
		}
		cloudSpecificR := resources.CloudSpecificResource{
			Cloud:             cloud,
			Resource:          r,
			ImplicitlyCreated: false,
		}
		result = append(result, cloudSpecificR)
		allVars, diags := resolveAttributes(attrs, cloudCtx)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "Unable to resolve input variables.")
		}
		allVars["id"] = cty.StringVal(cloudSpecificR.GetResourceId())
		outputVars := r.GetOutputValues(cloud)
		for key, val := range outputVars {
			allVars[key] = val
		}
		resultCtx.AddCloudDependentVars(map[string]cty.Value{
			resource.ID: cty.ObjectVal(allVars),
		}, cloud)

		if shouldCreateRg {
			result = append(result, resources.CloudSpecificResource{ImplicitlyCreated: true, Cloud: cloud, Resource: rg.DefaultResourceGroup(rgId)})
		}
	}
	return result, resultCtx
}

func decode(resource parser.MultyResource, ctx *hcl.EvalContext, rgId string, mhclProcessor mhcl.MHCLProcessor) (resources.Resource, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	if resource.ID == "" {
		validate.LogFatalWithSourceRange(resource.DefinitionRange, "found resource of type '%s' with no id", resource.Type)
	}

	commonParams := &resources.CommonResourceParams{
		ResourceId:      resource.ID,
		ResourceGroupId: rgId,
	}
	body := hclutil.PartialDecode(resource.HCLBody, ctx, commonParams)

	r, err := InitResource(resource.Type, commonParams)
	if err != nil {
		validate.LogFatalWithSourceRange(resource.DefinitionRange, err.Error())
	}

	body = mhclProcessor.Process(body, r, ctx)

	validationInfo, diags := decodeBody(r, body, ctx)
	commonParams.ResourceValidationInfo = validationInfo

	return r, diags
}

func decodeBody(t interface{}, body hcl.Body, ctx *hcl.EvalContext) (*validate.ResourceValidationInfo, hcl.Diagnostics) {
	schema, _ := gohcl.ImpliedBodySchema(t)
	content, diags := body.Content(schema)
	if diags != nil {
		return nil, diags
	}
	diags = gohcl.DecodeBody(body, ctx, t)
	return validate.NewResourceValidationInfoFromContent(content), diags
}

func getRgId(defaultRgId hcl.Expression, attrs hcl.Attributes, ctx *hcl.EvalContext, resource parser.MultyResource) (string, bool) {
	// TODO: refactor so that we don't do explicit checks here
	if resource.Type == "resource_group" || resource.Type == "route_table_association" {
		return "", false
	}
	shouldCreateRg := true
	// populate resource group name if empty
	var rgIdExpr = defaultRgId
	if rgIdAttr, ok := attrs["resource_group_id"]; ok {
		rgIdExpr = rgIdAttr.Expr
		shouldCreateRg = false
	}

	var rgVars = cty.ObjectVal(map[string]cty.Value{})
	var diags hcl.Diagnostics

	if rgVarsAttrs, ok := attrs["rg_vars"]; ok {
		rgVars, diags = rgVarsAttrs.Expr.Value(ctx)
		if diags != nil {
			validate.LogFatalWithDiags(diags, "Unable to resolve rg vars.")
		}
	}

	// create resource group if needed
	rgId := ""
	childCtx := hcl.EvalContext{
		Variables: ctx.Variables,
		Functions: ctx.Functions,
	}
	childCtx.Variables["resource_type"] = cty.StringVal(common.GetResourceTypeAbbreviation(resource.Type))
	childCtx.Variables["rg_vars"] = rgVars
	diags = gohcl.DecodeExpression(rgIdExpr, ctx, &rgId)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Unable to resolve resource group name for resource %s. "+
			"Are you missing a rg_var?", resource.ID)
	}
	return rgId, shouldCreateRg
}
