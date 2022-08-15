package database

import (
	"fmt"

	"github.com/multycloud/multy/resources/common"
)

const AzureMysqlResourceName = "azurerm_mysql_server"

type AzureMySqlServer struct {
	*common.AzResource           `hcl:",squash" default:"name=azurerm_mysql_server"`
	AzureDatabase                `hcl:",squash"`
	SslMinimalTlsVersionEnforced string `hcl:"ssl_minimal_tls_version_enforced"`
}

type AzureMySqlVirtualNetworkRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mysql_virtual_network_rule"`
	ServerName         string `hcl:"server_name,expr"`
	SubnetId           string `hcl:"subnet_id,expr"`
}

type AzureDbFirewallRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mysql_firewall_rule"`
	ServerName         string `hcl:"server_name,expr"`
	StartIpAddress     string `hcl:"start_ip_address"`
	EndIpAddress       string `hcl:"end_ip_address"`
}

func (db AzureMySqlServer) GetServerName() string {
	return fmt.Sprintf("azurerm_mysql_server.%s.name", db.ResourceId)
}
