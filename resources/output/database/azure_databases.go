package database

import (
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"strconv"
	"strings"
)

type AzureDbServer struct {
	*common.AzResource         `default:"name=azurerm_mysql_server"`
	AdministratorLogin         string
	AdministratorLoginPassword string
	SkuName                    string
	StorageMb                  int
	Version                    string
	Engine                     string
	SubnetIds                  []string
}

func NewAzureDatabase(server AzureDbServer) []output.TfBlock {
	switch strings.ToLower(server.Engine) {
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
		for i, subnetId := range server.SubnetIds {
			resources = append(
				resources, AzureMySqlVirtualNetworkRule{
					AzResource: &common.AzResource{
						TerraformResource: output.TerraformResource{ResourceId: server.ResourceId + strconv.Itoa(i)},
						ResourceGroupName: server.ResourceGroupName,
						Name:              server.Name + strconv.Itoa(i),
					},
					ServerName: mysqlServer.GetServerName(),
					SubnetId:   subnetId,
				},
			)
		}

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
	}
	return nil
}
