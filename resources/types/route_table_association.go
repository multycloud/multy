package types

import (
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/validate"
)

// route_table_association
type RouteTableAssociation struct {
	*resources.CommonResourceParams
	SubnetId     *Subnet     `mhcl:"ref=subnet_id"`
	RouteTableId *RouteTable `mhcl:"ref=route_table_id"`
}

func (r *RouteTableAssociation) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		return []output.TfBlock{
			route_table_association.AwsRouteTableAssociation{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				RouteTableId: resources.GetMainOutputId(r.RouteTableId, cloud),
				SubnetId:     resources.GetMainOutputId(r.SubnetId, cloud),
			},
		}
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			route_table_association.AzureRouteTableAssociation{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				RouteTableId: resources.GetMainOutputId(r.RouteTableId, cloud),
				SubnetId:     resources.GetMainOutputId(r.SubnetId, cloud),
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *RouteTableAssociation) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *RouteTableAssociation) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return route_table_association.AwsResourceName
	case common.AZURE:
		return route_table_association.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
