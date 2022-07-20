package public_ip

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_eip"

//  EIP may require IGW to exist prior to association. Use depends_on to set an explicit dependency on the IGW.
type AwsElasticIp struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_eip"`
	InstanceId          string `hcl:"instance" hcle:"omitempty"`
	Vpc                 bool   `hcl:"vpc,optional" hcle:"omitempty"`
	NetworkInterfaceId  string `hcl:"network_interface,expr" hcle:"omitempty"`

	PublicIp string `json:"public_ip" hcle:"omitempty"`
}
