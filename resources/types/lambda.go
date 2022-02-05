package types

import (
	"fmt"
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
	ResourceName string `hcl:",key"`
	ResourceId   string `hcl:",key"`
	Type         string `hcl:"type"`
	SourceDir    string `hcl:"source_dir"`
	OutputPath   string `hcl:"output_path"`
}

const SasExpirationDuration = time.Hour * 24 * 365

func (r *Lambda) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []any {
	if cloud == common.AWS {
		var result []any

		function := lambda.AwsLambdaFunction{
			AwsResource:  common.NewAwsResource(lambda.AwsResourceName, r.GetTfResourceId(cloud), r.FunctionName),
			FunctionName: r.FunctionName,
			Runtime:      r.Runtime,
			Role:         fmt.Sprintf("aws_iam_role.%s.arn", r.getAwsIamRoleName()),
			Handler:      "lambda_function.lambda_handler",
		}

		if r.SourceCodeDir != "" {
			result = append(
				result, output.DataSourceWrapper{R: lambdaZip{
					ResourceName: "archive_file",
					ResourceId:   r.GetTfResourceId(cloud),

					Type:       "zip",
					SourceDir:  r.SourceCodeDir,
					OutputPath: r.getSourceCodeZip(cloud),
				}},
			)
			function.SourceCodeHash = fmt.Sprintf("data.archive_file.%s.output_base64sha256", r.GetTfResourceId(cloud))
			function.Filename = r.getSourceCodeZip(cloud)
		} else {
			function.S3Bucket = r.SourceCodeObject.ObjectStorage.GetResourceName(cloud)
			function.S3Key = r.SourceCodeObject.GetS3Key()
		}
		result = append(result, function)
		result = append(
			result, lambda.AwsIamRole{
				AwsResource: common.NewAwsResource(
					lambda.AwsIamRoleResourceName, r.getAwsIamRoleName(), r.getAwsIamRoleName(),
				),
				Name:             r.getAwsIamRoleName(),
				AssumeRolePolicy: lambda.DefaultLambdaPolicy,
			},
		)
		return result
	} else if cloud == common.AZURE {
		rgName := rg.GetResourceGroupName(r.ResourceGroupId, cloud)
		var result []any
		function := lambda.AzureFunctionApp{
			AzResource: common.NewAzResource(
				lambda.AzureResourceName, r.GetTfResourceId(cloud), common.AlphanumericFormatFunc(r.FunctionName),
				rgName,
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			),
			// AWS only supports linux
			OperatingSystem:  "linux",
			AppServicePlanId: fmt.Sprintf("%s.%s.id", lambda.AzureAppServicePlanResourceName, r.GetTfResourceId(cloud)),
		}
		if r.SourceCodeDir != "" {
			result = append(
				result, output.DataSourceWrapper{R: lambdaZip{
					ResourceName: "archive_file",
					ResourceId:   r.GetTfResourceId(cloud),
					Type:         "zip",
					SourceDir:    r.SourceCodeDir,
					OutputPath:   r.getSourceCodeZip(cloud),
				}},
			)
			result = append(
				result, object_storage.AzureStorageAccount{
					AzResource: common.NewAzResource(
						object_storage.AzureResourceName, r.GetTfResourceId(cloud),
						common.UniqueId(r.FunctionName, "stac", common.LowercaseAlphanumericFormatFunc),
						rgName,
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
					AzResource: common.AzResource{
						ResourceName: "azurerm_storage_account_blob_container_sas",
						ResourceId:   r.GetTfResourceId(cloud),
					},
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
				result = append(result, output.DataSourceWrapper{R: sas})
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
					lambda.AzureAppServicePlanResourceName, r.GetTfResourceId(cloud),
					common.UniqueId(r.FunctionName, "svpl", common.LowercaseAlphanumericFormatFunc),
					rgName, ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
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
