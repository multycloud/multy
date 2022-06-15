package types

import (
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
	"net"
)

/*
Notes:
Azure: New subnets will be associated with a default route table to block traffic to internet
*/

type Subnet struct {
	resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.SubnetArgs]

	VirtualNetwork *VirtualNetwork
}

func (r *Subnet) Create(resourceId string, args *resourcespb.SubnetArgs, others *resources.Resources) error {
	return NewSubnet(r, resourceId, args, others)
}

func (r *Subnet) Update(args *resourcespb.SubnetArgs, _ *resources.Resources) error {
	r.Args = args
	return nil
}

func (r *Subnet) Import(resourceId string, args *resourcespb.SubnetArgs, others *resources.Resources) error {
	return NewSubnet(r, resourceId, args, others)
}

func (r *Subnet) Export(_ *resources.Resources) (*resourcespb.SubnetArgs, bool, error) {
	return r.Args, true, nil
}

func NewSubnet(s *Subnet, resourceId string, subnet *resourcespb.SubnetArgs, others *resources.Resources) error {
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, subnet.VirtualNetworkId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, s.ResourceId, "virtual_network_id"))
	}
	s.ChildResourceWithId = resources.NewChildResource(resourceId, vn, subnet)
	s.Parent = vn
	s.VirtualNetwork = vn
	return nil
}

func (r *Subnet) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.AvailbilityZone valid { return false }
	if len(r.Args.CidrBlock) == 0 { // max len?
		errs = append(errs, r.NewValidationError(fmt.Errorf("%s cidr_block length is invalid", r.ResourceId), "cidr_block"))
	}

	if _, vNetBlock, err := net.ParseCIDR(r.Args.CidrBlock); err == nil {
		if _, subnetBlock, err := net.ParseCIDR(r.Args.CidrBlock); err != nil {
			errs = append(errs, validate.ValidationError{
				ErrorMessage: err.Error(),
				ResourceId:   r.ResourceId,
				FieldName:    "cidr_block",
			})
		} else if err := cidr.VerifyNoOverlap([]*net.IPNet{subnetBlock}, vNetBlock); err != nil {
			errs = append(errs, validate.ValidationError{
				ErrorMessage: err.Error(),
				ResourceId:   r.ResourceId,
				FieldName:    "cidr_block",
			})
		}
	}

	return errs
}
