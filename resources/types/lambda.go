package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/lambda"
	"multy-go/resources/output/local_exec"
	"multy-go/resources/output/object_storage"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
	"strings"
)

type Lambda struct {
	*resources.CommonResourceParams
	FunctionName string `hcl:"function_name"`
	// This is only used for AWS since azure has a .json file within the source code folder.
	Runtime       string `hcl:"runtime"`
	SourceCodeDir string `hcl:"source_code_dir"`
}

type awsLambdaZip struct {
	common.AwsResource `hcl:",squash"`
	Type               string `hcl:"type"`
	SourceDir          string `hcl:"source_dir"`
	OutputPath         string `hcl:"output_path"`
}

func (r *Lambda) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []any {
	if cloud == common.AWS {
		return []any{
			lambda.AwsLambdaFunction{
				AwsResource:    common.NewAwsResource(lambda.AwsResourceName, r.GetTfResourceId(cloud), r.FunctionName),
				FunctionName:   r.FunctionName,
				Runtime:        r.Runtime,
				Filename:       r.getAwsZipFile(),
				SourceCodeHash: fmt.Sprintf("data.archive_file.%s.output_base64sha256", r.GetTfResourceId(cloud)),
				Role:           fmt.Sprintf("aws_iam_role.%s.arn", r.getAwsIamRoleName()),
				Handler:        "lambda_function.lambda_handler",
			},
			lambda.AwsIamRole{
				AwsResource: common.NewAwsResource(
					lambda.AwsIamRoleResourceName, r.getAwsIamRoleName(), r.getAwsIamRoleName(),
				),
				Name:             r.getAwsIamRoleName(),
				AssumeRolePolicy: lambda.DefaultLambdaPolicy,
			},
			output.DataSourceWrapper{R: awsLambdaZip{
				AwsResource: common.AwsResource{
					ResourceName: "archive_file",
					ResourceId:   r.GetTfResourceId(cloud),
				},
				Type:       "zip",
				SourceDir:  r.SourceCodeDir,
				OutputPath: r.getAwsZipFile(),
			}},
		}
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		return []any{
			object_storage.AzureStorageAccount{
				AzResource: common.NewAzResource(
					object_storage.AzureResourceName, r.ResourceId, fmt.Sprintf("%sstacct", r.ResourceId), rgName,
					ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				AccountTier:            "Standard",
				AccountReplicationType: "LRS",
			},
			lambda.AzureAppServicePlan{
				AzResource: common.NewAzResource(
					lambda.AzureAppServicePlanResourceName, r.ResourceId, fmt.Sprintf("%sservplan", r.ResourceId),
					rgName, ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				Kind:     "Linux",
				Reserved: true,
				Sku: lambda.AzureSku{
					Tier: "Dynamic",
					Size: "Y1",
				},
			},
			lambda.AzureFunctionApp{
				AzResource: common.NewAzResource(
					lambda.AzureResourceName, r.ResourceId, strings.ReplaceAll(r.FunctionName, "_", ""), rgName,
					ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				// AWS only supports linux
				OperatingSystem:    "linux",
				AppServicePlanId:   fmt.Sprintf("%s.%s.id", lambda.AzureAppServicePlanResourceName, r.ResourceId),
				StorageAccountName: fmt.Sprintf("%s.%s.name", object_storage.AzureResourceName, r.ResourceId),
				StorageAccountAccessKey: fmt.Sprintf(
					"%s.%s.primary_access_key", object_storage.AzureResourceName, r.ResourceId,
				),
				LocalExec: local_exec.New(
					local_exec.LocalExec{
						WorkingDir: r.SourceCodeDir,
						Command:    "func azure functionapp publish ${self.id}",
					},
				),
			},
		}
	}
	return nil
}

func (r *Lambda) Validate(ctx resources.MultyContext) {
	return
}

func (r *Lambda) getAwsIamRoleName() string {
	return fmt.Sprintf("iam_for_lambda_%s", r.FunctionName)
}

func (r *Lambda) getAwsZipFile() string {
	return fmt.Sprintf(".multy/tmp/%s.zip", r.FunctionName)
}

func (r *Lambda) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return lambda.AwsResourceName
	case common.AZURE:
		return lambda.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}
