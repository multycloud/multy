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

func (r *RouteTableAssociation) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	var rt *RouteTable
	if s, err := ctx.GetResource(r.RouteTableId); err != nil {
		r.LogFatal(r.ResourceId, "route_table_id", err.Error())
	} else {
		rt = s.Resource.(*RouteTable)
	}

	if cloud == common.AWS {
		return []interface{}{
			route_table_association.AwsRouteTableAssociation{
				AwsResource: common.AwsResource{
					ResourceName: route_table_association.AwsResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
				},
				RouteTableId: rt.GetId(cloud),
				SubnetId:     r.SubnetId,
			},
		}
	} else if cloud == common.AZURE {
		return []interface{}{
			route_table_association.AzureRouteTableAssociation{
				AzResource: common.AzResource{
					ResourceName: route_table_association.AzureResourceName,
					ResourceId:   r.GetTfResourceId(cloud),
				},
				RouteTableId: rt.GetId(cloud),
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
