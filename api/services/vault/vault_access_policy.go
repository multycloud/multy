package vault

import (
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultAccessPolicyService struct {
	Service services.Service[*resources.VaultAccessPolicyArgs, *resources.VaultAccessPolicyResource]
}

func (s VaultAccessPolicyService) Convert(resourceId string, args *resources.VaultAccessPolicyArgs, state *output.TfState) (*resources.VaultAccessPolicyResource, error) {

	return &resources.VaultAccessPolicyResource{
		CommonParameters: util.ConvertCommonChildParams(resourceId, args.CommonParameters),
		VaultId:          args.VaultId,
		Identity:         args.Identity,
		Access:           args.Access,
	}, nil
}

func NewVaultAccessPolicyService(database *db.Database) VaultAccessPolicyService {
	ni := VaultAccessPolicyService{
		Service: services.Service[*resources.VaultAccessPolicyArgs, *resources.VaultAccessPolicyResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
