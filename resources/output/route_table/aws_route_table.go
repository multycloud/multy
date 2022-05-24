package route_table

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_route_table"

type AwsRouteTable struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_route_table"`
	VpcId               string               `hcl:"vpc_id,expr" json:"vpc_id,omitempty"`
	Routes              []AwsRouteTableRoute `hcl:"route,blocks" json:"route,omitempty"`
}

type AwsDefaultRouteTable struct {
	*common.AwsResource `hcl:",squash"`
	DefaultRouteTableId string               `hcl:"default_route_table_id,expr"`
	Routes              []AwsRouteTableRoute `hcl:"route,blocks"`
}

type AwsRouteTableRoute struct {
	CidrBlock string `hcl:"cidr_block"`
	GatewayId string `hcl:"gateway_id,expr"`
}
