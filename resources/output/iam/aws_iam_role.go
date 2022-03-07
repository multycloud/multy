package iam

import (
	"encoding/json"
	"multy-go/resources/common"
	"multy-go/util"
	"multy-go/validate"
)

type AwsIamRole struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_role"`
	Name                string `hcl:"name"`
	AssumeRolePolicy    string `hcl:"assume_role_policy"`
}

type AwsIamRolePolicyAttachment struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_role_policy_attachment"`
	Role                string `hcl:"role,expr"`
	PolicyArn           string `hcl:"policy_arn"`
}

type AwsIamPolicyStatementPrincipal struct {
	Service string
}

type AwsIamPolicyStatement struct {
	Action    string
	Effect    string
	Principal AwsIamPolicyStatementPrincipal
}

type AwsIamPolicy struct {
	Statement []AwsIamPolicyStatement
	Version   string
}

func NewAssumeRolePolicy(services ...string) string {
	policy := AwsIamPolicy{
		Statement: util.MapSliceValues(services, func(service string) AwsIamPolicyStatement {
			return AwsIamPolicyStatement{
				Action:    "sts:AssumeRole",
				Effect:    "Allow",
				Principal: AwsIamPolicyStatementPrincipal{Service: service},
			}
		}),
		Version: "2012-10-17",
	}

	b, err := json.Marshal(policy)
	if err != nil {
		validate.LogInternalError("unable to encode aws policy: %s", err.Error())
	}
	return string(b)
}
