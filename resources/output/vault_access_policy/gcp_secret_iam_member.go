package vault_access_policy

import "github.com/multycloud/multy/resources/common"

const SecretAccessorRole = "roles/secretmanager.secretAccessor"
const SecretWriterRole = "roles/secretmanager.secretVersionManager"
const SecretOwnerRole = "roles/secretmanager.admin"

type GoogleSecretManagerSecretIamMember struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_secret_manager_secret_iam_member"`

	SecretId string `hcl:"secret_id,expr"`
	Role     string `hcl:"role" json:"role"`
	Member   string `hcl:"member" json:"member"`
}
