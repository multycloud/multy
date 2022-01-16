package virtual_machine

import (
	"multy-go/resources/common"
)

const AwsResourceName = "aws_instance"

type AwsEC2 struct {
	common.AwsResource       `hcl:",squash"`
	Ami                      string                   `hcl:"ami"`
	InstanceType             string                   `hcl:"instance_type"`
	AssociatePublicIpAddress bool                     `hcl:"associate_public_ip_address" hcle:"omitempty"`
	SubnetId                 string                   `hcl:"subnet_id" hcle:"omitempty"`
	UserDataBase64           string                   `hcl:"user_data_base64" hcle:"omitempty"`
	NetworkInterfaces        []AwsEc2NetworkInterface `hcl:"network_interface,blocks" hcle:"omitempty"`
	SecurityGroupIds         []string                 `hcl:"vpc_security_group_ids" hcle:"omitempty"`
	KeyName                  string                   `hcl:"key_name,expr" hcle:"omitempty"`
}

type AwsEc2NetworkInterface struct {
	NetworkInterfaceId string `hcl:"network_interface_id"`
	DeviceIndex        int    `hcl:"device_index"`
}

const AwsKeyPairResourceName = "aws_key_pair"

type AwsKeyPair struct {
	common.AwsResource `hcl:",squash"`
	KeyName            string `hcl:"key_name"`
	PublicKey          string `hcl:"public_key,expr"`
}
