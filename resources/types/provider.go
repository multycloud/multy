package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/resources/output/provider"
)

type Provider struct {
	Cloud             commonpb.CloudProvider
	Location          string
	IsDefaultProvider bool
	NumResources      int
	Credentials       *credspb.CloudCredentials
}

func (p *Provider) Translate() []any {
	if p.Cloud == commonpb.CloudProvider_AWS {
		return []any{provider.AwsProvider{
			ResourceName: provider.AwsResourceName,
			Region:       p.Location,
			Alias:        p.getAlias(),
			AccessKey:    p.Credentials.GetAwsCreds().GetAccessKey(),
			SecretKey:    p.Credentials.GetAwsCreds().GetSecretKey(),
		}}
	} else if p.Cloud == commonpb.CloudProvider_AZURE {
		if !p.IsDefaultProvider {
			return []any{}
		}
		return []any{provider.AzureProvider{
			ResourceName:   provider.AzureResourceName,
			Features:       provider.AzureProviderFeatures{},
			ClientId:       p.Credentials.GetAzureCreds().GetClientId(),
			SubscriptionId: p.Credentials.GetAzureCreds().GetSubscriptionId(),
			TenantId:       p.Credentials.GetAzureCreds().GetTenantId(),
			ClientSecret:   p.Credentials.GetAzureCreds().GetClientSecret(),
		}}
	}
	return nil
}

func (p *Provider) GetId() string {
	return fmt.Sprintf("%s.%s", p.Cloud, p.getAlias())
}

func (p *Provider) getAlias() string {
	if p.Cloud == commonpb.CloudProvider_AWS && !p.IsDefaultProvider {
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
