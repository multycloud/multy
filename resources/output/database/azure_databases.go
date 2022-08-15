package database

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
)

type AzureDatabaseEngine interface {
	GetServerName() string
}

type AzureDbServer struct {
	*common.AzResource         `default:"name=azurerm_mysql_server"`
	AdministratorLogin         string
	AdministratorLoginPassword string
	SkuName                    string
	StorageMb                  int
	Version                    string
	Engine                     resourcespb.DatabaseEngine
	SubnetId                   string
}

type AzureDatabase struct {
	AdministratorLogin         string `hcl:"administrator_login" json:"administrator_login"`
	AdministratorLoginPassword string `hcl:"administrator_login_password" json:"administrator_login_password"`
	SkuName                    string `hcl:"sku_name" json:"sku_name"`
	StorageMb                  int    `hcl:"storage_mb" json:"storage_mb"`
	Version                    string `hcl:"version" json:"version"`
	SslEnforcementEnabled      bool   `hcl:"ssl_enforcement_enabled" json:"ssl_enforcement_enabled"`

	Fqdn    string `hcl:"fqdn" hcle:"omitempty"`
	Id      string `hcl:"id" hcle:"omitempty"`
	NameOut string `json:"name" hcle:"omitempty"`
}

func NewAzureDatabase(server AzureDbServer) []output.TfBlock {
	azResource := &common.AzResource{
		TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
		ResourceGroupName: server.ResourceGroupName,
		Name:              server.Name,
		Location:          server.Location,
	}
	azDatabase := AzureDatabase{
		AdministratorLogin:         server.AdministratorLogin,
		AdministratorLoginPassword: server.AdministratorLoginPassword,
		SkuName:                    server.SkuName,
		StorageMb:                  server.StorageMb,
		Version:                    server.Version,
		SslEnforcementEnabled:      false,
	}
	switch server.Engine {
	// TODO: move to flexible mysql server
	case resourcespb.DatabaseEngine_MYSQL:
		mysqlServer := AzureMySqlServer{
			AzResource:                   azResource,
			AzureDatabase:                azDatabase,
			SslMinimalTlsVersionEnforced: "TLSEnforcementDisabled",
		}

		resources := []output.TfBlock{mysqlServer}
		resources = append(
			resources, AzureMySqlVirtualNetworkRule{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
					ResourceGroupName: server.ResourceGroupName,
					Name:              server.Name,
				},
				ServerName: mysqlServer.GetServerName(),
				SubnetId:   server.SubnetId,
			},
		)

		resources = append(resources, AzureDbFirewallRule{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              "public",
			},
			ServerName:     mysqlServer.GetServerName(),
			StartIpAddress: "0.0.0.0",
			EndIpAddress:   "255.255.255.255",
		})

		return resources

	case resourcespb.DatabaseEngine_POSTGRES:
		postgresqlServer := AzurePostgreSqlServer{
			AzResource:                   azResource,
			AzureDatabase:                azDatabase,
			SslMinimalTlsVersionEnforced: "TLSEnforcementDisabled",
		}

		resources := []output.TfBlock{postgresqlServer}
		resources = append(
			resources, AzurePostgreSqlVirtualNetworkRule{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
					ResourceGroupName: server.ResourceGroupName,
					Name:              server.Name,
				},
				ServerName: postgresqlServer.GetServerName(),
				SubnetId:   server.SubnetId,
			},
		)

		resources = append(resources, AzurePostgreSqlFirewallRule{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              "public",
			},
			ServerName:     postgresqlServer.GetServerName(),
			StartIpAddress: "0.0.0.0",
			EndIpAddress:   "255.255.255.255",
		})

		return resources

	case resourcespb.DatabaseEngine_MARIADB:
		mariaDbServer := AzureMariaDbServer{
			AzResource:    azResource,
			AzureDatabase: azDatabase,
		}

		resources := []output.TfBlock{mariaDbServer}
		resources = append(
			resources, AzureMariaDbVirtualNetworkRule{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
					ResourceGroupName: server.ResourceGroupName,
					Name:              server.Name,
				},
				ServerName: mariaDbServer.GetServerName(),
				SubnetId:   server.SubnetId,
			},
		)

		resources = append(resources, AzureMariaDbFirewallRule{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              "public",
			},
			ServerName:     mariaDbServer.GetServerName(),
			StartIpAddress: "0.0.0.0",
			EndIpAddress:   "255.255.255.255",
		})

		return resources
	default:
		panic(fmt.Sprintf("unknown engine %s", server.Engine.String()))
		return nil
	}
}
