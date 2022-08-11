package network_security_group

import (
	"github.com/multycloud/multy/resources/common"
)

const AzureNetworkSecurityGroupResourceName = "azurerm_network_security_group"

type AzureNsg struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_network_security_group"`
	Rules              []AzureRule `hcl:"security_rule,blocks" json:"security_rule"`
}

type AzureRule struct {
	Name                     string `hcl:"name" json:"name"`
	Protocol                 string `hcl:"protocol" json:"protocol"`
	Priority                 int    `hcl:"priority" json:"priority"`
	Access                   string `hcl:"access" json:"access"`
	SourcePortRange          string `hcl:"source_port_range" json:"source_port_range"`
	SourceAddressPrefix      string `hcl:"source_address_prefix" json:"source_address_prefix"`
	DestinationPortRange     string `hcl:"destination_port_range" json:"destination_port_range"`
	DestinationAddressPrefix string `hcl:"destination_address_prefix" json:"destination_address_prefix"`
	Direction                string `hcl:"direction" json:"direction"`
}

type AzureSubnetNsgAssociation struct {
	ResourceName string `hcl:",key"`
	ResourceId   string `hcl:",key"`
	SubnetId     string `hcl:"subnet_id,expr"`
	NsgId        string `hcl:"network_security_group_id,expr"`
}

type AzureNetworkInterfaceSecurityGroupAssociation struct {
	*common.AzResource     `hcl:",squash" default:"name=azurerm_network_interface_security_group_association"`
	NetworkInterfaceId     string `hcl:"network_interface_id,expr"`
	NetworkSecurityGroupId string `hcl:"network_security_group_id,expr"`
}
