package provider

const AzureResourceName = "azurerm"

type AzureProvider struct {
	ResourceName string                `hcl:",key"`
	Features     AzureProviderFeatures `hcl:"features"`
}

type AzureProviderFeatures struct {
	ResourceGroup string `hcl:"resource_group" hcle:"omitempty"`
}
