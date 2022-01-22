package lambda

import "multy-go/resources/common"

const AwsResourceName = "aws_lambda_function"

const AwsIamRoleResourceName = "aws_iam_role"

// TODO: support multiline strings in hclencoder
const DefaultLambdaPolicy = `{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}`

type AwsLambdaFunction struct {
	common.AwsResource `hcl:",squash"`
	FunctionName       string `hcl:"function_name"`
	Filename           string `hcl:"filename"`
	SourceCodeHash     string `hcl:"source_code_hash,expr"`
	Role               string `hcl:"role,expr"`
	Runtime            string `hcl:"runtime"`
	Handler            string `hcl:"handler"  hcle:"omitempty"`
}

type AwsIamRole struct {
	common.AwsResource `hcl:",squash"`
	AssumeRolePolicy   string `hcl:"assume_role_policy"`
}
