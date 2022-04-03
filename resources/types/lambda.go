package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/lambda"
	"github.com/multycloud/multy/resources/output/local_exec"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/output/object_storage_object"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
	"time"
)

type Lambda struct {
	resources.ResourceWithId[*resourcespb.LambdaArgs]

	SourceCodeObject *ObjectStorageObject `mhcl:"ref=source_code_object,optional"`
}

type lambdaZip struct {
	*output.TerraformDataSource `hcl:",squash"`
	Type                        string `hcl:"type"`
	SourceDir                   string `hcl:"source_dir"`
	OutputPath                  string `hcl:"output_path"`
}

const SasExpirationDuration = time.Hour * 24 * 365

func NewLambda(resourceId string, args *resourcespb.LambdaArgs, others resources.Resources) (*Lambda, error) {
	obj, _, err := GetOptional[*ObjectStorageObject](others, args.SourceCodeObjectId)
	if err != nil {
		return nil, err
	}

	return &Lambda{
		ResourceWithId: resources.ResourceWithId[*resourcespb.LambdaArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		SourceCodeObject: obj,
	}, nil
}

func (r *Lambda) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	sourceCodeDir := ""
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var result []output.TfBlock

		function := lambda.AwsLambdaFunction{
			AwsResource:  common.NewAwsResource(r.ResourceId, r.Args.Name),
			FunctionName: r.Args.Name,
			Runtime:      r.Args.Runtime,
			Role:         fmt.Sprintf("aws_iam_role.%s.arn", r.getAwsIamRoleName()),
			Handler:      "lambda_function.lambda_handler",
		}

		if false {
			result = append(
				result, lambdaZip{
					TerraformDataSource: &output.TerraformDataSource{
						ResourceName: "archive_file",
						ResourceId:   r.ResourceId,
					},

					Type:       "zip",
					SourceDir:  sourceCodeDir,
					OutputPath: r.getSourceCodeZip(),
				},
			)
			function.SourceCodeHash = fmt.Sprintf("data.archive_file.%s.output_base64sha256", r.ResourceId)
			function.Filename = r.getSourceCodeZip()
		} else {
			function.S3Bucket = r.SourceCodeObject.ObjectStorage.GetResourceName()
			function.S3Key = r.SourceCodeObject.GetS3Key()
		}
		result = append(result, function)
		result = append(
			result,
			iam.AwsIamRole{
				AwsResource:      common.NewAwsResource(r.getAwsIamRoleName(), r.getAwsIamRoleName()),
				Name:             r.getAwsIamRoleName(),
				AssumeRolePolicy: lambda.DefaultLambdaPolicy,
			},
			// this gives permission to write cloudwatch logs
			iam.AwsIamRolePolicyAttachment{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
				Role: fmt.Sprintf(
					"%s.%s.name", output.GetResourceName(iam.AwsIamRole{}), r.getAwsIamRoleName(),
				),
				PolicyArn: lambda.LambdaBasicExecutionRole,
			},
			// https://registry.terraform.io/providers/hashicorp/aws/2.34.0/docs/guides/serverless-with-aws-lambda-and-api-gateway
			lambda.AwsApiGatewayRestApi{
				AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
				Name:        r.Args.Name,
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
					"%s.%s.id", output.GetResourceName(lambda.AwsApiGatewayResource{}),
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
					"%s.%s.resource_id", output.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy", r.ResourceId),
				),
				HttpMethod: fmt.Sprintf(
					"%s.%s.http_method",
					output.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy", r.ResourceId),
				),
				IntegrationHttpMethod: "POST",
				Type:                  "AWS_PROXY",
				Uri: fmt.Sprintf(
					"%s.%s.invoke_arn", output.GetResourceName(lambda.AwsLambdaFunction{}),
					r.ResourceId,
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
					"%s.%s.resource_id", output.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy_root", r.ResourceId),
				),
				HttpMethod: fmt.Sprintf(
					"%s.%s.http_method",
					output.GetResourceName(lambda.AwsApiGatewayMethod{}),
					fmt.Sprintf("%s_proxy_root", r.ResourceId),
				),
				IntegrationHttpMethod: "POST",
				Type:                  "AWS_PROXY",
				Uri: fmt.Sprintf(
					"%s.%s.invoke_arn", output.GetResourceName(lambda.AwsLambdaFunction{}),
					r.ResourceId,
				),
			},
			lambda.AwsApiGatewayDeployment{
				AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
				RestApiId:   r.getAwsRestApiId(),
				StageName:   "api",
				DependsOn: []string{
					fmt.Sprintf(
						"%s.%s", output.GetResourceName(lambda.AwsApiGatewayIntegration{}),
						r.ResourceId+"_proxy",
					),
					fmt.Sprintf(
						"%s.%s", output.GetResourceName(lambda.AwsApiGatewayIntegration{}),
						r.ResourceId+"_proxy_root",
					),
				},
			},
			lambda.AwsLambdaPermission{
				AwsResource:  &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
				StatementId:  "AllowAPIGatewayInvoke",
				Action:       "lambda:InvokeFunction",
				FunctionName: r.Args.Name,
				Principal:    "apigateway.amazonaws.com",
				SourceArn: fmt.Sprintf(
					"${%s}/*/*", fmt.Sprintf(
						"%s.%s.execution_arn", output.GetResourceName(lambda.AwsApiGatewayRestApi{}),
						r.ResourceId,
					),
				),
			},
		)
		return result, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		rgName := rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
		var result []output.TfBlock
		function := lambda.AzureFunctionApp{
			AzResource: common.NewAzResource(
				r.ResourceId, common.AlphanumericFormatFunc(r.Args.Name), rgName,
				r.GetCloudSpecificLocation(),
			),
			// AWS only supports linux
			OperatingSystem:  "linux",
			AppServicePlanId: fmt.Sprintf("%s.%s.id", lambda.AzureAppServicePlanResourceName, r.ResourceId),
		}
		if false {
			result = append(
				result, lambdaZip{
					TerraformDataSource: &output.TerraformDataSource{
						ResourceName: "archive_file",
						ResourceId:   r.ResourceId,
					},
					Type:       "zip",
					SourceDir:  sourceCodeDir,
					OutputPath: r.getSourceCodeZip(),
				},
			)
			result = append(
				result, object_storage.AzureStorageAccount{
					AzResource: common.NewAzResource(
						r.ResourceId,
						common.UniqueId(r.Args.Name, "stac", common.LowercaseAlphanumericFormatFunc), rgName,
						r.GetCloudSpecificLocation(),
					),
					AccountTier:            "Standard",
					AccountReplicationType: "LRS",
				},
			)
			function.StorageAccountName = fmt.Sprintf(
				"%s.%s.name", object_storage.AzureResourceName, r.ResourceId,
			)
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName, r.ResourceId,
			)
			function.LocalExec = local_exec.New(
				local_exec.LocalExec{
					Command: fmt.Sprintf(
						"az functionapp deployment source config-zip -g ${self.resource_group_name} -n ${self."+
							"name} --src ${data.archive_file.%s.output_path}",
						r.ResourceId,
					),
				},
			)
		} else {
			function.StorageAccountName = r.SourceCodeObject.ObjectStorage.GetResourceName()
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName,
				r.SourceCodeObject.ObjectStorage.ResourceId,
			)
			if r.SourceCodeObject.IsPrivate() {
				sas := object_storage_object.AzureStorageAccountBlobSas{
					TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId},
					ConnectionString: fmt.Sprintf(
						"azurerm_storage_account.%s.primary_connection_string",
						r.SourceCodeObject.ObjectStorage.ResourceId,
					),
					ContainerName: r.SourceCodeObject.ObjectStorage.GetAssociatedPrivateContainerResourceName(),
					Start:         time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05Z"),
					Expiry:        time.Now().Add(SasExpirationDuration).Format("2006-01-02T15:04:05Z"),
					AzureStorageAccountBlobSasPermissions: object_storage_object.AzureStorageAccountBlobSasPermissions{
						Read: true,
					},
				}
				result = append(result, sas)
				function.AppSettings = map[string]string{
					"WEBSITE_RUN_FROM_PACKAGE": sas.GetSignedUrl(
						r.SourceCodeObject.ObjectStorage.GetResourceName(),
						r.SourceCodeObject.ObjectStorage.GetAssociatedPrivateContainerResourceName(),
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
					r.ResourceId,
					common.UniqueId(r.Args.Name, "svpl", common.LowercaseAlphanumericFormatFunc), rgName,
					r.GetCloudSpecificLocation(),
				),
				Kind:     "Linux",
				Reserved: true,
				Sku: lambda.AzureSku{
					Tier: "Dynamic",
					Size: "Y1",
				},
			},
		)
		return result, nil
	}
	return nil, nil
}

