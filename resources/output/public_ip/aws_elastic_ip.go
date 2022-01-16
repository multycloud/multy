package public_ip

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/validate"
)

const AwsResourceName = "aws_eip"

//  EIP may require IGW to exist prior to association. Use depends_on to set an explicit dependency on the IGW.
type AwsElasticIp struct {
	common.AwsResource `hcl:",squash"`
	InstanceId         string `hcl:"instance" hcle:"omitempty"`
	Vpc                bool   `hcl:"vpc,optional" hcle:"omitempty"`
	NetworkInterfaceId string `hcl:"network_interface" hcle:"omitempty"`
}

func (eIp AwsElasticIp) GetId(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("%s.%s.id", AwsResourceName, eIp.ResourceId)
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}
