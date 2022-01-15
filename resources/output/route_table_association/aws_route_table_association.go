package route_table_association

import "multy-go/resources/common"

const AwsResourceName = "aws_route_table_association"

type AwsRouteTableAssociation struct {
	common.AwsResource `hcl:",squash"`
	SubnetId           string `hcl:"subnet_id"`
	RouteTableId       string `hcl:"route_table_id,expr"`
}
