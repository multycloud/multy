package types

import (
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output/route_table_association"
	"multy-go/validate"
)

// route_table_association
type RouteTableAssociation struct {
	*resources.CommonResourceParams
	SubnetId     string `hcl:"subnet_id"`
	RouteTableId string `hcl:"route_table_id"`
}

func (r *RouteTableAssociation) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []any {
	if cloud == common.AWS {
		return []any{
			route_table_association.AwsRouteTableAssociation{
				AwsResource: common.AwsResource{
					ResourceName: route_table_association.AwsResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
				},
				RouteTableId: r.RouteTableId,
				SubnetId:     r.SubnetId,
			},
		}
	} else if cloud == common.AZURE {
		return []any{
			route_table_association.AzureRouteTableAssociation{
				AzResource: common.AzResource{
					ResourceName: route_table_association.AzureResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
				},
				RouteTableId: r.RouteTableId,
				SubnetId:     r.SubnetId,
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *RouteTableAssociation) Validate(ctx resources.MultyContext) {
	return
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
