package network_security_group

import "multy-go/resources/common"

type AwsAcl struct {
	common.AwsResource `hcl:",squash"`
	VpcId              string       `hcl:"vpc_id,expr"`
	SubnetIds          []string     `hcl:"subnet_ids"`
	Ingress            []AwsAclRule `hcl:"ingress"`
	Egress             []AwsAclRule `hcl:"egress"`
}

const AwsDefaultSecurityGroupResourceName = "aws_default_security_group"

type AwsDefaultSecurityGroup struct {
	common.AwsResource `hcl:",squash"`
	VpcId              string                 `hcl:"vpc_id,expr"`
	Ingress            []AwsSecurityGroupRule `hcl:"ingress,blocks"`
	Egress             []AwsSecurityGroupRule `hcl:"egress,blocks"`
}

const AwsSecurityGroupResourceName = "aws_security_group"

type AwsSecurityGroup struct {
	common.AwsResource `hcl:",squash"`
	VpcId              string                 `hcl:"vpc_id,expr"`
	Ingress            []AwsSecurityGroupRule `hcl:"ingress,blocks"`
	Egress             []AwsSecurityGroupRule `hcl:"egress,blocks"`
}

type AwsSecurityGroupRule struct {
	Protocol   string `hcl:"protocol"`
	FromPort   int    `hcl:"from_port"`
	ToPort     int    `hcl:"to_port"`
	CidrBlocks string `hcl:"cidr_blocks"`
	Self       bool   `hcl:"self" hcle:"omitempty"`
}

type AwsDefaultAcl struct {
	common.AwsResource  `hcl:",squash"`
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
