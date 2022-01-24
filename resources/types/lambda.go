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
	Runtime          string               `hcl:"runtime"`
	SourceCodeDir    string               `hcl:"source_code_dir,optional"`
	SourceCodeObject *ObjectStorageObject `mhcl:"ref=source_code_object,optional"`
}

type awsLambdaZip struct {
	common.AwsResource `hcl:",squash"`
	Type               string `hcl:"type"`
	SourceDir          string `hcl:"source_dir"`
	OutputPath         string `hcl:"output_path"`
}

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
				result, output.DataSourceWrapper{R: awsLambdaZip{
					AwsResource: common.AwsResource{
						ResourceName: "archive_file",
						ResourceId:   r.GetTfResourceId(cloud),
					},
					Type:       "zip",
					SourceDir:  r.SourceCodeDir,
					OutputPath: r.getAwsZipFile(),
				}},
			)
			function.SourceCodeHash = fmt.Sprintf("data.archive_file.%s.output_base64sha256", r.GetTfResourceId(cloud))
			function.Filename = r.getAwsZipFile()
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
				lambda.AzureResourceName, r.ResourceId, strings.ReplaceAll(r.FunctionName, "_", ""), rgName,
				ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
			),
			// AWS only supports linux
			OperatingSystem:  "linux",
			AppServicePlanId: fmt.Sprintf("%s.%s.id", lambda.AzureAppServicePlanResourceName, r.ResourceId),
		}
		if r.SourceCodeDir != "" {
			result = append(
				result, object_storage.AzureStorageAccount{
					AzResource: common.NewAzResource(
						object_storage.AzureResourceName, r.ResourceId, fmt.Sprintf("%sstacct", r.ResourceId), rgName,
						ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud),
					),
					AccountTier:            "Standard",
					AccountReplicationType: "LRS",
				},
			)
			function.StorageAccountName = fmt.Sprintf("%s.%s.name", object_storage.AzureResourceName, r.ResourceId)
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName, r.ResourceId,
			)
			function.LocalExec = local_exec.New(
				local_exec.LocalExec{
					WorkingDir: r.SourceCodeDir,
					Command:    "func azure functionapp publish ${self.id}",
				},
			)
		} else {
			function.StorageAccountName = r.SourceCodeObject.ObjectStorage.GetResourceName(cloud)
			function.StorageAccountAccessKey = fmt.Sprintf(
				"%s.%s.primary_access_key", object_storage.AzureResourceName,
				r.SourceCodeObject.ObjectStorage.GetTfResourceId(common.AZURE),
			)
			// TODO: create sas and test this
			//function.AppSettings = map[string]string{
			//	"WEBSITE_RUN_FROM_PACKAGE": fmt.Sprintf(
			//		"https://${%s}.blob.core.windows.net/${%s}/${%s}${data."+
			//			"azurerm_storage_account_blob_container_sas.storage_account_blob_container_sas.sas}",
			//		r.SourceCodeObject.ObjectStorage.GetResourceName(cloud),
			//		r.SourceCodeObject.ObjectStorage.GetAssociatedPrivateContainerResourceName(cloud),
			//		r.SourceCodeObject.GetAzureBlobName(),
			//	),
			//}
		}

		result = append(result, function)

		result = append(
			result, lambda.AzureAppServicePlan{
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
	return fmt.Sprintf("iam_for_lambda_%s", r.FunctionName)
}

func (r *Lambda) getAwsZipFile() string {
	if r.SourceCodeDir == "" {
		return ""
	}
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
