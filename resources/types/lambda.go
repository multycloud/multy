package types

import (
	"fmt"
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/lambda"
	"multy-go/resources/output/local_exec"
	"multy-go/resources/output/object_storage"
	"multy-go/resources/output/object_storage_object"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
	"time"
)

type Lambda struct {
	*resources.CommonResourceParams
	FunctionName string `hcl:"function_name"`
	// This is only used for AWS since azure has a .json file within the source code folder.
	Runtime          string               `hcl:"runtime"`
	SourceCodeDir    string               `hcl:"source_code_dir,optional"`
	SourceCodeObject *ObjectStorageObject `mhcl:"ref=source_code_object,optional"`
}

type lambdaZip struct {
	*output.TerraformDataSource `hcl:",squash"`
	Type                        string `hcl:"type"`
	SourceDir                   string `hcl:"source_dir"`
	OutputPath                  string `hcl:"output_path"`
}

const SasExpirationDuration = time.Hour * 24 * 365

func (r *Lambda) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		var result []output.TfBlock

		function := lambda.AwsLambdaFunction{
			AwsResource:  common.NewAwsResource(r.GetTfResourceId(cloud), r.FunctionName),
			FunctionName: r.FunctionName,
			Runtime:      r.Runtime,
			Role:         fmt.Sprintf("aws_iam_role.%s.arn", r.getAwsIamRoleName()),
			Handler:      "lambda_function.lambda_handler",
		}

		if r.SourceCodeDir != "" {
			result = append(
				result, lambdaZip{
					TerraformDataSource: &output.TerraformDataSource{
						ResourceName: "archive_file",
						ResourceId:   r.GetTfResourceId(cloud),
					},

					Type:       "zip",
					SourceDir:  r.SourceCodeDir,
					OutputPath: r.getSourceCodeZip(cloud),
				},
			)
			function.SourceCodeHash = fmt.Sprintf("data.archive_file.%s.output_base64sha256", r.GetTfResourceId(cloud))
			function.Filename = r.getSourceCodeZip(cloud)
		} else {
			function.S3Bucket = r.SourceCodeObject.ObjectStorage.GetResourceName(cloud)
			function.S3Key = r.SourceCodeObject.GetS3Key()
		}
		result = append(result, function)
		result = append(
			result,
			lambda.AwsIamRole{
				AwsResource:      common.NewAwsResource(r.getAwsIamRoleName(), r.getAwsIamRoleName()),
				Name:             r.getAwsIamRoleName(),
				AssumeRolePolicy: lambda.DefaultLambdaPolicy,
			},
			// this gives permission to write cloudwatch logs
			lambda.AwsIamRolePolicyAttachment{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)}},
				Role: fmt.Sprintf(
					"%s.%s.name", common.GetResourceName(lambda.AwsIamRole{}), r.getAwsIamRoleName(),
				),
				PolicyArn: lambda.LambdaBasicExecutionRole,
			},
			// https://registry.terraform.io/providers/hashicorp/aws/2.34.0/docs/guides/serverless-with-aws-lambda-and-api-gateway
			lambda.AwsApiGatewayRestApi{
				AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.FunctionName),
				Name:        r.FunctionName,
			},
			lambda.AwsApiGatewayResource{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: fmt.Sprintf(
					"%s_proxy", r.ResourceId,
				)}},
				RestApiId: r.getAwsRestApiId(),
				ParentId:  r.getAwsRestRootId(),
				PathPart:  "{proxy+}",
			},
			lambda.AwsApiGatewayMethod{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: fmt.Sprintf(
					"%s_proxy", r.ResourceId,
				)}},
				RestApiId: r.getAwsRestApiId(),
				ResourceId: fmt.Sprintf(
					"%s.%s.id", common.GetResourceName(lambda.AwsApiGatewayResource{}),
					fmt.Sprintf("%s_proxy", r.ResourceId),
				),
				HttpMethod:    "ANY",
				Authorization: "NONE",
			},
			lambda.AwsApiGatewayIntegration{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: fmt.Sprintf(
					"%s_proxy", r.ResourceId,
				)}},
				RestApiId: r.getAwsRestApiId(),
				ResourceId: fmt.Sprintf(
					"%s.%s.resource_id", common.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy", r.ResourceId),
				),
				HttpMethod: fmt.Sprintf(
					"%s.%s.http_method",
					common.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy", r.ResourceId),
				),
				IntegrationHttpMethod: "POST",
				Type:                  "AWS_PROXY",
				Uri: fmt.Sprintf(
					"%s.%s.invoke_arn", common.GetResourceName(lambda.AwsLambdaFunction{}),
					r.GetTfResourceId(cloud),
				),
			},
			lambda.AwsApiGatewayMethod{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: fmt.Sprintf(
					"%s_proxy_root", r.ResourceId,
				)}},
				RestApiId:     r.getAwsRestApiId(),
				ResourceId:    r.getAwsRestRootId(),
				HttpMethod:    "ANY",
				Authorization: "NONE",
			},
			lambda.AwsApiGatewayIntegration{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: fmt.Sprintf(
					"%s_proxy_root", r.ResourceId,
				)}},
				RestApiId: r.getAwsRestApiId(),
				ResourceId: fmt.Sprintf(
					"%s.%s.resource_id", common.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy_root", r.ResourceId),
				),
				HttpMethod: fmt.Sprintf(
					"%s.%s.http_method",
					common.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy_root", r.ResourceId),
				),
				IntegrationHttpMethod: "POST",
				Type:                  "AWS_PROXY",
				Uri: fmt.Sprintf(
					"%s.%s.invoke_arn", common.GetResourceName(lambda.AwsLambdaFunction{}),
					r.GetTfResourceId(cloud),
				),
			},
			lambda.AwsApiGatewayDeployment{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)}},
				RestApiId:   r.getAwsRestApiId(),
				StageName:   "api",
				DependsOn: []string{
					fmt.Sprintf(
						"%s.%s", common.GetResourceName(lambda.AwsApiGatewayIntegration{}),
						r.ResourceId+"_proxy",
					),
					fmt.Sprintf(
						"%s.%s", common.GetResourceName(lambda.AwsApiGatewayIntegration{}),
						r.ResourceId+"_proxy_root",
					),
				},
			},
			lambda.AwsLambdaPermission{
				AwsResource:  &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.GetTfResourceId(cloud)}},
				StatementId:  "AllowAPIGatewayInvoke",
				Action:       "lambda:InvokeFunction",
				FunctionName: r.FunctionName,
				Principal:    "apigateway.amazonaws.com",
				SourceArn: fmt.Sprintf(
					"${%s}/*/*", fmt.Sprintf(
						"%s.%s.execution_arn", common.GetResourceName(lambda.AwsApiGatewayRestApi{}),
						r.GetTfResourceId(common.AWS),
					),
				),
			},
		)
		return result
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		var result []output.TfBlock
		function := lambda.AzureFunctionApp{
			AzResource: common.NewAzResource(
				r.GetTfResourceId(cloud), common.AlphanumericFormatFunc(r.FunctionName), rgName,
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			),
			// AWS only supports linux
			OperatingSystem:  "linux",
			AppServicePlanId: fmt.Sprintf("%s.%s.id", lambda.AzureAppServicePlanResourceName, r.GetTfResourceId(cloud)),
		}
		if r.SourceCodeDir != "" {
			result = append(
				result, lambdaZip{
					TerraformDataSource: &output.TerraformDataSource{
						ResourceName: "archive_file",
						ResourceId:   r.GetTfResourceId(cloud),
					},
					Type:       "zip",
					SourceDir:  r.SourceCodeDir,
					OutputPath: r.getSourceCodeZip(cloud),
				},
			)
			result = append(
				result, object_storage.AzureStorageAccount{
					AzResource: common.NewAzResource(
						r.GetTfResourceId(cloud),
						common.UniqueId(r.FunctionName, "stac", common.LowercaseAlphanumericFormatFunc), rgName,
						ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
					),
					AccountTier:            "Standard",
					AccountReplicationType: "LRS",
				},
			)
			function.StorageAccountName = fmt.Sprintf(
				"%s.%s.name", object_storage.AzureResourceName, r.GetTfResourceId(cloud),
			)
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName, r.GetTfResourceId(cloud),
			)
			function.LocalExec = local_exec.New(
				local_exec.LocalExec{
					Command: fmt.Sprintf(
						"az functionapp deployment source config-zip -g ${self.resource_group_name} -n ${self."+
							"name} --src ${data.archive_file.%s.output_path}",
						r.GetTfResourceId(cloud),
					),
				},
			)
		} else {
			function.StorageAccountName = r.SourceCodeObject.ObjectStorage.GetResourceName(cloud)
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName,
				r.SourceCodeObject.ObjectStorage.GetTfResourceId(common.AZURE),
			)
			if r.SourceCodeObject.IsPrivate() {
				sas := object_storage_object.AzureStorageAccountBlobSas{
					TerraformDataSource: &output.TerraformDataSource{ResourceId: r.GetTfResourceId(cloud)},
					ConnectionString: fmt.Sprintf(
						"azurerm_storage_account.%s.primary_connection_string",
						r.SourceCodeObject.ObjectStorage.GetTfResourceId(cloud),
					),
					ContainerName: r.SourceCodeObject.ObjectStorage.GetAssociatedPrivateContainerResourceName(cloud),
					Start:         time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05Z"),
					Expiry:        time.Now().Add(SasExpirationDuration).Format("2006-01-02T15:04:05Z"),
					AzureStorageAccountBlobSasPermissions: object_storage_object.AzureStorageAccountBlobSasPermissions{
						Read: true,
					},
				}
				result = append(result, sas)
				function.AppSettings = map[string]string{
					"WEBSITE_RUN_FROM_PACKAGE": sas.GetSignedUrl(
						r.SourceCodeObject.ObjectStorage.GetResourceName(cloud),
						r.SourceCodeObject.ObjectStorage.GetAssociatedPrivateContainerResourceName(cloud),
						r.SourceCodeObject.GetAzureBlobName(),
					),
				}
			} else {
				function.AppSettings = map[string]string{
					"WEBSITE_RUN_FROM_PACKAGE": fmt.Sprintf("${%s}", r.SourceCodeObject.GetAzureBlobUrl()),
				}
			}
		}

		result = append(result, function)

		result = append(
			result, lambda.AzureAppServicePlan{
				AzResource: common.NewAzResource(
					r.GetTfResourceId(cloud),
					common.UniqueId(r.FunctionName, "svpl", common.LowercaseAlphanumericFormatFunc), rgName,
					ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
				),
				Kind:     "Linux",
				Reserved: true,
				Sku: lambda.AzureSku{
					Tier: "Dynamic",
					Size: "Y1",
				},
			},
		)
		return result
	}
	return nil
}

