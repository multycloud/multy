package database

import (
	"fmt"

	"github.com/multycloud/multy/resources/common"
)

const AzureMariaDbResourceName = "azurerm_mariadb_server"

type AzureMariaDbServer struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mariadb_server"`
	AzureDatabase      `hcl:",squash"`
}

type AzureMariaDbVirtualNetworkRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mariadb_virtual_network_rule"`
	ServerName         string `hcl:"server_name,expr"`
	SubnetId           string `hcl:"subnet_id,expr"`
}

type AzureMariaDbFirewallRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mariadb_firewall_rule"`
	ServerName         string `hcl:"server_name,expr"`
	StartIpAddress     string `hcl:"start_ip_address"`
	EndIpAddress       string `hcl:"end_ip_address"`
}

func (db AzureMariaDbServer) GetServerName() string {
	return fmt.Sprintf("azurerm_mariadb_server.%s.name", db.ResourceId)
}
