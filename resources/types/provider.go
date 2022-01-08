package types

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/resources/output/provider"
)

type Provider struct {
	Cloud             common.CloudProvider
	Location          string
	IsDefaultProvider bool
	NumResources      int
}

func (p *Provider) Translate() []interface{} {
	if p.Cloud == common.AWS {
		return []interface{}{provider.AwsProvider{
			ResourceName: provider.AwsResourceName,
			Region:       p.Location,
			Alias:        p.getAlias(),
		}}
	} else if p.Cloud == common.AZURE {
		if !p.IsDefaultProvider {
			return []interface{}{}
		}
		return []interface{}{provider.AzureProvider{
			ResourceName: provider.AzureResourceName,
			Features:     provider.AzureProviderFeatures{},
		}}
	}
	return nil
}

func (p *Provider) GetId() string {
	return fmt.Sprintf("%s.%s", p.Cloud, p.getAlias())
}

func (p *Provider) getAlias() string {
	if p.Cloud == common.AWS && !p.IsDefaultProvider {
		return p.Location
	} else {
		return ""
	}
}

func (p *Provider) GetResourceId() string {
	if p.getAlias() == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s", p.Cloud, p.getAlias())
}
