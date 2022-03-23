package virtual_machine

import (
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/resources/output"
)

type VirtualMachineService struct {
	Service services.Service[*resources.CloudSpecificVirtualMachineArgs, *resources.VirtualMachineResource]
}

func (s VirtualMachineService) Convert(resourceId string, args []*resources.CloudSpecificVirtualMachineArgs, state *output.TfState) (*resources.VirtualMachineResource, error) {
	var result []*resources.CloudSpecificVirtualMachineResource
	for _, r := range args {
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
		})
	}

	return &resources.VirtualMachineResource{
		CommonParameters: &common.CommonResourceParameters{ResourceId: resourceId},
		Resources:        result,
	}, nil
}

func (s VirtualMachineService) NewArg() *resources.CloudSpecificVirtualMachineArgs {
	return &resources.CloudSpecificVirtualMachineArgs{}
}

func (s VirtualMachineService) Nil() *resources.VirtualMachineResource {
	return nil
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
