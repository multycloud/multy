package encoder

import (
	"bytes"
	"log"
	"multy-go/decoder"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/types"
	"sort"

	"github.com/multy-dev/hclencoder"
)

type WithProvider struct {
	Resource      interface{} `hcl:",squash"`
	ProviderAlias string      `hcl:"provider" hcle:"omitempty"`
}

func Encode(decodedResources *decoder.DecodedResources) string {
	ctx := resources.MultyContext{Resources: decodedResources.Resources, Location: decodedResources.GlobalConfig.Location}
	var b bytes.Buffer

	providers := buildProviders(decodedResources, ctx)

	for _, r := range sortResources(decodedResources) {
		r.Resource.Validate(ctx)
		providerAlias := getProvider(providers, r, ctx).GetResourceId()
		for _, translated := range r.Translate(ctx) {
			var result interface{}
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

	biggestProviders := map[common.CloudProvider]*types.Provider{}
	for cloud, providerByLocation := range providers {
		for _, provider := range providerByLocation {
			if biggest, ok := biggestProviders[cloud]; ok {
				if biggest.NumResources < provider.NumResources || (biggest.NumResources == provider.NumResources &&
					biggest.Location > provider.Location) {
					biggestProviders[cloud] = provider
					provider.IsDefaultProvider = true
					biggest.IsDefaultProvider = false
				}
			} else {
				biggestProviders[cloud] = provider
				provider.IsDefaultProvider = true
			}
		}
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
	sort.Slice(result, func(a, b int) bool {
		return result[a].GetId() < result[b].GetId()
	})
	return result
}

type providerWrapper struct {
	P interface{} `hcl:"provider"`
}

func sortResources(decodedResources *decoder.DecodedResources) []resources.CloudSpecificResource {
	var keys []string
	for k := range decodedResources.Resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result []resources.CloudSpecificResource
	for _, k := range keys {
		result = append(result, decodedResources.Resources[k])
	}
	return result
}
