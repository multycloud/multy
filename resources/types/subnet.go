package types

import (
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

func (r *Subnet) Update(args *resourcespb.SubnetArgs, others *resources.Resources) error {
	return NewSubnet(r, r.ResourceId, args, others)
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
	if err := validate.NewWordWithDotHyphenUnder80Validator().Check(r.Args.Name, r.ResourceId); err != nil {
		errs = append(errs, r.NewValidationError(err, "name"))
	}
	if err := validate.NewCIDRIPv4Check().Check(r.Args.CidrBlock, r.ResourceId); err != nil {
		errs = append(errs, r.NewValidationError(err, "cidr_block"))
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
