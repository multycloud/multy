package types

import (
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/zclconf/go-cty/cty"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/network_interface"
	"multy-go/resources/output/network_security_group"
	"multy-go/resources/output/public_ip"
	"multy-go/resources/output/terraform"
	"multy-go/resources/output/virtual_machine"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
)

/*
Notes:
AWS: Can pass NICs and which overrides public_ip and subnet_id
Azure: To have a private IP by default, if no NIC is passed, one will be created.
       For PublicIp to be auto_assigned, public_ip is created an attached to default NIC
 	   NSG_NIC association
*/

type VirtualMachine struct {
	*resources.CommonResourceParams
	Name                    string   `hcl:"name"`
	OperatingSystem         string   `hcl:"os"`
	NetworkInterfaceIds     []string `hcl:"network_interface_ids,optional"`
	NetworkSecurityGroupIds []string `hcl:"network_security_group_ids,optional"`
	Size                    string   `hcl:"size"`
	UserData                string   `hcl:"user_data,optional"`
	SubnetId                string   `hcl:"subnet_id"`
	SshKeyFileName          string   `hcl:"ssh_key_file_path,optional"`
	PublicIpId              string   `hcl:"public_ip_id,optional"`
	// PublicIp auto-generate public IP
	PublicIp bool `hcl:"public_ip,optional"`
}

