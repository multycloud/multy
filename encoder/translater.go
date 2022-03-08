package encoder

import (
	"golang.org/x/exp/maps"
	"multy-go/decoder"
	"multy-go/mhcl"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/types"
	"multy-go/util"
	"multy-go/validate"
)

func TranslateResources(decodedResources *decoder.DecodedResources, ctx resources.MultyContext) map[resources.CloudSpecificResource][]output.TfBlock {
	defaultTagProcessor := mhcl.DefaultTagProcessor{}

	translationCache := map[resources.CloudSpecificResource][]output.TfBlock{}

	errors := map[validate.ValidationError]bool{}
	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		translationCache[r] = r.Translate(ctx)
		for _, translated := range translationCache[r] {
			defaultTagProcessor.Process(translated)
		}
		// we need to use a set here because errors are duplicated for multiple clouds
		for _, err := range r.Resource.Validate(ctx, r.Cloud) {
			errors[err] = true
		}
	}

	if len(errors) != 0 {
		validate.PrintAllAndExit(maps.Keys(errors))
	}
	return translationCache
}

func getProvider(providers map[common.CloudProvider]map[string]*types.Provider, r resources.CloudSpecificResource, ctx resources.MultyContext) *types.Provider {
	return providers[r.Cloud][r.GetLocation(ctx)]
}

func buildProviders(r *decoder.DecodedResources, ctx resources.MultyContext) map[common.CloudProvider]map[string]*types.Provider {
	providers := map[common.CloudProvider]map[string]*types.Provider{}

	for _, resource := range r.Resources {
		if _, ok := providers[resource.Cloud]; !ok {
			providers[resource.Cloud] = map[string]*types.Provider{}
		}
		if _, ok := providers[resource.Cloud][resource.GetLocation(ctx)]; !ok {
			provider := &types.Provider{
				Cloud:        resource.Cloud,
				Location:     resource.GetLocation(ctx),
				NumResources: 1,
			}
			providers[resource.Cloud][resource.GetLocation(ctx)] = provider
		} else {
			providers[resource.Cloud][resource.GetLocation(ctx)].NumResources += 1
		}
	}

	for _, providerByLocation := range providers {
		providerByLocation[util.MaxBy(
			providerByLocation, func(v *types.Provider) int {
				return v.NumResources
			},
		)].IsDefaultProvider = true
	}

	return providers
}

func flatten(p map[common.CloudProvider]map[string]*types.Provider) []*types.Provider {
	var result []*types.Provider
	for _, providers := range p {
		for _, provider := range providers {
			result = append(result, provider)
		}
	}
	util.SortResourcesById(result, func(p *types.Provider) string { return p.GetId() })
	return result
}
