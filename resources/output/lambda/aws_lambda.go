package lambda

import "multy-go/resources/common"

const AwsResourceName = "aws_lambda_function"

const AwsIamRoleResourceName = "aws_iam_role"

// TODO: support multiline strings in hclencoder
const DefaultLambdaPolicy = `{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": \"sts:AssumeRole\",\"Principal\": {\"Service\": \"lambda.amazonaws.com\"},\"Effect\": \"Allow\",\"Sid\": \"\"}]}`

type AwsLambdaFunction struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_lambda_function"`
	FunctionName        string `hcl:"function_name"`
	Role                string `hcl:"role,expr"`
	Filename            string `hcl:"filename" hcle:"omitempty"`
	SourceCodeHash      string `hcl:"source_code_hash,expr" hcle:"omitempty"`
	Runtime             string `hcl:"runtime" hcle:"omitempty"`
	Handler             string `hcl:"handler"  hcle:"omitempty"`
	S3Bucket            string `hcl:"s3_bucket,expr" hcle:"omitempty"`
	S3Key               string `hcl:"s3_key,expr" hcle:"omitempty"`
}

type AwsIamRole struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_iam_role"`
	Name                string `hcl:"name"`
	AssumeRolePolicy    string `hcl:"assume_role_policy"`
}
