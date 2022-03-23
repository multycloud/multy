package types

import (
	"encoding/json"
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/output/terraform"
	"github.com/multycloud/multy/resources/output/virtual_machine"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
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
	Name                    string                  `hcl:"name"`
	OperatingSystem         string                  `hcl:"os"`
	NetworkInterfaceIds     []*NetworkInterface     `mhcl:"ref=network_interface_ids,optional"`
	NetworkSecurityGroupIds []*NetworkSecurityGroup `mhcl:"ref=network_security_group_ids,optional"`
	Size                    string                  `hcl:"size"`
	UserData                string                  `hcl:"user_data,optional"`
	SubnetId                *Subnet                 `mhcl:"ref=subnet_id"`
	SshKeyFileName          string                  `hcl:"ssh_key_file_path,optional"`
	SshKey                  string                  `hcl:"ssh_key,optional"`
	PublicIpId              *PublicIp               `mhcl:"ref=public_ip_id,optional"`
	// PublicIp auto-generate public IP
	PublicIp bool `hcl:"public_ip,optional"`
}

type AwsCallerIdentityData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_caller_identity"`
}

type AwsRegionData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_region"`
}

func (vm *VirtualMachine) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if vm.UserData != "" {
		vm.UserData = fmt.Sprintf("%s(\"%s\")", "base64encode", []byte(hclencoder.EscapeString(vm.UserData)))
	}

	var subnetId = resources.GetMainOutputId(vm.SubnetId, cloud)

	sshKeyData := vm.SshKey
	if vm.SshKey == "" && vm.SshKeyFileName != "" {
		sshKeyData = fmt.Sprintf("${file(\"%s\")}", vm.SshKeyFileName)
	}
	if cloud == common.AWS {
		var awsResources []output.TfBlock
		var ec2NicIds []virtual_machine.AwsEc2NetworkInterface
		for i, ni := range vm.NetworkInterfaceIds {
			ec2NicIds = append(
				ec2NicIds, virtual_machine.AwsEc2NetworkInterface{
					NetworkInterfaceId: resources.GetMainOutputId(ni, cloud),
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
			SecurityGroupIds: util.MapSliceValues(vm.NetworkSecurityGroupIds, func(v *NetworkSecurityGroup) string {
				return resources.GetMainOutputId(v, cloud)
			}),
		}

		if len(ec2NicIds) != 0 {
			ec2.SubnetId = ""
			ec2.AssociatePublicIpAddress = false
		}

		iamRole := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(vm.GetTfResourceId(cloud), vm.Name),
			Name:             vm.getAwsIamRoleName(),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
		}

		if vault := getVaultAssociatedIdentity(ctx, iamRole.GetId()); vault != nil {
			awsResources = append(awsResources,
				AwsCallerIdentityData{TerraformDataSource: &output.TerraformDataSource{ResourceId: vm.GetTfResourceId(cloud)}},
				AwsRegionData{TerraformDataSource: &output.TerraformDataSource{ResourceId: vm.GetTfResourceId(cloud)}})

			policy, err := json.Marshal(iam.AwsIamPolicy{
				Statement: []iam.AwsIamPolicyStatement{{
					Action:   []string{"ssm:GetParameter*"},
					Effect:   "Allow",
					Resource: fmt.Sprintf("arn:aws:ssm:${data.aws_region.%s.name}:${data.aws_caller_identity.%s.account_id}:parameter/%s/*", vm.GetTfResourceId(cloud), vm.GetTfResourceId(cloud), vault.Name),
				}, {
					Action:   []string{"ssm:DescribeParameters"},
					Effect:   "Allow",
					Resource: "*",
				}},
				Version: "2012-10-17",
			})

			if err != nil {
				validate.LogInternalError("unable to encode aws policy: %s", err.Error())
			}

			iamRole.InlinePolicy = iam.AwsIamRoleInlinePolicy{
				Name:   "vault_policy",
				Policy: string(policy),
			}
		}

		iamInstanceProfile := iam.AwsIamInstanceProfile{
			AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: vm.GetTfResourceId(cloud)}},
			Name:        vm.getAwsIamRoleName(),
			Role: fmt.Sprintf(
				"%s.%s.name", output.GetResourceName(iam.AwsIamRole{}), iamRole.ResourceId,
			),
		}

		ec2.IamInstanceProfile = iamInstanceProfile.GetId()

		awsResources = append(awsResources,
			iamInstanceProfile,
			iamRole,
		)

		// this gives permission to write cloudwatch logs

		// adding ssh key to vm requires aws key pair resource
		// key pair will be added and referenced via key_name parameter
		if sshKeyData != "" {
			keyPair := virtual_machine.AwsKeyPair{
				AwsResource: common.NewAwsResource(vm.GetTfResourceId(cloud), vm.Name),
				KeyName:     fmt.Sprintf("%s_multy", vm.ResourceId),
				PublicKey:   sshKeyData,
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
		nicIds := util.MapSliceValues(vm.NetworkInterfaceIds, func(v *NetworkInterface) string {
			return resources.GetMainOutputId(v, cloud)
		})

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
			nicIds = append(nicIds, fmt.Sprintf("${%s.%s.id}", output.GetResourceName(nic), nic.ResourceId))
		}

		// TODO change this to multy nsg_nic_attachment resource and use aws_network_interface_sg_attachment
		if len(vm.NetworkSecurityGroupIds) != 0 {
			for _, nsg := range vm.NetworkSecurityGroupIds {
				for _, nicId := range nicIds {
					azResources = append(
						azResources, network_security_group.AzureNetworkInterfaceSecurityGroupAssociation{
							AzResource: &common.AzResource{
								TerraformResource: output.TerraformResource{
									ResourceId: vm.GetTfResourceId(cloud),
								},
							},
							NetworkInterfaceId:     nicId,
							NetworkSecurityGroupId: resources.GetMainOutputId(nsg, cloud),
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
		if sshKeyData != "" {
			azureSshKey = virtual_machine.AzureAdminSshKey{
				Username:  "adminuser",
				PublicKey: sshKeyData,
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
				Identity:                      virtual_machine.AzureIdentity{Type: "SystemAssigned"},
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

// check if VM identity is associated with vault (vault_access_policy)
// get Vault name by vault_access_policy that uses VM identity
func getVaultAssociatedIdentity(ctx resources.MultyContext, identity string) *Vault {
	for _, resource := range resources.GetAllResources[*VaultAccessPolicy](ctx) {
		if identity == resource.Identity {
			return resource.Vault
		}
	}
	return nil
}

func (vm *VirtualMachine) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.Size valid { return false }
	if vm.OperatingSystem != "linux" { // max len?
		errs = append(errs, vm.NewError("os", "invalid operating system"))
	}
	if vm.PublicIp && len(vm.NetworkInterfaceIds) != 0 {
		errs = append(errs, vm.NewError("public_ip", "public ip can't be set with network interface ids"))
	}
	if vm.PublicIp && vm.PublicIpId != nil {
		errs = append(errs, vm.NewError("public_ip", "conflict between public_ip and public_ip_id"))
	}
	if common.ValidateVmSize(vm.Size) {
		errs = append(errs, vm.NewError("size", fmt.Sprintf("\"%s\" is not []", vm.Size)))
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
						"${%s.%s.public_ip}", output.GetResourceName(virtual_machine.AwsEC2{}),
						vm.GetTfResourceId(cloud),
					),
				),
				"identity": cty.StringVal(
					fmt.Sprintf(
						"${%s.%s.id}", output.GetResourceName(iam.AwsIamRole{}),
						vm.GetTfResourceId(cloud),
					),
				),
			}
		} else if cloud == common.AZURE {
			return map[string]cty.Value{
				"public_ip": cty.StringVal(
					fmt.Sprintf(
						"${%s.%s.ip_address}", output.GetResourceName(public_ip.AzurePublicIp{}),
						vm.GetTfResourceId(cloud),
					),
				),
				"identity": cty.StringVal(
					fmt.Sprintf(
						"${%s.%s.identity[0].principal_id}", output.GetResourceName(virtual_machine.AzureVirtualMachine{}),
						vm.GetTfResourceId(cloud),
					),
				),
			}
		}
	}
	return map[string]cty.Value{}
}

func (vm *VirtualMachine) getAwsIamRoleName() string {
	return fmt.Sprintf("iam_for_vm_%s", vm.ResourceId)
}
