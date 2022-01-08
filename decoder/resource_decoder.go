package decoder

import (
	"fmt"
	"multy-go/hclutil"
	"multy-go/parser"
	"multy-go/resources"
	"multy-go/resources/common"
	rg "multy-go/resources/resource_group"
	"multy-go/resources/types"
	"multy-go/validate"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

type ResourceDecoder struct {
	globalConfig DecodedGlobalConfig
}

func (d ResourceDecoder) Decode(resource parser.MultyResource, ctx CloudSpecificContext) ([]resources.CloudSpecificResource, CloudSpecificContext) {
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
		r, diags := decode(resource, cloudCtx, rgId, attrs)
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
		outputVars := r.GetOutputValues(cloud)
		for key, val := range outputVars {
			allVars[key] = val
		}
		allVars["id"] = cty.StringVal(cloudSpecificR.GetResourceId())
		resultCtx.AddCloudDependentVars(map[string]cty.Value{
			resource.ID: cty.ObjectVal(allVars),
		}, cloud)

		if shouldCreateRg {
			result = append(result, resources.CloudSpecificResource{ImplicitlyCreated: true, Cloud: cloud, Resource: rg.DefaultResourceGroup(rgId)})
		}
	}
	return result, resultCtx
}

func decode(resource parser.MultyResource, ctx *hcl.EvalContext, rgId string, attrs hcl.Attributes) (resources.Resource, hcl.Diagnostics) {
	var r resources.Resource
	var diags hcl.Diagnostics
	if resource.ID == "" {
		validate.LogFatalWithSourceRange(resource.DefinitionRange, "found resource of type '%s' with no id", resource.Type)
	}

	validationInfo := validate.NewResourceValidationInfo(attrs)
	commonParams := &resources.CommonResourceParams{
		ResourceValidationInfo: validationInfo,
		ResourceId:             resource.ID,
		ResourceGroupId:        rgId,
	}

	body := hclutil.PartialDecode(resource.HCLBody, ctx, commonParams)

	if resource.Type == "virtual_network" {
		t := types.VirtualNetwork{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "subnet" {
		t := types.Subnet{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "network_security_group" {
		t := types.NetworkSecurityGroup{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "virtual_machine" {
		t := types.VirtualMachine{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "public_ip" {
		t := types.PublicIp{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "route_table" {
		t := types.RouteTable{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "route_table_association" {
		t := types.RouteTableAssociation{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "network_interface" {
		t := types.NetworkInterface{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "database" {
		t := types.Database{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "resource_group" {
		t := rg.Type{ResourceId: resource.ID}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "object_storage" {
		t := types.ObjectStorage{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
	} else if resource.Type == "object_storage_object" {
		t := types.ObjectStorageObject{CommonResourceParams: commonParams}
		diags = gohcl.DecodeBody(body, ctx, &t)
		r = &t
		//} else if resource.Type == "vault" {
		//	t := types.Vault{CommonResourceParams: commonParams}
		//	diags = gohcl.DecodeBody(body, ctx, &t)
		//	r = &t
		//} else if resource.Type == "vault_secret" {
		//	t := types.VaultSecret{CommonResourceParams: commonParams}
		//	diags = gohcl.DecodeBody(body, ctx, &t)
		//	r = &t
		//} else if resource.Type == "lambda" {
		//	t := types.Lambda{CommonResourceParams: commonParams}
		//	diags = gohcl.DecodeBody(body, ctx, &t)
		//	r = &t
	} else {
		validate.LogFatalWithSourceRange(resource.DefinitionRange, "unknown resource type %s", resource.Type)
	}

	return r, diags
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
