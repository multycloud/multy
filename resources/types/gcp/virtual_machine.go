package gcp_resources

import (
	"encoding/base64"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_machine"
	"github.com/multycloud/multy/resources/types"
)

const (
	defaultSshUser = "adminuser"
)

type GcpVirtualMachine struct {
	*types.VirtualMachine
}

func InitVirtualMachine(vn *types.VirtualMachine) resources.ResourceTranslator[*resourcespb.VirtualMachineResource] {
	return GcpVirtualMachine{vn}
}

func (r GcpVirtualMachine) FromState(state *output.TfState) (*resourcespb.VirtualMachineResource, error) {
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

	vm, err := output.GetParsedById[virtual_machine.GoogleComputeInstance](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	if r.Args.GeneratePublicIp {
		out.PublicIp = vm.NetworkInterface[0].AccessConfig[0].NatIp
	}

	sa, err := output.GetParsedById[iam.GoogleServiceAccount](state, r.ResourceId)
	if err != nil {
		return nil, err
	}
	out.IdentityId = sa.Email

	out.GcpOutputs = &resourcespb.VirtualMachineGcpOutputs{
		ComputeInstanceId:   vm.SelfLink,
		ServiceAccountEmail: sa.Email,
	}

	return out, nil
}

func (r GcpVirtualMachine) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	az := r.Args.AvailabilityZone
	if az == 0 {
		az = 1
	}

	zone, err := common.GetAvailabilityZone(r.GetLocation(), int(az), r.GetCloud())
	if err != nil {
		return nil, err
	}
	var tags []string
	if len(r.NetworkSecurityGroups) > 0 {
		for _, nsg := range r.NetworkSecurityGroups {
			tags = append(tags, GcpNetworkSecurityGroup{nsg}.getNsgTag()...)
		}
	} else {
		// default network sg
		tags = append(tags, GcpVirtualNetwork{r.Subnet.VirtualNetwork}.getVnTag())
	}
	tags = append(tags, GcpSubnet{r.Subnet}.getNetworkTags()...)
	var size string
	if r.Args.GetGcpOverride().GetMachineType() != "" {
		size = r.Args.GetGcpOverride().GetMachineType()
	} else {
		size, err = common.GetVmSize(r.Args.VmSize, r.GetCloud())
		if err != nil {
			return nil, err
		}
	}
	image, err := virtual_machine.GetLatestGcpImage(r.Args.ImageReference)
	if err != nil {
		return nil, err
	}

	userData, err := base64.StdEncoding.DecodeString(r.Args.UserDataBase64)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	if len(userData) > 0 {
		m["user-data"] = string(userData)
	}

	if len(r.Args.PublicSshKey) > 0 {
		m["ssh-keys"] = fmt.Sprintf("%s:%s", defaultSshUser, r.Args.PublicSshKey)
	}

	networkInterface := virtual_machine.GoogleNetworkInterface{
		Subnetwork: fmt.Sprintf("%s.%s.self_link", output.GetResourceName(subnet.GoogleComputeSubnetwork{}), r.Subnet.ResourceId),
	}

	if r.Args.GeneratePublicIp {
		networkInterface.AccessConfig = []virtual_machine.GoogleNetworkInterfaceAccessConfig{{NetworkTier: "STANDARD"}}
	} else if r.PublicIp != nil {
		networkInterface.AccessConfig = []virtual_machine.GoogleNetworkInterfaceAccessConfig{
			{
				NetworkTier: "STANDARD",
				NatIp:       GcpPublicIp{r.PublicIp}.GetAddress(),
			},
		}
	}

	serviceAccountId := r.getServiceAccountId()
	serviceAccount := &iam.GoogleServiceAccount{
		GcpResource: common.NewGcpResource(r.ResourceId, "", r.Args.GetGcpOverride().GetProject()),
		AccountId:   serviceAccountId,
		DisplayName: fmt.Sprintf("Service Account for VM %s", r.Args.Name),
	}

	vm := &virtual_machine.GoogleComputeInstance{
		GcpResource: common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
		MachineType: size,
		Zone:        zone,
		Tags:        tags,
		BootDisk: virtual_machine.GoogleBootDisk{
			InitializeParams: virtual_machine.GoogleBootDiskInitializeParams{
				Image: image,
			},
		},
		Metadata:         m,
		NetworkInterface: []virtual_machine.GoogleNetworkInterface{networkInterface},
		ServiceAccount: virtual_machine.GoogleComputeInstanceServiceAccount{
			Email:  fmt.Sprintf("%s.%s.email", output.GetResourceName(iam.GoogleServiceAccount{}), r.ResourceId),
			Scopes: []string{"cloud-platform"},
		},
	}
	return []output.TfBlock{serviceAccount, vm}, nil
}

func (r GcpVirtualMachine) GetMainResourceName() (string, error) {
	return output.GetResourceName(virtual_machine.GoogleComputeInstance{}), nil
}

func (r GcpVirtualMachine) getServiceAccountId() string {
	return common.UniqueId(fmt.Sprintf("%s-%s", r.Args.Name, r.ResourceId), "-sa-", common.LowercaseAlphanumericAndDashFormatFunc)
}
