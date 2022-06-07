package network_interface

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_network_interface"
const AwsEipAssociationResourceName = "aws_eip_association"

type AwsNetworkInterface struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_network_interface"`
	SubnetId            string `hcl:"subnet_id,expr"`
}

type AwsEipAssociation struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_eip_association"`
	AllocationId        string `hcl:"allocation_id,expr"`
	NetworkInterfaceId  string `hcl:"network_interface_id,expr"`
}
