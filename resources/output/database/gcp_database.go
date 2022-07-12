package database

import "github.com/multycloud/multy/resources/common"

type GoogleSqlDatabaseInstance struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_sql_database_instance"`
	DatabaseVersion     string                              `hcl:"database_version"`
	Settings            []GoogleSqlDatabaseInstanceSettings `hcl:"settings,blocks"`
	DeletionProtection  bool                                `hcl:"deletion_protection"`

	// outputs
	PublicIpAddress string `json:"public_ip_address" hcle:"omitempty"`
}

type GoogleSqlDatabaseInstanceSettings struct {
	Tier             string                `hcl:"tier"`
	AvailabilityType string                `hcl:"availability_type"`
	DiskAutoResize   bool                  `hcl:"disk_autoresize"`
	DiskSize         int                   `hcl:"disk_size"`
	IpConfiguration  GoogleIpConfiguration `hcl:"ip_configuration"`
}

type GoogleIpConfiguration struct {
	AuthorizedNetworks []GoogleAuthorizedNetwork `hcl:"authorized_networks,blocks"`
}

type GoogleAuthorizedNetwork struct {
	Value string `hcl:"value"`
}

type GoogleSqlUser struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_sql_user"`
	Instance            string `hcl:"instance,expr"`
	Password            string `hcl:"password"`
}
