package decoder

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/hclutil"
	"github.com/multycloud/multy/mhcl"
	"github.com/multycloud/multy/parser"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"

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
				validate.LogFatalWithSourceRange(
					cloudsAttr.Range, fmt.Sprintf("Unknown cloud provider %s.", cloudValue.AsString()),
				)
			}
		}
	}
	resultCtx := InitCloudSpecificContext(nil, nil)
	var result []resources.CloudSpecificResource
	for _, cloud := range clouds {
		cloudCtx := ctx.GetContext(cloud)
		rgId, shouldCreateRg := GetRgId(d.globalConfig.DefaultRgName, attrs, cloudCtx, resource.Type)
		r, diags := decode(resource, cloudCtx, rgId, mhclProcessor, cloud)
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
		allVars["id"] = cty.StringVal(cloudSpecificR.GetMainOutputId())
		allVars["multy_id"] = cty.StringVal(cloudSpecificR.GetResourceId())
		outputVars := r.GetOutputValues(cloud)
		for key, val := range outputVars {
			allVars[key] = val
		}
		resultCtx.AddCloudDependentVars(
			map[string]cty.Value{
				resource.ID: cty.ObjectVal(allVars),
			}, cloud,
		)

		if shouldCreateRg {
			result = append(
				result,
				resources.CloudSpecificResource{ImplicitlyCreated: true, Cloud: cloud, Resource: rg.DefaultResourceGroup(rgId)},
			)
		}
	}
	return result, resultCtx
}

func decode(resource parser.MultyResource, ctx *hcl.EvalContext, rgId string, mhclProcessor mhcl.MHCLProcessor, cloud common.CloudProvider) (resources.Resource, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	if resource.ID == "" {
		validate.LogFatalWithSourceRange(
			resource.DefinitionRange, "found resource of type '%s' with no id", resource.Type,
		)
	}

	commonParams := &resources.CommonResourceParams{
		ResourceId:      resource.ID,
		ResourceGroupId: rgId,
		DependsOn: util.FlatMapSliceValues(
			resource.Dependencies, func(dep parser.MultyResourceDependency) []string {
				if !dep.UserDeclared {
					return nil
				}
				return []string{dep.To.ID}
			},
		),
	}
	body := hclutil.PartialDecode(resource.HCLBody, ctx, commonParams)

	r, err := InitResource(resource.Type, commonParams)
	if err != nil {
		validate.LogFatalWithSourceRange(resource.DefinitionRange, err.Error())
	}

	body = mhclProcessor.Process(body, r, ctx, cloud)

	validationInfo, diags := decodeBody(r, body, ctx, resource.DefinitionRange, resource.ID)
	commonParams.ResourceValidationInfo = validationInfo

	return r, diags
}

func decodeBody(t any, body hcl.Body, ctx *hcl.EvalContext, definitionRange hcl.Range, resourceId string) (*validate.ResourceValidationInfo, hcl.Diagnostics) {
	schema, _ := gohcl.ImpliedBodySchema(t)
	content, diags := body.Content(schema)
	if diags != nil {
		return nil, diags
	}
	diags = gohcl.DecodeBody(body, ctx, t)
	return validate.NewResourceValidationInfoFromContent(content, definitionRange, resourceId), diags
}

func GetRgId(defaultRgId hcl.Expression, attrs hcl.Attributes, ctx *hcl.EvalContext, resourceType string) (string, bool) {
	// TODO: refactor so that we don't do explicit checks here
	if resourceType == "resource_group" || resourceType == "route_table_association" {
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
	childCtx.Variables["resource_type"] = cty.StringVal(common.GetResourceTypeAbbreviation(resourceType))
	childCtx.Variables["rg_vars"] = rgVars
	diags = gohcl.DecodeExpression(rgIdExpr, ctx, &rgId)
	if diags != nil {
		validate.LogFatalWithDiags(
			diags, "Unable to resolve resource group name for resource. "+
				"Are you missing a rg_var?",
		)
	}
	return rgId, shouldCreateRg
}
