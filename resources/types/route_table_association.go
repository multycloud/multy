package types

import (
	"fmt"
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

func (r *RouteTableAssociation) Translate(cloud common.CloudProvider, ctx resources.MultyContext) ([]output.TfBlock, error) {
	rtId, err := resources.GetMainOutputId(r.RouteTableId, cloud)
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(r.SubnetId, cloud)
	if err != nil {
		return nil, err
	}
	if cloud == common.AWS {
		return []output.TfBlock{
			route_table_association.AwsRouteTableAssociation{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				RouteTableId: rtId,
				SubnetId:     subnetId,
			},
		}, nil
	} else if cloud == common.AZURE {
		return []output.TfBlock{
			route_table_association.AzureRouteTableAssociation{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)},
				},
				RouteTableId: rtId,
				SubnetId:     subnetId,
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", cloud)
}

func (r *RouteTableAssociation) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	return errs
}

func (r *RouteTableAssociation) GetMainResourceName(cloud common.CloudProvider) (string, error) {
	switch cloud {
	case common.AWS:
		return route_table_association.AwsResourceName, nil
	case common.AZURE:
		return route_table_association.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", cloud)
	}
}
