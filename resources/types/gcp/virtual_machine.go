package gcp_resources

import (
	"encoding/base64"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
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
	}

	vm, err := output.GetParsedById[virtual_machine.GoogleComputeInstance](state, r.ResourceId)
	if err != nil {
		return nil, err
	}

	if r.Args.GeneratePublicIp {
		out.PublicIp = vm.NetworkInterface[0].AccessConfig[0].NatIp
	}

	return out, nil
}

func (r GcpVirtualMachine) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	zone, err := common.GetAvailabilityZone(r.GetLocation(), int(r.Subnet.Args.AvailabilityZone), r.GetCloud())
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, nsg := range r.NetworkSecurityGroups {
		tags = append(tags, GcpNetworkSecurityGroup{nsg}.getNsgTag()...)
	}
	tags = append(tags, GcpSubnet{r.Subnet}.getNetworkTag())
	size, err := common.GetVmSize(r.Args.VmSize, r.GetCloud())
	if err != nil {
		return nil, err
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
	}
	return []output.TfBlock{vm}, nil
}

func (r GcpVirtualMachine) GetMainResourceName() (string, error) {
	return output.GetResourceName(virtual_machine.GoogleComputeInstance{}), nil
}
