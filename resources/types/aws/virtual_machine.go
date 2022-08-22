package aws_resources

import (
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

func (r AwsVirtualMachine) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.VirtualMachineResource, error) {
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
	}

	if flags.DryRun {
		return out, nil
	}

	out.AwsOutputs = &resourcespb.VirtualMachineAwsOutputs{}
	statuses := map[string]commonpb.ResourceStatus_Status{}
	keyName := ""

	if vmResource, exists, err := output.MaybeGetParsedById[virtual_machine.AwsEC2](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.Ec2InstanceId = vmResource.ResourceId

		out.PublicIp = vmResource.PublicIp
		out.Name = vmResource.Tags["Name"]
		if len(r.Args.GetAwsOverride().GetInstanceType()) > 0 {
			out.AwsOverride.InstanceType = vmResource.InstanceType
		} else {
			out.VmSize = common.ParseVmSize(vmResource.InstanceType, common.AWS)
		}
		out.UserDataBase64 = vmResource.UserDataBase64
		out.GeneratePublicIp = vmResource.AssociatePublicIpAddress
		keyName = vmResource.KeyName
		if vmResource.KeyName == "" {
			out.PublicSshKey = ""
		}
		output.AddToStatuses(statuses, "aws_ec2", output.MaybeGetPlannedChageById[virtual_machine.AwsEC2](plan, r.ResourceId))
	} else {
		statuses["aws_ec2"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if iamRoleResource, exists, err := output.MaybeGetParsedById[iam.AwsIamRole](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.IdentityId = iamRoleResource.Id
		out.AwsOutputs.IamRoleArn = iamRoleResource.Arn
		output.AddToStatuses(statuses, "aws_iam_role", output.MaybeGetPlannedChageById[iam.AwsIamRole](plan, r.ResourceId))
	} else {
		statuses["aws_iam_role"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if iamInstanceProfileResource, exists, err := output.MaybeGetParsedById[iam.AwsIamInstanceProfile](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.IamInstanceProfileArn = iamInstanceProfileResource.Arn
		output.AddToStatuses(statuses, "aws_iam_instance_profile", output.MaybeGetPlannedChageById[iam.AwsIamInstanceProfile](plan, r.ResourceId))
	} else {
		statuses["aws_iam_instance_profile"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[virtual_machine.AwsKeyPair](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.AwsOutputs.KeyPairArn = stateResource.Arn
		if keyName == stateResource.KeyName {
			out.PublicSshKey = stateResource.PublicKey
		}
		output.AddToStatuses(statuses, "aws_key_pair", output.MaybeGetPlannedChageById[virtual_machine.AwsKeyPair](plan, r.ResourceId))
	} else if r.Args.PublicSshKey != "" {
		statuses["aws_key_pair"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

type AwsCallerIdentityData struct {
	*output.TerraformDataSource `hcl:",squash" default:"name=aws_caller_identity"`
}

func (r AwsVirtualMachine) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	if r.Args.UserDataBase64 != "" {
		r.Args.UserDataBase64 = fmt.Sprintf(hclencoder.EscapeString(r.Args.UserDataBase64))
	}

	subnetId, err := AwsSubnet{r.Subnet}.GetSubnetId(r.Args.AvailabilityZone)
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
			KeyName:     common.UniqueId(fmt.Sprintf("%s-%s", r.Args.Name, r.ResourceId), "-key-", common.LowercaseAlphanumericAndDashFormatFunc),
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
func (r AwsVirtualMachine) GetIdentity() string {
	return r.GetAwsIdentity()
}
