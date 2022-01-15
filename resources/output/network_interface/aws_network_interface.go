package network_interface

import "multy-go/resources/common"

const AwsResourceName = "aws_network_interface"

type AwsNetworkInterface struct {
	common.AwsResource `hcl:",squash"`
	SubnetId           string `hcl:"subnet_id"`
}