func (r *Lambda) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if r.SourceCodeObject == nil {
		errs = append(errs, r.NewValidationError("one of source_code_dir or source_code_object must be set", "source_code_object_id"))
	}
	return errs
}

func (r *Lambda) getAwsIamRoleName() string {
	return fmt.Sprintf("iam_for_lambda_%s", r.ResourceId)
}

func (r *Lambda) getAwsRestApiId() string {
	return fmt.Sprintf("%s.%s.id", output.GetResourceName(lambda.AwsApiGatewayRestApi{}), r.ResourceId)
}

func (r *Lambda) getAwsRestRootId() string {
	return fmt.Sprintf(
		"%s.%s.root_resource_id", output.GetResourceName(lambda.AwsApiGatewayRestApi{}),
		r.ResourceId,
	)
}

func (r *Lambda) getSourceCodeZip() string {
	return fmt.Sprintf(".multy/tmp/%s.zip", r.Args.Name)
}

func (r *Lambda) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return lambda.AwsResourceName, nil
	case common.AZURE:
		return lambda.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *Lambda) GetOutputValues(cloud commonpb.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"url": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.invoke_url}", output.GetResourceName(lambda.AwsApiGatewayDeployment{}),
					r.ResourceId,
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"url": cty.StringVal(
				fmt.Sprintf(
					"https://${%s.%s.default_hostname}",
					output.GetResourceName(lambda.AzureFunctionApp{}), r.ResourceId,
				),
			),
		}
	}
	return nil
}
