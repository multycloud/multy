package lambda

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output/iam"
)

const AwsResourceName = "aws_lambda_function"

const AwsIamRoleResourceName = "aws_iam_role"

var DefaultLambdaPolicy = iam.NewAssumeRolePolicy("lambda.amazonaws.com")

const LambdaBasicExecutionRole = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

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

type AwsApiGatewayRestApi struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_api_gateway_rest_api"`
	Name                string `hcl:"name"`
	Description         string `hcl:"description"`
}

type AwsApiGatewayResource struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_api_gateway_resource"`
	RestApiId           string `hcl:"rest_api_id,expr"`
	ParentId            string `hcl:"parent_id,expr"`
	PathPart            string `hcl:"path_part"`
}

type AwsApiGatewayMethod struct {
	*common.AwsResource `hcl:",squash"  default:"name=aws_api_gateway_method"`
	RestApiId           string `hcl:"rest_api_id,expr"`
	ResourceId          string `hcl:"resource_id,expr"`
	HttpMethod          string `hcl:"http_method"`
	Authorization       string `hcl:"authorization"`
}

type AwsApiGatewayIntegration struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_api_gateway_integration"`
	RestApiId           string `hcl:"rest_api_id,expr"`
	ResourceId          string `hcl:"resource_id,expr"`
	HttpMethod          string `hcl:"http_method,expr"`

	IntegrationHttpMethod string `hcl:"integration_http_method"`
	Type                  string `hcl:"type"`
	Uri                   string `hcl:"uri,expr"`
}

type AwsApiGatewayDeployment struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_api_gateway_deployment"`
	RestApiId           string `hcl:"rest_api_id,expr"`
	StageName           string `hcl:"stage_name"`

	DependsOn []string `hcl:"depends_on,expr"`
}

type AwsLambdaPermission struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_lambda_permission"`
	StatementId         string `hcl:"statement_id,expr"`
	Action              string `hcl:"action"`
	FunctionName        string `hcl:"function_name"`
	Principal           string `hcl:"principal"`
	SourceArn           string `hcl:"source_arn"`
}
