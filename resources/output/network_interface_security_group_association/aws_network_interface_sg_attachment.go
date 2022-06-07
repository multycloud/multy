package network_interface_security_group_association

import "github.com/multycloud/multy/resources/common"

const AwsResourceName = "aws_network_interface_sg_attachment"

type AwsNetworkInterfaceSecurityGroupAssociation struct {
	*common.AwsResource    `hcl:",squash" default:"name=aws_network_interface_sg_attachment"`
	NetworkInterfaceId     string `hcl:"network_interface_id,expr"`
	NetworkSecurityGroupId string `hcl:"security_group_id,expr"`
}
