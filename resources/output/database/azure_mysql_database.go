package database

import (
	"fmt"
	"multy-go/resources/common"
)

const AzureMysqlResourceName = "azurerm_mysql_server"

type AzureMySqlServer struct {
	common.AzResource          `hcl:",squash"`
	AdministratorLogin         string `hcl:"administrator_login"`
	AdministratorLoginPassword string `hcl:"administrator_login_password"`
	SkuName                    string `hcl:"sku_name"`
	StorageMb                  int    `hcl:"storage_mb"`
	Version                    string `hcl:"version"`
	SslEnforcementEnabled      bool   `hcl:"ssl_enforcement_enabled"`
}

type AzureMySqlVirtualNetworkRule struct {
	common.AzResource `hcl:",squash"`
	ServerName        string `hcl:"server_name,expr"`
	SubnetId          string `hcl:"subnet_id"`
}

func (db AzureMySqlServer) GetServerName() string {
	return fmt.Sprintf("azurerm_mysql_server.%s.name", db.ResourceId)
}
