package network_interface

import (
	"fmt"
	"multy-go/resources/common"
	"multy-go/validate"
)

const AzureResourceName = "azurerm_network_interface"

type AzureNetworkInterface struct {
	common.AzResource `hcl:",squash"`
	IpConfigurations  []AzureIpConfiguration `hcl:"ip_configuration,blocks"`
}

func (nic AzureNetworkInterface) GetResourceId(cloud common.CloudProvider) string {
	return fmt.Sprintf("%s_%s", nic.ResourceId, cloud)
}

func (nic AzureNetworkInterface) GetId(cloud common.CloudProvider) string {
	if cloud == common.AZURE {
		return fmt.Sprintf("${%s.%s.id}", AzureResourceName, nic.ResourceId)
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

type AzureIpConfiguration struct {
	Name                       string `hcl:"name"`
	PrivateIpAddressAllocation string `hcl:"private_ip_address_allocation"`
	SubnetId                   string `hcl:"subnet_id" hcle:"omitempty"`
	PublicIpAddressId          string `hcl:"public_ip_address_id,expr" hcle:"omitempty"`
	Primary                    bool   `hcl:"primary" hcle:"omitempty"`
}
