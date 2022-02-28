package encoder

import (
	"bytes"
	"github.com/multy-dev/hclencoder"
	"github.com/zclconf/go-cty/cty"
	"log"
	"multy-go/decoder"
	"multy-go/mhcl"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/types"
	"multy-go/util"
)

type WithProvider struct {
	Resource      output.TfBlock `hcl:",squash"`
	ProviderAlias string         `hcl:"provider" hcle:"omitempty"`
}

func (w WithProvider) GetFullResourceRef() string {
	return w.Resource.GetFullResourceRef()
}

func (w WithProvider) GetBlockType() string {
	return w.Resource.GetBlockType()
}

func (w WithProvider) AddDependency(s string) {
	w.Resource.AddDependency(s)
}

func Encode(decodedResources *decoder.DecodedResources) string {
	ctx := resources.MultyContext{Resources: decodedResources.Resources, Location: decodedResources.GlobalConfig.Location}
	var b bytes.Buffer
	defaultTagProcessor := mhcl.DefaultTagProcessor{}

	providers := buildProviders(decodedResources, ctx)

	translationCache := map[resources.CloudSpecificResource][]output.TfBlock{}

	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		translationCache[r] = r.Translate(ctx)
		for _, translated := range translationCache[r] {
			defaultTagProcessor.Process(translated)
		}
		r.Resource.Validate(ctx)
	}

	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		providerAlias := getProvider(providers, r, ctx).GetResourceId()
		for _, translated := range translationCache[r] {
			var result output.TfBlock
			result = WithProvider{
				Resource:      translated,
				ProviderAlias: providerAlias,
			}

			for _, dep := range r.Resource.GetDependencies(ctx) {
				for _, translatedDep := range translationCache[dep] {
					translated.AddDependency(translatedDep.GetFullResourceRef())
				}
			}

			// If not already wrapped in a tf block, assume it's a resource.
			hcl, err := hclencoder.Encode(output.WrapWithBlockType(result))
			if err != nil {
				log.Fatal("unable to encode: ", err)
			}
			b.Write(hcl)
		}

	}

	for _, p := range flatten(providers) {
		for _, translatedProvider := range p.Translate() {
			hcl, err := hclencoder.Encode(
				providerWrapper{
					P: translatedProvider,
				},
			)
			if err != nil {
				log.Fatal("unable to encode: ", err)
			}
			b.Write(hcl)
		}
	}

	type outputStruct struct {
		ResourceId string `hcl:",key"`
		Value      string `hcl:"value"`
	}

	for outputId, outputVal := range decodedResources.Outputs {
		if outputVal.Type() != cty.String {
			log.Fatalf("non-string outputs are currently not supported")
		}
		hclOutput, err := hclencoder.Encode(
			outputWrapper{
				O: outputStruct{
					ResourceId: outputId,
					Value:      outputVal.AsString(),
				},
			},
		)
		if err != nil {
			log.Fatal("unable to encode: ", err)
		}
		b.Write(hclOutput)
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

type providerWrapper struct {
	P any `hcl:"provider"`
}
type outputWrapper struct {
	O any `hcl:"output"`
}
