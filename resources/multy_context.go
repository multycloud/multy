package resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
)

type MultyContext struct {
	Resources map[string]Resource
}

func GetAllResources[T Resource](ctx MultyContext) []T {
	var result []T
	for _, r := range ctx.Resources {
		if casted, canCast := r.(T); canCast {
			result = append(result, casted)
		}
	}
	return result
}

func GetAllResourcesInCloud[T Resource](ctx MultyContext, cloud commonpb.CloudProvider) []T {
	var result []T
	for _, r := range ctx.Resources {
		if r.GetCloud() != cloud {
			continue
		}
		if casted, canCast := r.(T); canCast {
			result = append(result, casted)
		}
	}
	return result
}
