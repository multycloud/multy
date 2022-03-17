package virtual_network

import (
	"fmt"
	"multy/resources/common"
)

const AwsResourceName = "aws_vpc"
const AwsInternetGatewayName = "aws_internet_gateway"

type AwsVpc struct {
	*common.AwsResource `hcl:",squash"  default:"name=aws_vpc"`
	CidrBlock           string `hcl:"cidr_block"`
	EnableDnsHostnames  bool   `hcl:"enable_dns_hostnames"` // needed for publicly accessible rds
}

// AwsInternetGateway : by default, Internet Gateway is associated with VPC
type AwsInternetGateway struct {
	*common.AwsResource `hcl:",squash"  default:"name=aws_internet_gateway"`
	VpcId               string `hcl:"vpc_id,expr"`
}

func (vpc *AwsVpc) GetDefaultRouteTableId() string {
	return fmt.Sprintf("aws_vpc.%s.default_route_table_id", vpc.ResourceId)
}

func (igw *AwsInternetGateway) GetId() string {
	return fmt.Sprintf("aws_internet_gateway.%s.id", igw.ResourceId)
}
