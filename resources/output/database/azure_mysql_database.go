package database

import (
	"fmt"
	"multy/resources/common"
)

const AzureMysqlResourceName = "azurerm_mysql_server"

type AzureMySqlServer struct {
	*common.AzResource         `hcl:",squash" default:"name=azurerm_mysql_server"`
	AdministratorLogin         string `hcl:"administrator_login"`
	AdministratorLoginPassword string `hcl:"administrator_login_password"`
	SkuName                    string `hcl:"sku_name"`
	StorageMb                  int    `hcl:"storage_mb"`
	Version                    string `hcl:"version"`
	SslEnforcementEnabled      bool   `hcl:"ssl_enforcement_enabled"`
}

type AzureMySqlVirtualNetworkRule struct {
	*common.AzResource `hcl:",squash" default:"name=azurerm_mysql_virtual_network_rule"`
	ServerName         string `hcl:"server_name,expr"`
	SubnetId           string `hcl:"subnet_id"`
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
