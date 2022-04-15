package virtual_machine

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/output/virtual_machine"
)

type VirtualMachineService struct {
	Service services.Service[*resourcespb.VirtualMachineArgs, *resourcespb.VirtualMachineResource]
}

func (s VirtualMachineService) Convert(resourceId string, args *resourcespb.VirtualMachineArgs, state *output.TfState) (*resourcespb.VirtualMachineResource, error) {
	var err error
	ip := "dryrun"
	identityId := "dryrun"
	if !flags.DryRun {
		ip, err = getPublicIp(resourceId, state, args.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
		identityId, err = getIdentityId(resourceId, state, args.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
	}

	return &resourcespb.VirtualMachineResource{
		CommonParameters:        util.ConvertCommonParams(resourceId, args.CommonParameters),
		Name:                    args.Name,
		OperatingSystem:         args.OperatingSystem,
		NetworkInterfaceIds:     args.NetworkInterfaceIds,
		NetworkSecurityGroupIds: args.NetworkSecurityGroupIds,
		VmSize:                  args.VmSize,
		UserData:                args.UserData,
		SubnetId:                args.SubnetId,
		PublicSshKey:            args.PublicSshKey,
		PublicIpId:              args.PublicIpId,
		GeneratePublicIp:        args.GeneratePublicIp,

		PublicIp:   ip,
		IdentityId: identityId,
	}, nil
}

func getPublicIp(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(virtual_machine.AwsEC2{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["public_ip"].(string), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(public_ip.AzurePublicIp{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["ip_address"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func getIdentityId(resourceId string, state *output.TfState, cloud commonpb.CloudProvider) (string, error) {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		values, err := state.GetValues(iam.AwsIamRole{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["id"].(string), nil
	case commonpb.CloudProvider_AZURE:
		values, err := state.GetValues(virtual_machine.AzureVirtualMachine{}, resourceId)
		if err != nil {
			return "", err
		}
		return values["identity"].([]interface{})[0].(map[string]interface{})["principal_id"].(string), nil
	}

	return "", fmt.Errorf("unknown cloud: %s", cloud.String())
}

func NewVirtualMachineService(database *db.Database) VirtualMachineService {
	VirtualMachine := VirtualMachineService{
		Service: services.Service[*resourcespb.VirtualMachineArgs, *resourcespb.VirtualMachineResource]{
			Db:         database,
			Converters: nil,
		},
	}
	VirtualMachine.Service.Converters = &VirtualMachine
	return VirtualMachine
}
