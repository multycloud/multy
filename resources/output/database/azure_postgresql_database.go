package database

import (
	"fmt"

	"github.com/multycloud/multy/resources/common"
)

const AzurePostgresqlResourceName = "azurerm_postgresql_server"

type AzurePostgreSqlServer struct {
	*common.AzResource           `hcl:",squash" default:"name=azurerm_postgresql_server"`
	AzureDatabase                `hcl:",squash"`
	SslMinimalTlsVersionEnforced string `hcl:"ssl_minimal_tls_version_enforced"`
}

type AzurePostgreSqlVirtualNetworkRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_postgresql_virtual_network_rule"`
	ServerName         string `hcl:"server_name,expr"`
	SubnetId           string `hcl:"subnet_id,expr"`
}

type AzurePostgreSqlFirewallRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_postgresql_firewall_rule"`
	ServerName         string `hcl:"server_name,expr"`
	StartIpAddress     string `hcl:"start_ip_address"`
	EndIpAddress       string `hcl:"end_ip_address"`
}

func (db AzurePostgreSqlServer) GetServerName() string {
	return fmt.Sprintf("azurerm_postgresql_server.%s.name", db.ResourceId)
}
