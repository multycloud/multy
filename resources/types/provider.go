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
	GcpProject        string
}

func (p *Provider) Translate() []any {
	if p.Cloud == commonpb.CloudProvider_AWS {
		return []any{provider.AwsProvider{
			ResourceName: p.getResourceName(),
			Region:       p.Location,
			Alias:        p.getAlias(),
			AccessKey:    p.Credentials.GetAwsCreds().GetAccessKey(),
			SecretKey:    p.Credentials.GetAwsCreds().GetSecretKey(),
			SessionToken: p.Credentials.GetAwsCreds().GetSessionToken(),
		}}
	} else if p.Cloud == commonpb.CloudProvider_AZURE {
		if !p.IsDefaultProvider {
			return []any{}
		}
		return []any{provider.AzureProvider{
			ResourceName:   p.getResourceName(),
			Features:       provider.AzureProviderFeatures{},
			ClientId:       p.Credentials.GetAzureCreds().GetClientId(),
			SubscriptionId: p.Credentials.GetAzureCreds().GetSubscriptionId(),
			TenantId:       p.Credentials.GetAzureCreds().GetTenantId(),
			ClientSecret:   p.Credentials.GetAzureCreds().GetClientSecret(),
		}}
	} else if p.Cloud == commonpb.CloudProvider_GCP {
		return []any{provider.GcpProvider{
			ResourceName: p.getResourceName(),
			Region:       p.Location,
			Alias:        p.getAlias(),
			Credentials:  p.Credentials.GetGcpCreds().GetCredentials(),
		}}
	}
	return nil
}

func (p *Provider) getResourceName() string {
	switch p.Cloud {
	case commonpb.CloudProvider_AWS:
		return provider.AwsResourceName
	case commonpb.CloudProvider_AZURE:
		return provider.AzureResourceName
	case commonpb.CloudProvider_GCP:
		return provider.GcpResourceName
	default:
		panic(fmt.Sprintf("unhandled cloud %s", p.Cloud))
	}
}

func (p *Provider) GetId() string {
	return fmt.Sprintf("%s.%s", p.Cloud, p.getAlias())
}

func (p *Provider) getAlias() string {
	if p.Cloud == commonpb.CloudProvider_AWS || p.Cloud == commonpb.CloudProvider_GCP {
		return p.Location
	}

	return ""
}

func (p *Provider) GetResourceId() string {
	if p.getAlias() == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s", p.getResourceName(), p.getAlias())
}
