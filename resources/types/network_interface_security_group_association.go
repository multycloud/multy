package types

import (
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
)

type NetworkInterfaceSecurityGroupAssociation struct {
	resources.ChildResourceWithId[*NetworkInterface, *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs]

	NetworkInterface     *NetworkInterface
	NetworkSecurityGroup *NetworkSecurityGroup
}

func (r *NetworkInterfaceSecurityGroupAssociation) Create(resourceId string, args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) error {
	return NewNetworkInterfaceSecurityGroupAssociation(r, resourceId, args, others)
}

func (r *NetworkInterfaceSecurityGroupAssociation) Update(args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) error {
	return NewNetworkInterfaceSecurityGroupAssociation(r, r.ResourceId, args, others)
}

func (r *NetworkInterfaceSecurityGroupAssociation) Import(resourceId string, args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) error {
	return NewNetworkInterfaceSecurityGroupAssociation(r, resourceId, args, others)
}

func (r *NetworkInterfaceSecurityGroupAssociation) Export(_ *resources.Resources) (*resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, bool, error) {
	return r.Args, true, nil
}

func NewNetworkInterfaceSecurityGroupAssociation(r *NetworkInterfaceSecurityGroupAssociation, resourceId string, args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs, others *resources.Resources) error {
	nic, err := resources.Get[*NetworkInterface](resourceId, others, args.NetworkInterfaceId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, r.ResourceId, "network_interface_id"))
	}
	r.ChildResourceWithId = resources.NewChildResource(resourceId, nic, args)
	r.NetworkInterface = nic

	nsg, err := resources.Get[*NetworkSecurityGroup](resourceId, others, args.SecurityGroupId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, r.ResourceId, "security_group_id"))
	}
	r.NetworkSecurityGroup = nsg
	return nil
}

func (r *NetworkInterfaceSecurityGroupAssociation) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	return nil
}

func (r *NetworkInterfaceSecurityGroupAssociation) ParseCloud(args *resourcespb.NetworkInterfaceSecurityGroupAssociationArgs) commonpb.CloudProvider {
	return common.ParseCloudFromResourceId(args.NetworkInterfaceId)
}
