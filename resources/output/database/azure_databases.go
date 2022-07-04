package database

import (
	"strings"

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
	Engine                     string
	SubnetId                   string
}

func NewAzureDatabase(server AzureDbServer) []output.TfBlock {
	switch strings.ToLower(server.Engine) {
	// TODO: move to flexible mysql server
	case "mysql":
		mysqlServer := AzureMySqlServer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              server.Name,
				Location:          server.Location,
			},
			AdministratorLogin:           server.AdministratorLogin,
			AdministratorLoginPassword:   server.AdministratorLoginPassword,
			SkuName:                      server.SkuName,
			StorageMb:                    server.StorageMb,
			Version:                      server.Version,
			SslEnforcementEnabled:        false,
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

	case "postgres":
		postgresqlServer := AzurePostgreSqlServer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              server.Name,
				Location:          server.Location,
			},
			AdministratorLogin:           server.AdministratorLogin,
			AdministratorLoginPassword:   server.AdministratorLoginPassword,
			SkuName:                      server.SkuName,
			StorageMb:                    server.StorageMb,
			Version:                      server.Version,
			SslEnforcementEnabled:        false,
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

	case "mariadb":
		mariaDbServer := AzureMariaDbServer{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: server.ResourceId},
				ResourceGroupName: server.ResourceGroupName,
				Name:              server.Name,
				Location:          server.Location,
			},
			AdministratorLogin:         server.AdministratorLogin,
			AdministratorLoginPassword: server.AdministratorLoginPassword,
			SkuName:                    server.SkuName,
			StorageMb:                  server.StorageMb,
			Version:                    server.Version,
			SslEnforcementEnabled:      false,
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

	}
	return nil
}
