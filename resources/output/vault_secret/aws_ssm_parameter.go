package vault_secret

import "multy-go/resources/common"

const AwsResourceName = "aws_ssm_parameter"

type AwsSsmParameter struct {
	common.AwsResource `hcl:",squash"`
	Type               string `hcl:"type"` // Valid types are String, StringList and SecureString.
	Value              string `hcl:"value"`
}
