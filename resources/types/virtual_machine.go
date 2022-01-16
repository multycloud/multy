package types

import (
	"encoding/base64"
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
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

func (vm *VirtualMachine) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []interface{} {
	if vm.UserData != "" {
		vm.UserData = base64.StdEncoding.EncodeToString([]byte(vm.UserData))
	}

	var subnetId = vm.SubnetId

	var nicIds []string
	for _, id := range vm.NetworkInterfaceIds {
		if n, err := ctx.GetResource(id); err != nil {
			vm.LogFatal(vm.ResourceId, "network_interface_ids", err.Error())
		} else {
			nicIds = append(nicIds, n.Resource.(*NetworkInterface).GetId(cloud))
		}
	}

	var nsgIds []string
	for _, id := range vm.NetworkSecurityGroupIds {
		if n, err := ctx.GetResource(id); err != nil {
			vm.LogFatal(vm.ResourceId, "network_security_group_ids", err.Error())
		} else {
			nsgIds = append(nsgIds, n.Resource.(*NetworkSecurityGroup).GetId(cloud))
		}
	}

	if cloud == common.AWS {
		var awsResources []interface{}
		var ec2NicIds []virtual_machine.AwsEc2NetworkInterface
		for i, id := range nicIds {
			ec2NicIds = append(ec2NicIds, virtual_machine.AwsEc2NetworkInterface{
				NetworkInterfaceId: id,
				DeviceIndex:        i,
			})
		}

		ec2 := virtual_machine.AwsEC2{
			AwsResource: common.AwsResource{
				ResourceName: virtual_machine.AwsResourceName,
				ResourceId:   vm.GetTfResourceId(cloud),
				Tags:         map[string]string{"Name": vm.Name},
			},
			Ami:                      "ami-09d4a659cdd8677be", // eu-west-2 "ami-0fc15d50d39e4503c", // https://cloud-images.ubuntu.com/locator/ec2/
			InstanceType:             common.VMSIZE[common.MICRO][cloud],
			AssociatePublicIpAddress: vm.PublicIp,
			UserDataBase64:           vm.UserData,
			SubnetId:                 subnetId,
			NetworkInterfaces:        ec2NicIds,
			SecurityGroupIds:         nsgIds,
		}

		if len(ec2NicIds) != 0 {
			ec2.SubnetId = ""
			ec2.AssociatePublicIpAddress = false
		}

		// adding ssh key to vm requires aws key pair resource
		// key pair will be added and referenced via key_name parameter
		if vm.SshKeyFileName != "" {
			keyPair := virtual_machine.AwsKeyPair{
				AwsResource: common.AwsResource{
					ResourceName: virtual_machine.AwsKeyPairResourceName,
					ResourceId:   vm.GetTfResourceId(cloud),
					Tags:         map[string]string{"Name": vm.Name},
				},
				KeyName:   fmt.Sprintf("%s_multy", vm.ResourceId),
				PublicKey: fmt.Sprintf("file(\"%s\")", vm.SshKeyFileName),
			}
			ec2.KeyName = vm.GetAssociatedKeyPairName(cloud)
			awsResources = append(awsResources, keyPair)
		}

		awsResources = append(awsResources, ec2)
		return awsResources
	} else if cloud == common.AZURE {
		// TODO validate that NIC is on the same VNET
		var azResources []interface{}
		rgName := rg.GetResourceGroupName(vm.ResourceGroupId, cloud)

		if len(vm.NetworkInterfaceIds) == 0 {
			nic := network_interface.AzureNetworkInterface{
				AzResource: common.AzResource{
					ResourceName:      network_interface.AzureResourceName,
					ResourceId:        vm.GetTfResourceId(cloud),
					ResourceGroupName: rgName,
					Name:              vm.Name,
					Location:          ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud),
				},
				IpConfigurations: []network_interface.AzureIpConfiguration{{
					Name:                       "internal", // this name shouldn't be vm.name
					PrivateIpAddressAllocation: "Dynamic",
					SubnetId:                   subnetId,
					Primary:                    true,
				}},
			}
			azResources = append(azResources, &nic)

			if vm.PublicIp {
				pIp := public_ip.AzurePublicIp{
					AzResource: common.AzResource{
						ResourceName:      public_ip.AzureResourceName,
						ResourceId:        vm.GetTfResourceId(cloud),
						ResourceGroupName: rgName,
						Name:              vm.Name,
						Location:          ctx.GetLocationFromCommonParams(vm.CommonResourceParams, cloud),
					},
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
			nicIds = append(nicIds, nic.GetId(cloud))
		}

		// TODO change this to multy nsg_nic_attachment resource and use aws_network_interface_sg_attachment
		if len(vm.NetworkSecurityGroupIds) != 0 {
			for _, nsgId := range nsgIds {
				for _, nicId := range nicIds {
					azResources = append(azResources, network_security_group.AzureNetworkInterfaceSecurityGroupAssociation{
						ResourceName:           network_security_group.AzureNicNsgAssociation,
						ResourceId:             vm.GetTfResourceId(cloud),
						NetworkInterfaceId:     nicId,
						NetworkSecurityGroupId: nsgId,
					})
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
				ResourceName: terraform.TerraformResourceName,
				ResourceId:   vm.GetTfResourceId(cloud),
				Length:       16,
				Special:      true,
				Upper:        true,
				Lower:        true,
				Number:       true,
			}
			vmPassword = randomPassword.GetResult()
			azResources = append(azResources, randomPassword)
		}

		azResources = append(azResources, virtual_machine.AzureVirtualMachine{
			AzResource: common.AzResource{
				ResourceName:      virtual_machine.AzureResourceName,
				ResourceId:        vm.GetTfResourceId(cloud),
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
		})

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

func (vm *VirtualMachine) Validate(ctx resources.MultyContext) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.Size valid { return false }
	if vm.OperatingSystem != "linux" { // max len?
		vm.LogFatal(vm.ResourceId, "os", "invalid operating system")
	}
	if vm.PublicIp && len(vm.NetworkInterfaceIds) != 0 {
		vm.LogFatal(vm.ResourceId, "public_ip", "public ip can't be set with network interface ids")
	}
	if vm.PublicIp && vm.PublicIpId != "" {
		vm.LogFatal(vm.ResourceId, "public_ip", "conflict between public_ip and public_ip_id")
	}
	if common.ValidateVmSize(vm.Size) {
		vm.LogFatal(vm.ResourceId, "size", fmt.Sprintf("\"%s\" is not []", vm.Size))
	}
	return
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
