package network_security_group

import (
	"github.com/multycloud/multy/resources/common"
)

type AwsAcl struct {
	*common.AwsResource `hcl:",squash"`
	VpcId               string       `hcl:"vpc_id,expr"`
	SubnetIds           []string     `hcl:"subnet_ids,expr"`
	Ingress             []AwsAclRule `hcl:"ingress"`
	Egress              []AwsAclRule `hcl:"egress"`
}

const AwsDefaultSecurityGroupResourceName = "aws_default_security_group"

type AwsDefaultSecurityGroup struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_default_security_group"`
	VpcId               string                 `hcl:"vpc_id,expr"`
	Ingress             []AwsSecurityGroupRule `hcl:"ingress,blocks"`
	Egress              []AwsSecurityGroupRule `hcl:"egress,blocks"`
}

const AwsSecurityGroupResourceName = "aws_security_group"

type AwsSecurityGroup struct {
	*common.AwsResource `hcl:",squash"  default:"name=aws_security_group"`
	VpcId               string                 `hcl:"vpc_id,expr"`
	Name                string                 `hcl:"name"`
	Description         string                 `hcl:"description"`
	Ingress             []AwsSecurityGroupRule `hcl:"ingress,blocks"`
	Egress              []AwsSecurityGroupRule `hcl:"egress,blocks"`
}

type AwsSecurityGroupRule struct {
	Protocol   string   `hcl:"protocol" json:"protocol"`
	FromPort   int      `hcl:"from_port" json:"from_port"`
	ToPort     int      `hcl:"to_port" json:"to_port"`
	CidrBlocks []string `hcl:"cidr_blocks" json:"cidr_blocks"`
	Self       bool     `hcl:"self" hcle:"omitempty" json:"self"`
}

// checks if 2 rules are equal, assuming they have at most 1 cidr block
func (r AwsSecurityGroupRule) Equals(other AwsSecurityGroupRule) bool {
	return r.Protocol == other.Protocol && r.FromPort == other.FromPort && r.ToPort == other.ToPort &&
		r.Self == other.Self &&
		len(r.CidrBlocks) == len(other.CidrBlocks) && (len(r.CidrBlocks) == 0 || r.CidrBlocks[0] == other.CidrBlocks[0])
}

type AwsDefaultAcl struct {
	*common.AwsResource `hcl:",squash"`
	DefaultNetworkAclId string       `hcl:"default_network_acl_id,expr"`
	Ingress             []AwsAclRule `hcl:"ingress"`
	Egress              []AwsAclRule `hcl:"egress"`
}

type AwsAclRule struct {
	Protocol   string   `hcl:"protocol"` // ALL: Aws = -1 (FromPort & ToPort must be 0) / SubnetAz = "*"
	RuleNumber int      `hcl:"rule_no"`
	Action     string   `hcl:"action"`
	FromPort   int      `hcl:"from_port"`
	ToPort     int      `hcl:"to_port"`
	CidrBlock  []string `hcl:"cidr_block"`
}
