package encoder

import (
	"bytes"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
	"log"
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

	translatedResources := TranslateResources(decodedResources, ctx)
	providers := buildProviders(decodedResources, ctx)

	var b bytes.Buffer
	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		providerAlias := getProvider(providers, r, ctx).GetResourceId()
		for _, translated := range translatedResources[r] {
			var result output.TfBlock
			result = WithProvider{
				Resource:      translated,
				ProviderAlias: providerAlias,
			}

			for _, dep := range r.Resource.GetDependencies(ctx) {
				for _, translatedDep := range translatedResources[dep] {
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
		ResourceId string    `hcl:",key"`
		Value      cty.Value `hcl:"value"`
	}

	for _, outputVal := range util.GetSortedMapValues(decodedResources.Outputs) {
		hclOutput, err := hclencoder.Encode(
			outputWrapper{
				O: outputStruct{
					ResourceId: outputVal.OutputId,
					Value:      outputVal.Value,
				},
			},
		)
		if err != nil {
			validate.LogFatalWithSourceRange(outputVal.DefinitionRange, "unable to encode output: %s", err.Error())
		}
		b.Write(hclOutput)
	}

	return b.String()
}

type providerWrapper struct {
	P any `hcl:"provider"`
}
type outputWrapper struct {
	O any `hcl:"output"`
}
