package vault

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultAccessPolicyService struct {
	Service services.Service[*resourcespb.VaultAccessPolicyArgs, *resourcespb.VaultAccessPolicyResource]
}

func (s VaultAccessPolicyService) Convert(resourceId string, args *resourcespb.VaultAccessPolicyArgs, state *output.TfState) (*resourcespb.VaultAccessPolicyResource, error) {

	return &resourcespb.VaultAccessPolicyResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		VaultId:          args.VaultId,
		Identity:         args.Identity,
		Access:           args.Access,
	}, nil
}

func NewVaultAccessPolicyService(database *db.Database) VaultAccessPolicyService {
	ni := VaultAccessPolicyService{
		Service: services.Service[*resourcespb.VaultAccessPolicyArgs, *resourcespb.VaultAccessPolicyResource]{
			Db:           database,
			Converters:   nil,
			ResourceName: "vault_access_policy",
		},
	}
	ni.Service.Converters = &ni
	return ni
}
