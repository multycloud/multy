package network_interface_security_group_association

import "github.com/multycloud/multy/resources/common"

const AzureResourceName = "azurerm_network_interface_security_group_association"

type AzureNetworkInterfaceSecurityGroupAssociation struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_network_interface_security_group_association"`
	SecurityGroupId    string `hcl:"network_security_group_id,expr"`
	NetworkInterfaceId string `hcl:"network_interface_id,expr"`
}
