package subnet

import (
	"multy-go/resources/common"
)

const AwsResourceName = "aws_subnet"

type AwsSubnet struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_subnet"`
	CidrBlock           string `hcl:"cidr_block"`
	VpcId               string `hcl:"vpc_id,expr"`
	AvailabilityZone    string `hcl:"availability_zone,optional" hcle:"omitempty"`
	MapPublicIpOnLaunch bool   `hcl:"map_public_ip_on_launch"  hcle:"omitempty"`
}
