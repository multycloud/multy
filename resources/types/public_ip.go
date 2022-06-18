package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

/*
Notes
AWS: NIC association done on public_ip
Azure: NIC association done on NIC creation
*/

type PublicIp struct {
	resources.ResourceWithId[*resourcespb.PublicIpArgs]
}

func (r *PublicIp) Create(resourceId string, args *resourcespb.PublicIpArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		// todo: maybe put in the same RG as VM?
		rgId, err := NewRg("pip", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewPublicIp(r, resourceId, args)
}

func (r *PublicIp) Update(args *resourcespb.PublicIpArgs, _ *resources.Resources) error {
	return NewPublicIp(r, r.ResourceId, args)
}

func (r *PublicIp) Import(resourceId string, args *resourcespb.PublicIpArgs, _ *resources.Resources) error {
	return NewPublicIp(r, resourceId, args)
}

func (r *PublicIp) Export(_ *resources.Resources) (*resourcespb.PublicIpArgs, bool, error) {
	return r.Args, true, nil
}

func NewPublicIp(r *PublicIp, resourceId string, args *resourcespb.PublicIpArgs) error {
	r.ResourceWithId = resources.ResourceWithId[*resourcespb.PublicIpArgs]{
		ResourceId: resourceId,
		Args:       args,
	}
	return nil
}

func (r *PublicIp) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}
