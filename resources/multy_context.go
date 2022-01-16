package resources

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/validate"
)

type MultyContext struct {
	Resources map[string]CloudSpecificResource
	Location  string
}

func (ctx *MultyContext) GetResource(id string) (*CloudSpecificResource, error) {
	if r, ok := ctx.Resources[id]; ok {
		return &r, nil
	}
	return nil, fmt.Errorf("resource %s not found", id)
}

func (ctx *MultyContext) GetLocationFromCommonParams(commonParams *CommonResourceParams, cloud common.CloudProvider) string {
	location := ctx.Location
	if commonParams.Location != "" {
		location = commonParams.Location
	}

	if result, err := common.GetCloudLocation(location, cloud); err != nil {
		if location == commonParams.Location {
			commonParams.LogFatal(commonParams.ResourceId, "location", err.Error())
		} else {
			// TODO: throw a user error when validating global config
			validate.LogInternalError(err.Error())
		}
		return ""
	} else {
		return result
	}
}

func (ctx *MultyContext) GetLocation(specifiedLocation string, cloud common.CloudProvider) string {
	location := ctx.Location
	if specifiedLocation != "" {
		location = specifiedLocation
	}

	if result, err := common.GetCloudLocation(location, cloud); err != nil {
		validate.LogInternalError(err.Error())
		return ""
	} else {
		return result
	}
}
