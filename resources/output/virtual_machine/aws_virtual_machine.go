package virtual_machine

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_instance"

type AwsEC2 struct {
	*common.AwsResource      `hcl:",squash"  default:"name=aws_instance"`
	Ami                      string                   `hcl:"ami"`
	InstanceType             string                   `hcl:"instance_type"`
	AssociatePublicIpAddress bool                     `hcl:"associate_public_ip_address" hcle:"omitempty"`
	SubnetId                 string                   `hcl:"subnet_id,expr" hcle:"omitempty"`
	UserDataBase64           string                   `hcl:"user_data_base64" hcle:"omitempty"`
	NetworkInterfaces        []AwsEc2NetworkInterface `hcl:"network_interface,blocks" hcle:"omitempty"`
	SecurityGroupIds         []string                 `hcl:"vpc_security_group_ids,expr" hcle:"omitempty"`
	KeyName                  string                   `hcl:"key_name,expr" hcle:"omitempty"`
	IamInstanceProfile       string                   `hcl:"iam_instance_profile,expr" hcle:"omitempty"`
}

type AwsEc2NetworkInterface struct {
	NetworkInterfaceId string `hcl:"network_interface_id,expr"`
	DeviceIndex        int    `hcl:"device_index"`
}

const AwsKeyPairResourceName = "aws_key_pair"

type AwsKeyPair struct {
	*common.AwsResource `hcl:",squash"  default:"name=aws_key_pair"`
	KeyName             string `hcl:"key_name"`
	PublicKey           string `hcl:"public_key"`
}
