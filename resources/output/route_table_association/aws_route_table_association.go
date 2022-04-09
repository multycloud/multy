package route_table_association

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_route_table_association"

type AwsRouteTableAssociation struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_route_table_association"`
	SubnetId            string `hcl:"subnet_id,expr"`
	RouteTableId        string `hcl:"route_table_id,expr"`
}
