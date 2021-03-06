package encoder

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/mhcl"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"golang.org/x/exp/maps"
)

type cloudSpecificResources map[string]resources.CloudSpecificResourceTranslator

func TranslateResources(decodedResources cloudSpecificResources, ctx resources.MultyContext) (map[string][]output.TfBlock, []validate.ValidationError, error) {
	defaultTagProcessor := mhcl.DefaultTagProcessor{}

	translationCache := map[string][]output.TfBlock{}

	errors := map[validate.ValidationError]bool{}
	for _, cr := range util.GetSortedMapValues(decodedResources) {
		var err error
		// we need to use a set here because errors are duplicated for multiple clouds
		validationErrors := cr.Validate(ctx)
		for _, err := range validationErrors {
			errors[err] = true
		}
		if len(validationErrors) == 0 {
			translationCache[cr.GetResourceId()], err = cr.Translate(ctx)
			if err != nil {
				return translationCache, nil, err
			}
			for _, translated := range translationCache[cr.GetResourceId()] {
				defaultTagProcessor.Process(translated)
			}
		}
	}

	return translationCache, maps.Keys(errors), nil
}

func getProvider(providers map[commonpb.CloudProvider]map[string]*types.Provider, r resources.Resource) *types.Provider {
	return providers[r.GetCloud()][r.GetCloudSpecificLocation()]
}

type providersMap map[commonpb.CloudProvider]map[string]*types.Provider

func buildProviders(providers providersMap, r cloudSpecificResources, credentials *credspb.CloudCredentials) providersMap {
	for _, resource := range r {
		if _, ok := providers[resource.GetCloud()]; !ok {
			providers[resource.GetCloud()] = map[string]*types.Provider{}
		}
		location := resource.GetCloudSpecificLocation()
		if _, ok := providers[resource.GetCloud()][location]; !ok {
			provider := &types.Provider{
				Cloud:        resource.GetCloud(),
				Location:     location,
				NumResources: 1,
			}
			providers[resource.GetCloud()][location] = provider
		} else {
			providers[resource.GetCloud()][location].NumResources += 1
		}
	}

	for _, providerByLocation := range providers {
		providerByLocation[util.MaxBy(
			providerByLocation, func(v *types.Provider) int {
				return v.NumResources
			},
		)].IsDefaultProvider = true
	}

	for _, perRegion := range providers {
		for _, p := range perRegion {
			p.Credentials = credentials
		}
	}

	return providers
}

func flatten(p map[commonpb.CloudProvider]map[string]*types.Provider) []*types.Provider {
	var result []*types.Provider
	for _, providers := range p {
		for _, provider := range providers {
			result = append(result, provider)
		}
	}
	util.SortResourcesById(result, func(p *types.Provider) string { return p.GetId() })
	return result
}
