package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

/*
Notes:
AWS: Can pass NICs and which overrides public_ip and subnet_id
Azure: To have a private IP by default, if no NIC is passed, one will be created.
       For PublicIp to be auto_assigned, public_ip is created an attached to default NIC
 	   NSG_NIC association
*/

type VirtualMachine struct {
	resources.ResourceWithId[*resourcespb.VirtualMachineArgs]

	NetworkInterface      []*NetworkInterface
	NetworkSecurityGroups []*NetworkSecurityGroup
	Subnet                *Subnet
	PublicIp              *PublicIp
}

func (r *VirtualMachine) Create(resourceId string, args *resourcespb.VirtualMachineArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return err
		}
		rgId, err := NewRgFromParent("vm", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	if args.AvailabilityZone == 0 {
		args.AvailabilityZone = 1
		// if there's a network interface attached, place vm in same zone
		if len(args.NetworkInterfaceIds) > 0 {
			if ni, err := resources.Get[*NetworkInterface](resourceId, others, args.NetworkInterfaceIds[0]); err == nil {
				args.AvailabilityZone = ni.Args.AvailabilityZone
			}
		}
	}
	return NewVirtualMachine(r, resourceId, args, others)
}

func (r *VirtualMachine) Update(args *resourcespb.VirtualMachineArgs, others *resources.Resources) error {
	return NewVirtualMachine(r, r.ResourceId, args, others)
}
func (r *VirtualMachine) Import(resourceId string, args *resourcespb.VirtualMachineArgs, others *resources.Resources) error {
	return NewVirtualMachine(r, resourceId, args, others)
}

func (r *VirtualMachine) Export(_ *resources.Resources) (*resourcespb.VirtualMachineArgs, bool, error) {
	return r.Args, true, nil
}

func NewVirtualMachine(vm *VirtualMachine, resourceId string, args *resourcespb.VirtualMachineArgs, others *resources.Resources) error {
	networkInterfaces, err := util.MapSliceValuesErr(args.NetworkInterfaceIds, func(id string) (*NetworkInterface, error) {
		return resources.Get[*NetworkInterface](resourceId, others, id)
	})
	if err != nil {
		return err
	}
	vm.NetworkInterface = networkInterfaces

	networkSecurityGroups, err := util.MapSliceValuesErr(args.NetworkSecurityGroupIds, func(id string) (*NetworkSecurityGroup, error) {
		return resources.Get[*NetworkSecurityGroup](resourceId, others, id)
	})
	if err != nil {
		return err
	}
	vm.NetworkSecurityGroups = networkSecurityGroups

	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return err
	}
	vm.Subnet = subnet

	publicIp, err := resources.GetOptional[*PublicIp](resourceId, others, args.PublicIpId)
	if err != nil {
		return err
	}
	vm.PublicIp = publicIp

	if args.GetImageReference() == nil {
		args.ImageReference = &resourcespb.ImageReference{
			Os:      resourcespb.ImageReference_UBUNTU,
			Version: "18.04",
		}
	}

	vm.ResourceWithId = resources.ResourceWithId[*resourcespb.VirtualMachineArgs]{
		ResourceId: resourceId,
		Args:       args,
	}

	return nil
}

func (r *VirtualMachine) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.Size valid { return false }
	if r.Args.GeneratePublicIp && len(r.NetworkInterface) != 0 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("generate public ip can't be set with network interface ids"), "generate_public_ip"))
	}
	if r.Args.GeneratePublicIp && r.PublicIp != nil {
		errs = append(errs, r.NewValidationError(fmt.Errorf("conflict between generate_public_ip and public_ip_id"), "generate_public_ip"))
	}
	return errs
}

func (r *VirtualMachine) GetAwsIdentity() string {
	return fmt.Sprintf("%s-vm-role", r.ResourceId)
}
