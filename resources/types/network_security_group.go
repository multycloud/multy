package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

/*
Notes
NSG can only be applied to NIC (currently done in VM creation, to be changed to separate resource)
When NSG is applied, only rules specified are allowed.
AWS: VPC traffic is always added as an extra rule
*/

type NetworkSecurityGroup struct {
	resources.ResourceWithId[*resourcespb.NetworkSecurityGroupArgs]

	VirtualNetwork *VirtualNetwork
}

func (r *NetworkSecurityGroup) Create(resourceId string, args *resourcespb.NetworkSecurityGroupArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
		if err != nil {
			return err
		}
		rgId, err := NewRgFromParent("nsg", vn.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewNetworkSecurityGroup(r, resourceId, args, others)
}

func (r *NetworkSecurityGroup) Update(args *resourcespb.NetworkSecurityGroupArgs, others *resources.Resources) error {
	return NewNetworkSecurityGroup(r, r.ResourceId, args, others)
}

func (r *NetworkSecurityGroup) Import(resourceId string, args *resourcespb.NetworkSecurityGroupArgs, others *resources.Resources) error {
	return NewNetworkSecurityGroup(r, resourceId, args, others)
}

func (r *NetworkSecurityGroup) Export(_ *resources.Resources) (*resourcespb.NetworkSecurityGroupArgs, bool, error) {
	return r.Args, true, nil
}

type RuleType struct {
	Protocol  string `cty:"protocol"`
	Priority  int    `cty:"priority"`
	FromPort  string `cty:"from_port"`
	ToPort    string `cty:"to_port"`
	CidrBlock string `cty:"cidr_block"`
	Direction string `cty:"direction"`
}

func NewNetworkSecurityGroup(nsg *NetworkSecurityGroup, resourceId string, args *resourcespb.NetworkSecurityGroupArgs, others *resources.Resources) error {
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
	if err != nil {
		return err
	}
	nsg.VirtualNetwork = vn
	nsg.ResourceWithId = resources.ResourceWithId[*resourcespb.NetworkSecurityGroupArgs]{
		ResourceId: resourceId,
		Args:       args,
	}
	return nil
}

func (r *NetworkSecurityGroup) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if err := validate.NewWordWithDotHyphenUnder80Validator().Check(r.Args.Name, r.ResourceId); err != nil {
		errs = append(errs, r.NewValidationError(err, "name"))
	}
	cidrValidator := validate.NewCIDRIPv4Check()
	portValidator := validate.NewPortCheck()
	protoValidator := validate.NewProtocolCheck()
	for _, rule := range r.Args.Rules {
		if err := portValidator.Check(rule.PortRange.To, "to_port"); err != nil {
			errs = append(errs, r.NewValidationError(err, "rules"))
		}
		if err := portValidator.Check(rule.PortRange.From, "from_port"); err != nil {
			errs = append(errs, r.NewValidationError(err, "rules"))
		}
		if err := cidrValidator.Check(rule.CidrBlock, "rule cidr"); err != nil {
			errs = append(errs, r.NewValidationError(err, "rules"))
		}
		if err := protoValidator.Check(rule.Protocol, "rule protocol"); err != nil {
			errs = append(errs, r.NewValidationError(err, "rules"))
		}
	}
	// TODO validate location matches with VN location
	return errs
}
