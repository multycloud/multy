package encoder

import (
	"bytes"
	"github.com/multy-dev/hclencoder"
	"log"
	"multy-go/decoder"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/types"
	"multy-go/util"
)

type WithProvider struct {
	Resource      any    `hcl:",squash"`
	ProviderAlias string `hcl:"provider" hcle:"omitempty"`
}

func Encode(decodedResources *decoder.DecodedResources) string {
	ctx := resources.MultyContext{Resources: decodedResources.Resources, Location: decodedResources.GlobalConfig.Location}
	var b bytes.Buffer

	providers := buildProviders(decodedResources, ctx)

	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		r.Resource.Validate(ctx)
		providerAlias := getProvider(providers, r, ctx).GetResourceId()
		for _, translated := range r.Translate(ctx) {
			var result any
			result = WithProvider{
				Resource:      translated,
				ProviderAlias: providerAlias,
			}
			// If not already wrapped in a tf block, assume it's a resource.
			if !output.IsTerraformBlock(translated) {
				result = output.ResourceWrapper{
					R: result,
				}
			}
			hcl, err := hclencoder.Encode(result)
			if err != nil {
				log.Fatal("unable to encode: ", err)
			}
			b.Write(hcl)
		}

	}

	for _, p := range flatten(providers) {
		for _, translatedProvider := range p.Translate() {
			hcl, err := hclencoder.Encode(providerWrapper{
				P: translatedProvider,
			})
			if err != nil {
				log.Fatal("unable to encode: ", err)
			}
			b.Write(hcl)
		}
	}

	return b.String()
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
		providerByLocation[util.MaxBy(providerByLocation, func(v *types.Provider) int {
			return v.NumResources
		})].IsDefaultProvider = true
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

type providerWrapper struct {
	P any `hcl:"provider"`
}