func (vm *VirtualMachine) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if vm.UserData != "" {
		vm.UserData = fmt.Sprintf("%s(\"%s\")", "base64encode", []byte(hclencoder.EscapeString(vm.UserData)))
	}

	var subnetId = vm.SubnetId

	if cloud == common.AWS {
		var awsResources []output.TfBlock
		var ec2NicIds []virtual_machine.AwsEc2NetworkInterface
		for i, id := range vm.NetworkInterfaceIds {
			ec2NicIds = append(
				ec2NicIds, virtual_machine.AwsEc2NetworkInterface{
					NetworkInterfaceId: id,
					DeviceIndex:        i,
				},
			)
		}

		ec2 := virtual_machine.AwsEC2{
			AwsResource:              common.NewAwsResource(vm.GetTfResourceId(cloud), vm.Name),
			Ami:                      common.AMIMAP[ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud)],
			InstanceType:             common.VMSIZE[common.MICRO][cloud],
			AssociatePublicIpAddress: vm.PublicIp,
			UserDataBase64:           vm.UserData,
			SubnetId:                 subnetId,
			NetworkInterfaces:        ec2NicIds,
			SecurityGroupIds:         vm.NetworkSecurityGroupIds,
		}

		if len(ec2NicIds) != 0 {
			ec2.SubnetId = ""
			ec2.AssociatePublicIpAddress = false
		}

		// adding ssh key to vm requires aws key pair resource
		// key pair will be added and referenced via key_name parameter
		if vm.SshKeyFileName != "" {
			keyPair := virtual_machine.AwsKeyPair{
				AwsResource: common.NewAwsResource(vm.GetTfResourceId(cloud), vm.Name),
				KeyName:     fmt.Sprintf("%s_multy", vm.ResourceId),
				PublicKey:   fmt.Sprintf("file(\"%s\")", vm.SshKeyFileName),
			}
			ec2.KeyName = vm.GetAssociatedKeyPairName(cloud)
			awsResources = append(awsResources, keyPair)
		}

		awsResources = append(awsResources, ec2)
		return awsResources
	} else if cloud == common.AZURE {
		// TODO validate that NIC is on the same VNET
		var azResources []output.TfBlock
		rgName := rg.GetResourceGroupName(vm.ResourceGroupId, cloud)
		nicIds := vm.NetworkInterfaceIds

		if len(vm.NetworkInterfaceIds) == 0 {
			nic := network_interface.AzureNetworkInterface{
				AzResource: common.NewAzResource(
					vm.GetTfResourceId(cloud), vm.Name, rgName,
					ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud),
				),
				IpConfigurations: []network_interface.AzureIpConfiguration{{
					Name:                       "internal", // this name shouldn't be vm.name
					PrivateIpAddressAllocation: "Dynamic",
					SubnetId:                   subnetId,
					Primary:                    true,
				}},
			}

			if vm.PublicIp {
				pIp := public_ip.AzurePublicIp{
					AzResource: common.NewAzResource(
						vm.GetTfResourceId(cloud), vm.Name, rgName,
						ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud),
					),
					AllocationMethod: "Static",
				}
				nic.IpConfigurations = []network_interface.AzureIpConfiguration{{
					Name:                       "external", // this name shouldn't be vm.name
					PrivateIpAddressAllocation: "Dynamic",
					SubnetId:                   subnetId,
					PublicIpAddressId:          pIp.GetId(cloud),
					Primary:                    true,
				}}
				azResources = append(azResources, &pIp)
			}
			azResources = append(azResources, nic)
			nicIds = append(nicIds, nic.GetId(cloud))
		}

		// TODO change this to multy nsg_nic_attachment resource and use aws_network_interface_sg_attachment
		if len(vm.NetworkSecurityGroupIds) != 0 {
			for _, nsgId := range vm.NetworkSecurityGroupIds {
				for _, nicId := range nicIds {
					azResources = append(
						azResources, network_security_group.AzureNetworkInterfaceSecurityGroupAssociation{
							AzResource: &common.AzResource{
								TerraformResource: output.TerraformResource{
									ResourceId: vm.GetTfResourceId(cloud),
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
		if vm.SshKeyFileName != "" {
			azureSshKey = virtual_machine.AzureAdminSshKey{
				Username:  "adminuser",
				PublicKey: fmt.Sprintf("file(\"%s\")", vm.SshKeyFileName),
			}
			disablePassAuth = true
		} else {
			randomPassword := terraform.RandomPassword{
				TerraformResource: &output.TerraformResource{
					ResourceId: vm.GetTfResourceId(cloud),
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

		azResources = append(
			azResources, virtual_machine.AzureVirtualMachine{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: vm.GetTfResourceId(cloud)},
					ResourceGroupName: rgName,
					Name:              vm.Name,
				},
				Location:            ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud),
				Size:                common.VMSIZE[common.MICRO][cloud],
				NetworkInterfaceIds: nicIds,
				CustomData:          vm.UserData,
				OsDisk: virtual_machine.AzureOsDisk{
					Caching:            "None",
					StorageAccountType: "Standard_LRS",
				},
				AdminUsername: "adminuser",
				AdminPassword: vmPassword,
				AdminSshKey:   azureSshKey,
				SourceImageReference: virtual_machine.AzureSourceImageReference{
					Publisher: "OpenLogic",
					Offer:     "CentOS",
					Sku:       "7_9-gen2",
					Version:   "latest",
				},
				DisablePasswordAuthentication: disablePassAuth,
			},
		)

		return azResources
	}

	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (vm *VirtualMachine) GetAssociatedKeyPairName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return fmt.Sprintf("%s.%s.key_name", virtual_machine.AwsKeyPairResourceName, vm.GetTfResourceId(common.AWS))
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (vm *VirtualMachine) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.Size valid { return false }
	if vm.OperatingSystem != "linux" { // max len?
		vm.NewError("os", "invalid operating system")
	}
	if vm.PublicIp && len(vm.NetworkInterfaceIds) != 0 {
		vm.NewError("public_ip", "public ip can't be set with network interface ids")
	}
	if vm.PublicIp && vm.PublicIpId != "" {
		vm.NewError("public_ip", "conflict between public_ip and public_ip_id")
	}
	if common.ValidateVmSize(vm.Size) {
		vm.NewError("size", fmt.Sprintf("\"%s\" is not []", vm.Size))
	}
	return errs
}

func (vm *VirtualMachine) GetMainResourceName(cloud common.CloudProvider) string {
	switch cloud {
	case common.AWS:
		return virtual_machine.AwsResourceName
	case common.AZURE:
		return virtual_machine.AzureResourceName
	default:
		validate.LogInternalError("unknown cloud %s", cloud)
	}
	return ""
}

func (vm *VirtualMachine) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	if vm.PublicIp {
		if cloud == common.AWS {
			return map[string]cty.Value{
				"public_ip": cty.StringVal(
					fmt.Sprintf(
						"${%s.%s.public_ip}", common.GetResourceName(virtual_machine.AwsEC2{}),
						vm.GetTfResourceId(cloud),
					),
				),
			}
		} else if cloud == common.AZURE {
			return map[string]cty.Value{
				"public_ip": cty.StringVal(
					fmt.Sprintf(
						"${%s.%s.ip_address}", common.GetResourceName(public_ip.AzurePublicIp{}),
						vm.GetTfResourceId(cloud),
					),
				),
			}
		}
	}
	return map[string]cty.Value{}
}
