package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type RouteTableAssociation struct {
	resources.ChildResourceWithId[*RouteTable, *resourcespb.RouteTableAssociationArgs]

	RouteTable *RouteTable
	Subnet     *Subnet
}

func (r *RouteTableAssociation) Create(resourceId string, args *resourcespb.RouteTableAssociationArgs, others *resources.Resources) error {
	return NewRouteTableAssociation(r, resourceId, args, others)
}

func (r *RouteTableAssociation) Update(args *resourcespb.RouteTableAssociationArgs, _ *resources.Resources) error {
	r.Args = args
	return nil
}

func (r *RouteTableAssociation) Import(resourceId string, args *resourcespb.RouteTableAssociationArgs, others *resources.Resources) error {
	return NewRouteTableAssociation(r, resourceId, args, others)
}

func (r *RouteTableAssociation) Export(_ *resources.Resources) (*resourcespb.RouteTableAssociationArgs, bool, error) {
	return r.Args, true, nil
}

func NewRouteTableAssociation(rta *RouteTableAssociation, resourceId string, args *resourcespb.RouteTableAssociationArgs, others *resources.Resources) error {
	rt, err := resources.Get[*RouteTable](resourceId, others, args.RouteTableId)
	if err != nil {
		return errors.ValidationErrors([]validate.ValidationError{rta.NewValidationError(err, "virtual_network_id")})
	}
	rta.ChildResourceWithId = resources.NewChildResource(resourceId, rt, args)
	rta.Parent = rt
	rta.RouteTable = rt

	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return errors.ValidationErrors([]validate.ValidationError{rta.NewValidationError(err, "subnet_id")})
	}
	rta.Subnet = subnet
	return nil
}
func (r *RouteTableAssociation) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if r.RouteTable.VirtualNetwork.ResourceId != r.Subnet.VirtualNetwork.ResourceId {
		errs = append(errs, r.NewValidationError(fmt.Errorf(
			"cannot associate subnet %s to route_table %s because they are in different virtual networks",
			r.Subnet.ResourceId, r.RouteTable.ResourceId),
			"subnet_id"))
	}
	return errs
}
