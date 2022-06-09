package types

import (
	"fmt"

	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

/*
Notes:
AWS: Internet route to IGW
Azure: Internet route nextHop Internet
*/

type RouteTable struct {
	resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.RouteTableArgs]

	VirtualNetwork *VirtualNetwork `mhcl:"ref=virtual_network"`
}

func (r *RouteTable) Create(resourceId string, args *resourcespb.RouteTableArgs, others *resources.Resources) error {
	return NewRouteTable(r, resourceId, args, others)
}

func (r *RouteTable) Update(args *resourcespb.RouteTableArgs, _ *resources.Resources) error {
	r.Args = args
	return nil
}

func (r *RouteTable) Import(resourceId string, args *resourcespb.RouteTableArgs, others *resources.Resources) error {
	return NewRouteTable(r, resourceId, args, others)
}

func (r *RouteTable) Export(_ *resources.Resources) (*resourcespb.RouteTableArgs, bool, error) {
	return r.Args, true, nil
}

func NewRouteTable(rt *RouteTable, resourceId string, args *resourcespb.RouteTableArgs, others *resources.Resources) error {
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, rt.ResourceId, "virtual_network_id"))
	}
	rt.ChildResourceWithId = resources.NewChildResource(resourceId, vn, args)
	rt.Parent = vn
	rt.VirtualNetwork = vn
	return nil
}
func (r *RouteTable) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if len(r.Args.Routes) > 20 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("\"%d\" exceeds routes limit is 20", len(r.Args.Routes)), "routes"))
	}
	for _, route := range r.Args.Routes {
		if route.Destination == resourcespb.RouteDestination_UNKNOWN_DESTINATION {
			errs = append(errs, r.NewValidationError(fmt.Errorf("unknown route destination"), "route"))
		}
		//	if route.CidrBlock valid CIDR
	}
	return errs
}
