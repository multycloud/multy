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

func TranslateResources(decodedResources *DecodedResources, ctx resources.MultyContext) (map[resources.Resource][]output.TfBlock, []validate.ValidationError, error) {
	defaultTagProcessor := mhcl.DefaultTagProcessor{}

	translationCache := map[resources.Resource][]output.TfBlock{}

	errors := map[validate.ValidationError]bool{}
	for _, r := range util.GetSortedMapValues(decodedResources.Resources.ResourceMap) {
		var err error
		// we need to use a set here because errors are duplicated for multiple clouds
		validationErrors := r.Validate(ctx)
		for _, err := range validationErrors {
			errors[err] = true
		}
		if len(validationErrors) == 0 {
			translationCache[r], err = r.Translate(ctx)
			if err != nil {
				return translationCache, nil, err
			}
			for _, translated := range translationCache[r] {
				defaultTagProcessor.Process(translated)
			}
		}
	}

	return translationCache, maps.Keys(errors), nil
}

func getProvider(providers map[commonpb.CloudProvider]map[string]*types.Provider, r resources.Resource) *types.Provider {
	return providers[r.GetCloud()][r.GetCloudSpecificLocation()]
}

func buildProviders(r *DecodedResources, credentials *credspb.CloudCredentials) map[commonpb.CloudProvider]map[string]*types.Provider {
	providers := r.Providers

	for _, resource := range r.Resources.ResourceMap {
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
