package functions

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/zclconf/go-cty/cty/function"
)

func GetAllFunctions(cloud common.CloudProvider) map[string]function.Function {
	functions := GetAllCloudAgnosticFunctions()
	functions["cloud_specific_value"] = GetCloudSpecificValueFunction(cloud, true)
	return functions
}

func GetAllCloudAgnosticFunctions() map[string]function.Function {
	return map[string]function.Function{
		"file":         getTerraformFunction("file"),
		"templatefile": getTerraformFunction("templatefile"),
	}
}

// GetAllFunctionsForVarEvaluation returns the only functions that are available when parsing variable blocks.
func GetAllFunctionsForVarEvaluation(cloud common.CloudProvider) map[string]function.Function {
	return map[string]function.Function{
		"cloud_specific_value": GetCloudSpecificValueFunction(cloud, false),
	}
}
