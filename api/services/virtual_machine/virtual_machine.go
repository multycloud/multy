package virtual_machine

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/output/virtual_machine"
	common_util "github.com/multycloud/multy/util"
)

type VirtualMachineService struct {
	Service services.Service[*resources.CloudSpecificVirtualMachineArgs, *resources.VirtualMachineResource]
}

func (s VirtualMachineService) Convert(resourceId string, args []*resources.CloudSpecificVirtualMachineArgs, state *output.TfState) (*resources.VirtualMachineResource, error) {
	var result []*resources.CloudSpecificVirtualMachineResource
	for _, r := range args {
		ip, err := getPublicIp(resourceId, state, r.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
		identityId, err := getIdentityId(resourceId, state, r.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
		result = append(result, &resources.CloudSpecificVirtualMachineResource{
			CommonParameters:        util.ConvertCommonParams(r.CommonParameters),
			Name:                    r.Name,
			OperatingSystem:         r.OperatingSystem,
			NetworkInterfaceIds:     r.NetworkInterfaceIds,
			NetworkSecurityGroupIds: r.NetworkSecurityGroupIds,
			VmSize:                  r.VmSize,
			UserData:                r.UserData,
			SubnetId:                r.SubnetId,
			PublicSshKey:            r.PublicSshKey,
			PublicIpId:              r.PublicIpId,
			GeneratePublicIp:        r.GeneratePublicIp,

			PublicIp:   ip,
			IdentityId: identityId,
		})
	}

	return &resources.VirtualMachineResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func getPublicIp(resourceId string, state *output.TfState, cloud common.CloudProvider) (string, error) {
	rId := common_util.GetTfResourceId(resourceId, cloud.String())
	switch cloud {
	case common.CloudProvider_AWS:
		values, err := state.GetValues(virtual_machine.AwsEC2{}, rId)
		if err != nil {
			return "", err
		}
		return values["public_ip"].(string), nil
	case common.CloudProvider_AZURE:
		values, err := state.GetValues(public_ip.AzurePublicIp{}, rId)
		if err != nil {
			return "", err
		}
		return values["ip_address"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func getIdentityId(resourceId string, state *output.TfState, cloud common.CloudProvider) (string, error) {
	rId := common_util.GetTfResourceId(resourceId, cloud.String())
	switch cloud {
	case common.CloudProvider_AWS:
		values, err := state.GetValues(iam.AwsIamRole{}, rId)
		if err != nil {
			return "", err
		}
		return values["id"].(string), nil
	case common.CloudProvider_AZURE:
		values, err := state.GetValues(virtual_machine.AzureVirtualMachine{}, rId)
		if err != nil {
			return "", err
		}
		return values["identity"].([]interface{})[0].(map[string]interface{})["principal_id"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func NewVirtualMachineService(database *db.Database) VirtualMachineService {
	VirtualMachine := VirtualMachineService{
		Service: services.Service[*resources.CloudSpecificVirtualMachineArgs, *resources.VirtualMachineResource]{
			Db:         database,
			Converters: nil,
		},
	}
	VirtualMachine.Service.Converters = &VirtualMachine
	return VirtualMachine
}
