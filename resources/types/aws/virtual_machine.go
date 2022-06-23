package aws_resources

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
	"github.com/multycloud/multy/resources/output/virtual_machine"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
)

type AwsVirtualMachine struct {
	*types.VirtualMachine
}

func InitVirtualMachine(vn *types.VirtualMachine) resources.ResourceTranslator[*resourcespb.VirtualMachineResource] {
	return AwsVirtualMachine{vn}
}

func (r AwsVirtualMachine) FromState(state *output.TfState) (*resourcespb.VirtualMachineResource, error) {
	var ip string
	identityId := "dryrun"
	if r.Args.GeneratePublicIp {
		ip = "dryrun"
	}

	if !flags.DryRun {
		if r.Args.GeneratePublicIp {
			vmResource, err := output.GetParsedById[virtual_machine.AwsEC2](state, r.ResourceId)
			if err != nil {
				return nil, err
			}
			ip = vmResource.PublicIp
		}

		iamRoleResource, err := output.GetParsedById[iam.AwsIamRole](state, r.ResourceId)
		if err != nil {
			return nil, err
		}
		identityId = iamRoleResource.Id
	}

	// TODO: handle default values on create
	if r.Args.ImageReference == nil {
		r.Args.ImageReference = &resourcespb.ImageReference{
			Os:      resourcespb.ImageReference_UBUNTU,
			Version: "16.04",
		}
	}

	return &resourcespb.VirtualMachineResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
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
		PublicIp:                ip,
		IdentityId:              identityId,
	}, nil
}

type AwsCallerIdentityData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_caller_identity"`
}

type AwsRegionData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_region"`
}

func (r AwsVirtualMachine) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	if r.Args.UserDataBase64 != "" {
		r.Args.UserDataBase64 = fmt.Sprintf(hclencoder.EscapeString(r.Args.UserDataBase64))
	}

	subnetId, err := resources.GetMainOutputId(AwsSubnet{r.Subnet})
	if err != nil {
		return nil, err
	}

	var awsResources []output.TfBlock
	var ec2NicIds []virtual_machine.AwsEc2NetworkInterface
	for i, ni := range r.NetworkInterface {
		niId, err := resources.GetMainOutputId(AwsNetworkInterface{ni})
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

	nsgIds, err := util.MapSliceValuesErr(r.NetworkSecurityGroups, func(v *types.NetworkSecurityGroup) (string, error) {
		return resources.GetMainOutputId(AwsNetworkSecurityGroup{v})
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
		Name:             r.GetIdentity(),
		AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
	}

	if vault := getVaultAssociatedIdentity(ctx, r.GetIdentity()); vault != nil {
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
			return nil, fmt.Errorf("unable to encode aws policy: %s", err)
		}

		iamRole.InlinePolicy = iam.AwsIamRoleInlinePolicy{
			Name: "vault_policy",
			// we need to have an expression here because we use template strings within the policy json
			Policy: fmt.Sprintf("\"%s\"", hclencoder.EscapeString(string(policy))),
		}
	}

	iamInstanceProfile := iam.AwsIamInstanceProfile{
		AwsResource: &common.AwsResource{TerraformResource: output.TerraformResource{ResourceId: r.ResourceId}},
		Name:        r.GetIdentity(),
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
		ec2.KeyName = fmt.Sprintf("%s.%s.key_name", virtual_machine.AwsKeyPairResourceName, r.ResourceId)
		awsResources = append(awsResources, keyPair)
	}

	awsAmi, err := virtual_machine.LatestAwsAmi(r.Args.ImageReference, r.ResourceId)
	if err != nil {
		return nil, err
	}

	ec2.Ami = fmt.Sprintf("%s.id", awsAmi.GetFullResourceRef())

	awsResources = append(awsResources, awsAmi, ec2)
	return awsResources, nil

}

func (r AwsVirtualMachine) GetMainResourceName() (string, error) {
	return virtual_machine.AwsResourceName, nil
}

func getVaultAssociatedIdentity(ctx resources.MultyContext, identity string) *types.Vault {
	for _, resource := range resources.GetAllResources[*types.VaultAccessPolicy](ctx) {
		if identity == resource.Args.Identity {
			return resource.Vault
		}
	}
	return nil
}

func (r AwsVirtualMachine) GetIdentity() string {
	return r.GetAwsIdentity()
}
