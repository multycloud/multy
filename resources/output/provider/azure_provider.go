package provider

const AzureResourceName = "azurerm"

type AzureProvider struct {
	ResourceName   string                `hcl:",key"`
	Features       AzureProviderFeatures `hcl:"features"`
	ClientId       string                `hcl:"client_id" hcle:"omitempty"`
	SubscriptionId string                `hcl:"subscription_id" hcle:"omitempty"`
	TenantId       string                `hcl:"tenant_id" hcle:"omitempty"`
	ClientSecret   string                `hcl:"client_secret" hcle:"omitempty"`
}

type AzureProviderFeatures struct {
	ResourceGroup string `hcl:"resource_group" hcle:"omitempty"`
}
