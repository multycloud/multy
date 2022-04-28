package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/validate"
)

// route_table_association
type RouteTableAssociation struct {
	resources.ChildResourceWithId[*RouteTable, *resourcespb.RouteTableAssociationArgs]

	RouteTable *RouteTable
	Subnet     *Subnet
}

func NewRouteTableAssociation(resourceId string, args *resourcespb.RouteTableAssociationArgs, others resources.Resources) (*RouteTableAssociation, error) {
	rta := &RouteTableAssociation{
		ChildResourceWithId: resources.ChildResourceWithId[*RouteTable, *resourcespb.RouteTableAssociationArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}
	rt, err := resources.Get[*RouteTable](resourceId, others, args.RouteTableId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{rta.NewValidationError(err, "virtual_network_id")})
	}
	rta.Parent = rt
	rta.RouteTable = rt

	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{rta.NewValidationError(err, "subnet_id")})
	}
	rta.Subnet = subnet
	return rta, nil
}

func (r *RouteTableAssociation) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	rtId, err := resources.GetMainOutputId(r.RouteTable)
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(r.Subnet)
	if err != nil {
		return nil, err
	}
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{
			route_table_association.AwsRouteTableAssociation{
				AwsResource: &common.AwsResource{
					TerraformResource: output.TerraformResource{ResourceId: r.Subnet.ResourceId},
				},
				RouteTableId: rtId,
				SubnetId:     subnetId,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return []output.TfBlock{
			route_table_association.AzureRouteTableAssociation{
				// Here we use the subnet id so that it is the same as the one that is created by default in the subnet.
				// This ensures that if a RTA is created after the default RTA is created, they will have the same ID and
				// terraform will either update it in place or destroy it before creating it.
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.Subnet.ResourceId},
				},
				RouteTableId: rtId,
				SubnetId:     subnetId,
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *RouteTableAssociation) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	return errs
}

func (r *RouteTableAssociation) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return route_table_association.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return route_table_association.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
