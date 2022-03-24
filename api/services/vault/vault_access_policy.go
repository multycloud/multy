package vault

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VaultAccessPolicyService struct {
	Service services.Service[*resources.CloudSpecificVaultAccessPolicyArgs, *resources.VaultAccessPolicyResource]
}

func (s VaultAccessPolicyService) Convert(resourceId string, args []*resources.CloudSpecificVaultAccessPolicyArgs, state *output.TfState) (*resources.VaultAccessPolicyResource, error) {
	var result []*resources.CloudSpecificVaultAccessPolicyResource
	for _, r := range args {
		result = append(result, &resources.CloudSpecificVaultAccessPolicyResource{
			CommonParameters: util.ConvertCommonParams(r.CommonParameters),
			VaultId:          r.VaultId,
			Identity:         r.Identity,
			Access:           r.Access,
		})
	}

	return &resources.VaultAccessPolicyResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func NewVaultAccessPolicyService(database *db.Database) VaultAccessPolicyService {
	ni := VaultAccessPolicyService{
		Service: services.Service[*resources.CloudSpecificVaultAccessPolicyArgs, *resources.VaultAccessPolicyResource]{
			Db:         database,
			Converters: nil,
		},
	}
	ni.Service.Converters = &ni
	return ni
}
