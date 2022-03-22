package iam

import (
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

type AwsIamRole struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_role"`
	Name                string                 `hcl:"name"`
	AssumeRolePolicy    string                 `hcl:"assume_role_policy"`
	InlinePolicy        AwsIamRoleInlinePolicy `hcl:"inline_policy,optional" hcle:"omitempty"`
}

type AwsIamRoleInlinePolicy struct {
	Name   string `hcl:"name"`
	Policy string `hcl:"policy"`
}

type AwsIamRolePolicyAttachment struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_role_policy_attachment"`
	Role                string `hcl:"role,expr"`
	PolicyArn           string `hcl:"policy_arn"`
}

type AwsIamPolicyStatementPrincipal struct {
	Service string `json:"Service,omitempty"`
}

type AwsIamPolicyStatement struct {
	Action    []string
	Effect    string
	Resource  string                          `json:"Resource,omitempty"`
	Principal *AwsIamPolicyStatementPrincipal `json:"Principal,omitempty"`
}

type AwsIamPolicy struct {
	Statement []AwsIamPolicyStatement
	Version   string
}

type AwsIamInstanceProfile struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_instance_profile"`
	Name                string `hcl:"name"`
	Role                string `hcl:"role,expr"`
}

func NewAssumeRolePolicy(services ...string) string {
	policy := AwsIamPolicy{
		Statement: util.MapSliceValues(services, func(service string) AwsIamPolicyStatement {
			return AwsIamPolicyStatement{
				Action:    []string{"sts:AssumeRole"},
				Effect:    "Allow",
				Principal: &AwsIamPolicyStatementPrincipal{Service: service},
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

func NewRoleResourcePolicy(resource ...string) string {
	policy := AwsIamPolicy{
		Statement: util.MapSliceValues(resource, func(resource string) AwsIamPolicyStatement {
			return AwsIamPolicyStatement{
				Action:   []string{"sts:AssumeRole"},
				Effect:   "Allow",
				Resource: resource,
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

func (r *AwsIamRole) GetId() string {
	return fmt.Sprintf("${%s.%s.id}", output.GetResourceName(AwsIamRole{}), r.ResourceId)
}

func (r *AwsIamInstanceProfile) GetId() string {
	return fmt.Sprintf("%s.%s.id", output.GetResourceName(AwsIamInstanceProfile{}), r.ResourceId)
}
