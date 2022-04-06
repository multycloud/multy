package encoder

import (
	"bytes"
	"github.com/hashicorp/hcl/v2"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
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

type DecodedResources struct {
	Resources map[string]resources.Resource
	Outputs   map[string]DecodedOutput
	Providers map[commonpb.CloudProvider]map[string]*types.Provider
}

type DecodedOutput struct {
	OutputId        string
	Value           cty.Value
	DefinitionRange hcl.Range
}

func Encode(decodedResources *DecodedResources, credentials *credspb.CloudCredentials) (string, []validate.ValidationError, error) {
	ctx := resources.MultyContext{Resources: decodedResources.Resources}

	translatedResources, errs, err := TranslateResources(decodedResources, ctx)
	if len(errs) > 0 || err != nil {
		return "", errs, err
	}
	providers := buildProviders(decodedResources, credentials)

	var b bytes.Buffer
	for _, r := range util.GetSortedMapValues(decodedResources.Resources) {
		providerAlias := getProvider(providers, r).GetResourceId()
		for _, translated := range translatedResources[r] {
			var result output.TfBlock
			result = WithProvider{
				Resource:      translated,
				ProviderAlias: providerAlias,
			}

			// If not already wrapped in a tf block, assume it's a resource.
			block, err := output.WrapWithBlockType(result)
			if err != nil {
				return "", nil, err
			}
			hclStr, err := hclencoder.Encode(block)
			if err != nil {
				return "", nil, errors.InternalServerErrorWithMessage("unexpected error encoding resource", err)
			}
			b.Write(hclStr)
		}
	}

	for _, p := range flatten(providers) {
		for _, translatedProvider := range p.Translate() {
			hclStr, err := hclencoder.Encode(
				providerWrapper{
					P: translatedProvider,
				},
			)
			if err != nil {
				return "", nil, errors.InternalServerErrorWithMessage("unexpected error encoding providers", err)
			}
			b.Write(hclStr)
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
			if err != nil {
				return "", nil, errors.InternalServerErrorWithMessage("unexpected error encoding outputs", err)
			}
		}
		b.Write(hclOutput)
	}

	return b.String(), nil, nil
}

type providerWrapper struct {
	P any `hcl:"provider"`
}
type outputWrapper struct {
	O any `hcl:"output"`
}