func (r *Lambda) Validate(ctx resources.MultyContext) {
	if r.SourceCodeDir == "" && r.SourceCodeObject == nil {
		r.LogFatal(r.ResourceId, "", "one of source_code_dir or source_code_object must be set")
	}
	if r.SourceCodeDir != "" && r.SourceCodeObject != nil {
		r.LogFatal(r.ResourceId, "source_code_dir", "only one of source_code_dir or source_code_object can be set")
	}
	return
}

func (r *Lambda) getAwsIamRoleName() string {
	return fmt.Sprintf("iam_for_lambda_%s", r.ResourceId)
}

func (r *Lambda) getAwsRestApiId() string {
	return fmt.Sprintf("%s.%s.id", common.GetResourceName(lambda.AwsApiGatewayRestApi{}), r.GetTfResourceId(common.AWS))
}

func (r *Lambda) getAwsRestRootId() string {
	return fmt.Sprintf(
		"%s.%s.root_resource_id", common.GetResourceName(lambda.AwsApiGatewayRestApi{}),
		r.GetTfResourceId(common.AWS),
	)
}

func (r *Lambda) getSourceCodeZip(cloud common.CloudProvider) string {
	if r.SourceCodeDir == "" {
		return ""
	}
	return fmt.Sprintf(".multy/tmp/%s_%s.zip", r.FunctionName, cloud)
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

func (r *Lambda) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"url": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.invoke_url}", common.GetResourceName(lambda.AwsApiGatewayDeployment{}),
					r.GetTfResourceId(cloud),
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"url": cty.StringVal(
				fmt.Sprintf(
					"https://${%s.%s.default_hostname}",
					common.GetResourceName(lambda.AzureFunctionApp{}), r.GetTfResourceId(cloud),
				),
			),
		}
	}

	validate.LogInternalError("unknown cloud %s", cloud)
	return nil
}
