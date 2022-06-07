package types

import (
	"encoding/json"
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/resources/output/network_security_group"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/output/terraform"
	"github.com/multycloud/multy/resources/output/virtual_machine"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"log"
	"regexp"
)

/*
Notes:
AWS: Can pass NICs and which overrides public_ip and subnet_id
Azure: To have a private IP by default, if no NIC is passed, one will be created.
       For PublicIp to be auto_assigned, public_ip is created an attached to default NIC
 	   NSG_NIC association
*/

var virtualMachineMetadata = resources.ResourceMetadata[*resourcespb.VirtualMachineArgs, *VirtualMachine, *resourcespb.VirtualMachineResource]{
	CreateFunc:        CreateVirtualMachine,
	UpdateFunc:        UpdateVirtualMachine,
	ReadFromStateFunc: VirtualMachineFromState,
	ExportFunc: func(r *VirtualMachine, _ *resources.Resources) (*resourcespb.VirtualMachineArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewVirtualMachine,
	AbbreviatedName: "vm",
}

type VirtualMachine struct {
	resources.ResourceWithId[*resourcespb.VirtualMachineArgs]

	NetworkInterface      []*NetworkInterface
	NetworkSecurityGroups []*NetworkSecurityGroup
	Subnet                *Subnet
	PublicIp              *PublicIp
}

func (r *VirtualMachine) GetMetadata() resources.ResourceMetadataInterface {
	return &virtualMachineMetadata
}

func CreateVirtualMachine(resourceId string, args *resourcespb.VirtualMachineArgs, others *resources.Resources) (*VirtualMachine, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return nil, err
		}
		rgId, err := NewRgFromParent("vm", subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return nil, err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewVirtualMachine(resourceId, args, others)
}

func UpdateVirtualMachine(resource *VirtualMachine, vn *resourcespb.VirtualMachineArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func VirtualMachineFromState(resource *VirtualMachine, state *output.TfState) (*resourcespb.VirtualMachineResource, error) {
	var err error
	var ip string
	identityId := "dryrun"
	if resource.Args.GeneratePublicIp {
		ip = "dryrun"
	}

	if !flags.DryRun {
		if resource.Args.GeneratePublicIp {
			ip, err = getPublicIp(resource.ResourceId, state, resource.Args.CommonParameters.CloudProvider)
			if err != nil {
				return nil, err
			} else {
				ip = ""
			}
		}
		identityId, err = getIdentityId(resource.ResourceId, state, resource.Args.CommonParameters.CloudProvider)
		if err != nil {
			return nil, err
		}
	}

	// TODO: handle default values on create
	if resource.Args.ImageReference == nil {
		resource.Args.ImageReference = &resourcespb.ImageReference{
			Os:      resourcespb.ImageReference_UBUNTU,
			Version: "16.04",
		}
	}

	return &resourcespb.VirtualMachineResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:                    resource.Args.Name,
		NetworkInterfaceIds:     resource.Args.NetworkInterfaceIds,
		NetworkSecurityGroupIds: resource.Args.NetworkSecurityGroupIds,
		VmSize:                  resource.Args.VmSize,
		UserDataBase64:          resource.Args.UserDataBase64,
		SubnetId:                resource.Args.SubnetId,
		PublicSshKey:            resource.Args.PublicSshKey,
		PublicIpId:              resource.Args.PublicIpId,
		GeneratePublicIp:        resource.Args.GeneratePublicIp,
		ImageReference:          resource.Args.ImageReference,
		AwsOverride:             resource.Args.AwsOverride,
		AzureOverride:           resource.Args.AzureOverride,

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

func NewVirtualMachine(resourceId string, args *resourcespb.VirtualMachineArgs, others *resources.Resources) (*VirtualMachine, error) {
	networkInterfaces, err := util.MapSliceValuesErr(args.NetworkInterfaceIds, func(id string) (*NetworkInterface, error) {
		return resources.Get[*NetworkInterface](resourceId, others, id)
	})
	if err != nil {
		return nil, err
	}
	networkSecurityGroups, err := util.MapSliceValuesErr(args.NetworkSecurityGroupIds, func(id string) (*NetworkSecurityGroup, error) {
		return resources.Get[*NetworkSecurityGroup](resourceId, others, id)
	})
	if err != nil {
		return nil, err
	}
	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return nil, err
	}
	publicIp, _, err := resources.GetOptional[*PublicIp](resourceId, others, args.PublicIpId)
	if err != nil {
		return nil, err
	}
	if args.GetImageReference() == nil {
		args.ImageReference = &resourcespb.ImageReference{
			Os:      resourcespb.ImageReference_UBUNTU,
			Version: "16.04",
		}
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
	if r.Args.UserDataBase64 != "" {
		r.Args.UserDataBase64 = fmt.Sprintf(hclencoder.EscapeString(r.Args.UserDataBase64))
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

		var vmSize string
		if r.Args.AwsOverride.GetInstanceType() != "" {
			vmSize = r.Args.AwsOverride.GetInstanceType()
		} else {
			vmSize = common.VMSIZE[r.Args.VmSize][r.GetCloud()]
		}

		ec2 := virtual_machine.AwsEC2{
			AwsResource:              common.NewAwsResource(r.ResourceId, r.Args.Name),
			InstanceType:             vmSize,
			AssociatePublicIpAddress: r.Args.GeneratePublicIp,
			UserDataBase64:           r.Args.UserDataBase64,
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
			Name:             r.GetAwsIdentity(),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
		}

		if vault := getVaultAssociatedIdentity(ctx, r.GetAwsIdentity()); vault != nil {
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
				return nil, fmt.Errorf("unable to encode aws policy: %s", err.Error())
			}

			iamRole.InlinePolicy = iam.AwsIamRoleInlinePolicy{
				Name: "vault_policy",
				// we need to have an expression here because we use template strings within the policy json
				Policy: fmt.Sprintf("\"%s\"", hclencoder.EscapeString(string(policy))),
			}
		}

		iamInstanceProfile := iam.AwsIamInstanceProfile{
			AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
			Name:        r.GetAwsIdentity(),
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

		awsAmi, err := virtual_machine.LatestAwsAmi(r.Args.ImageReference, r.ResourceId)
		if err != nil {
			return nil, err
		}

		ec2.Ami = fmt.Sprintf("%s.id", awsAmi.GetFullResourceRef())

		awsResources = append(awsResources, awsAmi, ec2)
		return awsResources, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		// TODO validate that NIC is on the same VNET
		var azResources []output.TfBlock
		rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
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
			nicIds = append(nicIds, fmt.Sprintf("%s.%s.id", output.GetResourceName(nic), nic.ResourceId))
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

		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		computerName := reg.ReplaceAllString(r.Args.Name, "")

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
		errs = append(errs, r.NewValidationError(fmt.Errorf("generate public ip can't be set with network interface ids"), "generate_public_ip"))
	}
	if r.Args.GeneratePublicIp && r.PublicIp != nil {
		errs = append(errs, r.NewValidationError(fmt.Errorf("conflict between generate_public_ip and public_ip_id"), "generate_public_ip"))
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

func (r *VirtualMachine) GetAwsIdentity() string {
	return getIdentity(r.ResourceId)
}

func getIdentity(resourceId string) string {
	return fmt.Sprintf("multy-vm-%s-role", resourceId)
}
