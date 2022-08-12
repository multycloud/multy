package vault_secret

import (
	"github.com/multycloud/multy/resources/common"
)

const AwsResourceName = "aws_ssm_parameter"

type AwsSsmParameter struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_ssm_parameter"`
	Name                string `hcl:"name"`
	Type                string `hcl:"type"` // Valid types are String, StringList and SecureString.
	Value               string `hcl:"value" json:"value"`

	Arn string `json:"arn" hcle:"omitempty"`
}
