package azure_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/output/terraform"
	"github.com/multycloud/multy/resources/output/virtual_machine"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
	"regexp"
)

type AzureVirtualMachine struct {
	*types.VirtualMachine
}

func InitVirtualMachine(vn *types.VirtualMachine) resources.ResourceTranslator[*resourcespb.VirtualMachineResource] {
	return AzureVirtualMachine{vn}
}

func (r AzureVirtualMachine) FromState(state *output.TfState) (*resourcespb.VirtualMachineResource, error) {
	out := &resourcespb.VirtualMachineResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
		},
		Name:                    r.Args.Name,
		NetworkInterfaceIds:     r.Args.NetworkInterfaceIds,
		NetworkSecurityGroupIds: r.Args.NetworkSecurityGroupIds,
		VmSize:                  r.Args.VmSize,
		UserDataBase64:          r.Args.UserDataBase64,
		SubnetId:                r.Args.SubnetId,
		PublicSshKey:            r.Args.PublicSshKey,
		PublicIpId:              r.Args.PublicIpId,
		GeneratePublicIp:        r.Args.GeneratePublicIp,
		ImageReference:          r.Args.ImageReference,
		AwsOverride:             r.Args.AwsOverride,
		AzureOverride:           r.Args.AzureOverride,
		GcpOverride:             r.Args.GcpOverride,
		AvailabilityZone:        r.Args.AvailabilityZone,
		IdentityId:              "dryrun",
	}

	if flags.DryRun {
		return out, nil
	}

	vmResource, err := output.GetParsedById[virtual_machine.AzureVirtualMachine](state, r.ResourceId)
	if err != nil {
		return nil, err
	}
	out.IdentityId = vmResource.Identities[0].PrincipalId

	out.AzureOutputs = &resourcespb.VirtualMachineAzureOutputs{
		VirtualMachineId: vmResource.ResourceId,
	}

	if r.Args.GeneratePublicIp {
		ipResource, err := output.GetParsedById[public_ip.AzurePublicIp](state, r.ResourceId)
		if err != nil {
			return nil, err
		}
		out.PublicIp = ipResource.IpAddress
		out.AzureOutputs.PublicIpId = ipResource.ResourceId
	}

	if stateResource, exists, err := output.MaybeGetParsedById[network_interface.AzureNetworkInterface](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AzureOutputs.NetworkInterfaceId = stateResource.ResourceId
	}

	return out, nil
}

func (r AzureVirtualMachine) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	subnetId, err := resources.GetMainOutputId(AzureSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}
	// TODO validate that NIC is on the same VNET
	var azResources []output.TfBlock
	rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
	nicIds, err := util.MapSliceValuesErr(r.NetworkInterface, func(v *types.NetworkInterface) (string, error) {
		return resources.GetMainOutputId(AzureNetworkInterface{v})
	})
	if err != nil {
		return nil, err
	}

	if len(r.NetworkInterface) == 0 {
		nic := network_interface.AzureNetworkInterface{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rgName,
				r.GetCloudSpecificLocation(),
			),
			IpConfigurations: []network_interface.AzureIpConfiguration{{
				Name:                       "internal", // this name shouldn't be r.name
				PrivateIpAddressAllocation: "Dynamic",
				SubnetId:                   subnetId,
				Primary:                    true,
			}},
		}

		if r.Args.GeneratePublicIp {
			pIp := public_ip.AzurePublicIp{
				AzResource: common.NewAzResource(
					r.ResourceId, r.Args.Name, rgName,
					r.GetCloudSpecificLocation(),
				),
				// TODO: this should be Dynamic, and then use the data source to get the public IP
				AllocationMethod: "Static",
				Sku:              "Standard",
			}
			nic.IpConfigurations = []network_interface.AzureIpConfiguration{{
				Name:                       "external", // this name shouldn't be r.name
				PrivateIpAddressAllocation: "Dynamic",
				SubnetId:                   subnetId,
				PublicIpAddressId:          pIp.GetId(),
				Primary:                    true,
			}}
			azResources = append(azResources, &pIp)
		}
		azResources = append(azResources, nic)
		nicIds = append(nicIds, fmt.Sprintf("%s.%s.id", output.GetResourceName(nic), nic.ResourceId))
	}

	// TODO change this to multy nsg_nic_attachment resource and use aws_network_interface_sg_attachment
	if len(r.NetworkSecurityGroups) != 0 {
		for _, nsg := range r.NetworkSecurityGroups {
			for _, nicId := range nicIds {
				nsgId, err := resources.GetMainOutputId(AzureNetworkSecurityGroup{nsg})
				if err != nil {
					return nil, err
				}
				azResources = append(
					azResources, network_security_group.AzureNetworkInterfaceSecurityGroupAssociation{
						AzResource: &common.AzResource{
							TerraformResource: output.TerraformResource{
								ResourceId: r.ResourceId,
							},
						},
						NetworkInterfaceId:     nicId,
						NetworkSecurityGroupId: nsgId,
					},
				)
			}
		}
	}

	// if ssh key is specified, add admin_ssh param
	// ssh authentication will replace password authentication
	// if no ssh key is passed, password is required
	// random_password will be used
	var azureSshKey virtual_machine.AzureAdminSshKey
	var vmPassword string
	disablePassAuth := false
	if r.Args.PublicSshKey != "" {
		azureSshKey = virtual_machine.AzureAdminSshKey{
			Username:  "adminuser",
			PublicKey: r.Args.PublicSshKey,
		}
		disablePassAuth = true
	} else {
		randomPassword := terraform.RandomPassword{
			TerraformResource: &output.TerraformResource{
				ResourceId: r.ResourceId,
			},
			Length:  16,
			Special: true,
			Upper:   true,
			Lower:   true,
			Number:  true,
		}
		vmPassword = randomPassword.GetResult()
		azResources = append(azResources, randomPassword)
	}

	computerName := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(r.Args.Name, "")

	sourceImg, err := virtual_machine.GetLatestAzureSourceImageReference(r.Args.ImageReference)
	if err != nil {
		return nil, err
	}

	var vmSize string
	if r.Args.AzureOverride.GetSize() != "" {
		vmSize = r.Args.AzureOverride.GetSize()
	} else {
		vmSize = common.VMSIZE[r.Args.VmSize][r.GetCloud()]
	}

	var zone string
	if r.Args.AvailabilityZone != 0 {
		zone, err = common.GetAvailabilityZone(r.GetLocation(), int(r.Args.AvailabilityZone), r.GetCloud())
		if err != nil {
			return nil, err
		}
	}

	azResources = append(
		azResources, virtual_machine.AzureVirtualMachine{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				ResourceGroupName: rgName,
				Name:              r.Args.Name,
			},
			Location:            r.GetCloudSpecificLocation(),
			Size:                vmSize,
			NetworkInterfaceIds: nicIds,
			CustomData:          r.Args.UserDataBase64,
			OsDisk: virtual_machine.AzureOsDisk{
				Caching:            "None",
				StorageAccountType: "Standard_LRS",
			},
			AdminUsername:                 "adminuser",
			AdminPassword:                 vmPassword,
			AdminSshKey:                   azureSshKey,
			SourceImageReference:          sourceImg,
			DisablePasswordAuthentication: disablePassAuth,
			Identity:                      virtual_machine.AzureIdentity{Type: "SystemAssigned"},
			ComputerName:                  computerName,
			Zone:                          zone,
		},
	)

	return azResources, nil

}

func (r AzureVirtualMachine) GetMainResourceName() (string, error) {
	return virtual_machine.AzureResourceName, nil
}
