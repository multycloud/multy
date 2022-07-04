package iam

import "github.com/multycloud/multy/resources/common"

type GoogleServiceAccount struct {
	*common.GcpResource `hcl:",squash"  default:"name=google_service_account"`
	AccountId           string `hcl:"account_id"`
	DisplayName         string `hcl:"display_name"`
}
