package virtual_machine

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
)

const AwsResourceName = "aws_instance"

type AwsEC2 struct {
	*common.AwsResource      `hcl:",squash"  default:"name=aws_instance"`
	Ami                      string                   `hcl:"ami,expr"`
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

type AwsAmi struct {
	*output.TerraformDataSource `hcl:",squash"  default:"name=aws_ami"`

	Owners     []string `hcl:"owners"`
	MostRecent bool     `hcl:"most_recent"`
	NameRegex  string   `hcl:"name_regex"  hcle:"omitempty"`

	Filters []AwsAmiFilter `hcl:"filter,blocks"`
}

type AwsAmiFilter struct {
	Name   string   `hcl:"name"`
	Values []string `hcl:"values"`
}

// Ubuntu: ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20210928
// Debian: debian-11-arm64-20220503-998
// CentOS: CentOS Stream 9 x86_64 20211006

func LatestAwsAmi(image *resourcespb.ImageReference, resourceId string) (AwsAmi, error) {
	name, err := GetNamePattern(image)
	if err != nil {
		return AwsAmi{}, err
	}

	return AwsAmi{
		TerraformDataSource: &output.TerraformDataSource{ResourceId: resourceId, ResourceName: "aws_ami"},
		Owners:              []string{common.AwsAmiOwners[image.Os]},
		MostRecent:          true,
		Filters: []AwsAmiFilter{
			{
				Name:   "name",
				Values: []string{name},
			},
			{
				Name:   "root-device-type",
				Values: []string{"ebs"},
			},
			{
				Name:   "virtualization-type",
				Values: []string{"hvm"},
			},
		},
	}, nil
}

func GetNamePattern(image *resourcespb.ImageReference) (string, error) {
	switch image.Os {
	case resourcespb.ImageReference_UBUNTU:
		return fmt.Sprintf("ubuntu*-%s-amd64-server-*", image.Version), nil
	case resourcespb.ImageReference_DEBIAN:
		return fmt.Sprintf("debian-%s-amd64-*", image.Version), nil
	case resourcespb.ImageReference_CENT_OS:
		return fmt.Sprintf("CentOS %s* x86_64", image.Version), nil
	default:
		return "", fmt.Errorf("unknown operating system distibution %s", image.Os)
	}
}
