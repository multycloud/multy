package functions

import (
	"github.com/zclconf/go-cty/cty/function"
	"multy-go/resources/common"
)

func GetAllFunctions(cloud common.CloudProvider) map[string]function.Function {
	return map[string]function.Function{
		"cloud_specific_value": GetCloudSpecificValueFunction(cloud, true),
	}
}

// GetAllFunctionsForVarEvaluation returns the only functions that are available when parsing variable blocks.
func GetAllFunctionsForVarEvaluation(cloud common.CloudProvider) map[string]function.Function {
	return map[string]function.Function{
		"cloud_specific_value": GetCloudSpecificValueFunction(cloud, false),
	}
}
