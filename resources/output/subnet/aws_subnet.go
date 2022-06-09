package subnet

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_subnet"

type AwsSubnet struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_subnet"`
	CidrBlock           string `hcl:"cidr_block" json:"cidr_block"`
	VpcId               string `hcl:"vpc_id,expr" json:"vpc_id"`
	AvailabilityZone    string `hcl:"availability_zone,optional" hcle:"omitempty" json:"availability_zone"`
	MapPublicIpOnLaunch bool   `hcl:"map_public_ip_on_launch"  hcle:"omitempty" json:"map_public_ip_on_launch"`
}
