package network_security_group

import (
	"multy-go/resources/common"
)

const AzureNetworkSecurityGroupResourceName = "azurerm_network_security_group"

type AzureNetworkSecurityGroup struct {
	Nsg                  AzureNsg                  `hcl:"resource"`
	SubnetNsgAssociation AzureSubnetNsgAssociation `hcl:"resource"`
}

type AzureNsg struct {
	common.AzResource `hcl:",squash"`
	Rules             []AzureRule `hcl:"security_rule,blocks"`
}

type AzureRule struct {
	Name                     string `hcl:"name"`
	Protocol                 string `hcl:"protocol"`
	Priority                 int    `hcl:"priority"`
	Access                   string `hcl:"access"`
	SourcePortRange          string `hcl:"source_port_range"`
	SourceAddressPrefix      string `hcl:"source_address_prefix"`
	DestinationPortRange     string `hcl:"destination_port_range"`
	DestinationAddressPrefix string `hcl:"destination_address_prefix"`
	Direction                string `hcl:"direction"`
}

type AzureSubnetNsgAssociation struct {
	ResourceName string `hcl:",key"`
	ResourceId   string `hcl:",key"`
	SubnetId     string `hcl:"subnet_id"`
	NsgId        string `hcl:"network_security_group_id,expr"`
}

const AzureNicNsgAssociation = "azurerm_network_interface_security_group_association"

type AzureNetworkInterfaceSecurityGroupAssociation struct {
	ResourceName           string `hcl:",key"`
	ResourceId             string `hcl:",key"`
	NetworkInterfaceId     string `hcl:"network_interface_id,expr"`
	NetworkSecurityGroupId string `hcl:"network_security_group_id,expr"`
}
