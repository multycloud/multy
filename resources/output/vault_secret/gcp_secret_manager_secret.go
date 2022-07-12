package vault_secret

import "github.com/multycloud/multy/resources/common"

type GoogleSecretManagerSecret struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_secret_manager_secret"`

	SecretId    string                               `hcl:"secret_id"`
	Replication GoogleSecretManagerSecretReplication `hcl:"replication"`
}

type GoogleSecretManagerSecretReplication struct {
	Automatic   bool                                            `hcl:"automatic" hcle:"omitempty"`
	UserManaged GoogleSecretManagerSecretReplicationUserManaged `hcl:"user_managed"  hcle:"omitempty"`
}

type GoogleSecretManagerSecretReplicationUserManaged struct {
	Replicas []GoogleSecretManagerSecretReplicationReplica `hcl:"replicas,blocks"`
}
type GoogleSecretManagerSecretReplicationReplica struct {
	Location string `hcl:"location"`
}

func NewManagedSecretReplication(locations ...string) GoogleSecretManagerSecretReplication {
	result := GoogleSecretManagerSecretReplication{
		Automatic:   false,
		UserManaged: GoogleSecretManagerSecretReplicationUserManaged{},
	}
	for _, location := range locations {
		result.UserManaged.Replicas = append(result.UserManaged.Replicas,
			GoogleSecretManagerSecretReplicationReplica{Location: location})
	}
	return result
}

type GoogleSecretManagerSecretVersion struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_secret_manager_secret_version"`

	SecretId   string `hcl:"secret,expr"`
	SecretData string `hcl:"secret_data"`
}
