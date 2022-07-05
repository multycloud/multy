package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
	"net"
)

/*
Virtual network traffic is only internal
AWS: Default security group defaults to allow all traffic to  mirror Azure
	IGW created by default (to be changed)
AZ: Route table created to restrict traffic on vnet

*/

type VirtualNetwork struct {
	resources.ResourceWithId[*resourcespb.VirtualNetworkArgs]
}

func (r *VirtualNetwork) Create(resourceId string, args *resourcespb.VirtualNetworkArgs, others *resources.Resources) error {
	if args.GetCommonParameters().GetResourceGroupId() == "" {
		rgId, err := NewRg("vn", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewVirtualNetwork(r, resourceId, args)
}

func (r *VirtualNetwork) Update(args *resourcespb.VirtualNetworkArgs, _ *resources.Resources) error {
	return NewVirtualNetwork(r, r.ResourceId, args)
}

func (r *VirtualNetwork) Import(resourceId string, args *resourcespb.VirtualNetworkArgs, _ *resources.Resources) error {
	return NewVirtualNetwork(r, resourceId, args)
}

func (r *VirtualNetwork) Export(_ *resources.Resources) (*resourcespb.VirtualNetworkArgs, bool, error) {
	return r.Args, true, nil
}

func NewVirtualNetwork(r *VirtualNetwork, resourceId string, vn *resourcespb.VirtualNetworkArgs) error {
	r.ResourceWithId = resources.NewResource(resourceId, vn)
	return nil
}

// TODO validate commonparams
func (r *VirtualNetwork) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	//if r.Name contains not letters,numbers,_,- { return false }
	//if r.Name length? { return false }
	//if r.CidrBlock valid CIDR { return false }
	if len(r.Args.CidrBlock) == 0 { // max len?
		errs = append(errs, validate.ValidationError{
			ErrorMessage: "cidr_block length is invalid",
			ResourceId:   r.ResourceId,
			FieldName:    "cidr_block",
		})
	}
	if _, _, err := net.ParseCIDR(r.Args.CidrBlock); err != nil {
		errs = append(errs, validate.ValidationError{
			ErrorMessage: err.Error(),
			ResourceId:   r.ResourceId,
			FieldName:    "cidr_block",
		})
	}
	return errs
}

/*
Virtual network is a private address space when resources can be placed.

By default, resources inside `virtual_network` cannot access the internet. To enable internet access look at`route_table`

Mapping:
AWS:
`aws_vpc` - VPC
`aws_internet_gateway` - Internet Gateway attached to VPC. Default route table does not route traffic to outside VPC

Azure:
`azurerm_virtual_network` - Vnet
`azurerm_route_table` - Default route table to block internet access

Inputs:

cidr_block `[]string` (required): Address range of virtual network.
name `string` (required): Name of virtual network. Read about multy naming [here](xxx)
location `string` (optional): Region to deploy resource into. See available locations [here](xxx)
clouds `[]string` (optional): Clouds where services will be deployed to. Learn about multy clouds [here](xxx)
*/
