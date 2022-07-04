package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type NetworkInterface struct {
	resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]

	Subnet   *Subnet
	PublicIp *PublicIp
}

func (r *NetworkInterface) Create(resourceId string, args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return err
		}
		rgId, err := NewRgFromParent("nic", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	if args.AvailabilityZone == 0 {
		args.AvailabilityZone = 1
	}

	return NewNetworkInterface(r, resourceId, args, others)
}

func (r *NetworkInterface) Update(args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) error {
	return NewNetworkInterface(r, r.ResourceId, args, others)
}

func (r *NetworkInterface) Import(resourceId string, args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) error {
	return NewNetworkInterface(r, resourceId, args, others)
}

func (r *NetworkInterface) Export(_ *resources.Resources) (*resourcespb.NetworkInterfaceArgs, bool, error) {
	return r.Args, true, nil
}
func NewNetworkInterface(r *NetworkInterface, resourceId string, args *resourcespb.NetworkInterfaceArgs, others *resources.Resources) error {
	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return err
	}
	r.Subnet = subnet
	pIp, err := resources.GetOptional[*PublicIp](resourceId, others, args.PublicIpId)
	if err != nil {
		return err
	}
	r.PublicIp = pIp

	r.ResourceWithId = resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]{
		ResourceId: resourceId,
		Args:       args,
	}
	return nil
}

func (r *NetworkInterface) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}
