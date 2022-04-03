package types

import (
	"encoding/json"
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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
)

/*
Notes:
AWS: Can pass NICs and which overrides public_ip and subnet_id
Azure: To have a private IP by default, if no NIC is passed, one will be created.
       For PublicIp to be auto_assigned, public_ip is created an attached to default NIC
 	   NSG_NIC association
*/

type VirtualMachine struct {
	resources.ResourceWithId[*resourcespb.VirtualMachineArgs]

	NetworkInterface      []*NetworkInterface
	NetworkSecurityGroups []*NetworkSecurityGroup
	Subnet                *Subnet
	PublicIp              *PublicIp
}

func NewVirtualMachine(resourceId string, args *resourcespb.VirtualMachineArgs, others resources.Resources) (*VirtualMachine, error) {
	networkInterfaces, err := util.MapSliceValuesErr(args.NetworkInterfaceIds, func(id string) (*NetworkInterface, error) {
		return Get[*NetworkInterface](others, id)
	})
	if err != nil {
		return nil, err
	}
	networkSecurityGroups, err := util.MapSliceValuesErr(args.NetworkSecurityGroupIds, func(id string) (*NetworkSecurityGroup, error) {
		return Get[*NetworkSecurityGroup](others, id)
	})
	if err != nil {
		return nil, err
	}
	subnet, err := Get[*Subnet](others, args.SubnetId)
	if err != nil {
		return nil, err
	}
	publicIp, _, err := GetOptional[*PublicIp](others, args.PublicIpId)
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{
		ResourceWithId: resources.ResourceWithId[*resourcespb.VirtualMachineArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		NetworkInterface:      networkInterfaces,
		NetworkSecurityGroups: networkSecurityGroups,
		Subnet:                subnet,
		PublicIp:              publicIp,
	}, nil
}

type AwsCallerIdentityData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_caller_identity"`
}

type AwsRegionData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_region"`
}

func (r *VirtualMachine) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	if r.Args.UserData != "" {
		r.Args.UserData = fmt.Sprintf("%s(\"%s\")", "base64encode", []byte(hclencoder.EscapeString(r.Args.UserData)))
	}

	subnetId, err := resources.GetMainOutputId(r.Subnet)
	if err != nil {
		return nil, err
	}

	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var awsResources []output.TfBlock
		var ec2NicIds []virtual_machine.AwsEc2NetworkInterface
		for i, ni := range r.NetworkInterface {
			niId, err := resources.GetMainOutputId(ni)
			if err != nil {
				return nil, err
			}
			ec2NicIds = append(
				ec2NicIds, virtual_machine.AwsEc2NetworkInterface{
					NetworkInterfaceId: niId,
					DeviceIndex:        i,
				},
			)
		}

		nsgIds, err := util.MapSliceValuesErr(r.NetworkSecurityGroups, func(v *NetworkSecurityGroup) (string, error) {
			return resources.GetMainOutputId(v)
		})
		if err != nil {
			return nil, err
		}

		ec2 := virtual_machine.AwsEC2{
			AwsResource:              common.NewAwsResource(r.ResourceId, r.Args.Name),
			Ami:                      common.AMIMAP[r.GetCloudSpecificLocation()],
			InstanceType:             common.VMSIZE[r.Args.VmSize][r.GetCloud()],
			AssociatePublicIpAddress: r.Args.GeneratePublicIp,
			UserDataBase64:           r.Args.UserData,
			SubnetId:                 subnetId,
			NetworkInterfaces:        ec2NicIds,
			SecurityGroupIds:         nsgIds,
		}

		if len(ec2NicIds) != 0 {
			ec2.SubnetId = ""
			ec2.AssociatePublicIpAddress = false
		}

		iamRole := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.ResourceId, r.Args.Name),
			Name:             r.getAwsIamRoleName(),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
		}

		if vault := getVaultAssociatedIdentity(ctx, iamRole.GetId()); vault != nil {
			awsResources = append(awsResources,
				AwsCallerIdentityData{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}},
				AwsRegionData{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}})

			policy, err := json.Marshal(iam.AwsIamPolicy{
				Statement: []iam.AwsIamPolicyStatement{{
					Action:   []string{"ssm:GetParameter*"},
					Effect:   "Allow",
					Resource: fmt.Sprintf("arn:aws:ssm:${data.aws_region.%s.name}:${data.aws_caller_identity.%s.account_id}:parameter/%s/*", r.ResourceId, r.ResourceId, vault.Args.Name),
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
			AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
			Name:        r.getAwsIamRoleName(),
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

		// adding ssh key to r requires aws key pair resource
		// key pair will be added and referenced via key_name parameter
		if r.Args.PublicSshKey != "" {
			keyPair := virtual_machine.AwsKeyPair{
				AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
				KeyName:     fmt.Sprintf("%s_multy", r.ResourceId),
				PublicKey:   r.Args.PublicSshKey,
			}
			ec2.KeyName, err = r.GetAssociatedKeyPairName()
			if err != nil {
				return nil, err
			}
			awsResources = append(awsResources, keyPair)
		}

		awsResources = append(awsResources, ec2)
		return awsResources, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		// TODO validate that NIC is on the same VNET
		var azResources []output.TfBlock
		rgName := rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
		nicIds, err := util.MapSliceValuesErr(r.NetworkInterface, func(v *NetworkInterface) (string, error) {
			return resources.GetMainOutputId(v)
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
					AllocationMethod: "Static",
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
			nicIds = append(nicIds, fmt.Sprintf("${%s.%s.id}", output.GetResourceName(nic), nic.ResourceId))
		}

		// TODO change this to multy nsg_nic_attachment resource and use aws_network_interface_sg_attachment
		if len(r.NetworkSecurityGroups) != 0 {
			for _, nsg := range r.NetworkSecurityGroups {
				for _, nicId := range nicIds {
					nsgId, err := resources.GetMainOutputId(nsg)
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

		azResources = append(
			azResources, virtual_machine.AzureVirtualMachine{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
					ResourceGroupName: rgName,
					Name:              r.Args.Name,
				},
				Location:            r.GetCloudSpecificLocation(),
				Size:                common.VMSIZE[r.Args.VmSize][r.GetCloud()],
				NetworkInterfaceIds: nicIds,
				CustomData:          r.Args.UserData,
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

		return azResources, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *VirtualMachine) GetAssociatedKeyPairName() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return fmt.Sprintf("%s.%s.key_name", virtual_machine.AwsKeyPairResourceName, r.ResourceId), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

// check if VM identity is associated with vault (vault_access_policy)
// get Vault name by vault_access_policy that uses VM identity
func getVaultAssociatedIdentity(ctx resources.MultyContext, identity string) *Vault {
	for _, resource := range resources.GetAllResources[*VaultAccessPolicy](ctx) {
		if identity == resource.Args.Identity {
			return resource.Vault
		}
	}
	return nil
}

func (r *VirtualMachine) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.Size valid { return false }
	if r.Args.GeneratePublicIp && len(r.NetworkInterface) != 0 {
		errs = append(errs, r.NewValidationError("generate public ip can't be set with network interface ids", "generate_public_ip"))
	}
	if r.Args.GeneratePublicIp && r.PublicIp != nil {
		errs = append(errs, r.NewValidationError("conflict between generate_public_ip and public_ip_id", "generate_public_ip"))
	}
	return errs
}

func (r *VirtualMachine) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return virtual_machine.AwsResourceName, nil
	case common.AZURE:
		return virtual_machine.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}

func (r *VirtualMachine) getAwsIamRoleName() string {
	return fmt.Sprintf("iam_for_vm_%s", r.ResourceId)
}
