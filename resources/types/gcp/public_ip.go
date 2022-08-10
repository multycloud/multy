package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/types"
)

type GcpPublicIp struct {
	*types.PublicIp
}

func InitPublicIp(vn *types.PublicIp) resources.ResourceTranslator[*resourcespb.PublicIpResource] {
	return GcpPublicIp{vn}
}

func (r GcpPublicIp) FromState(state *output.TfState) (*resourcespb.PublicIpResource, error) {
	out := &resourcespb.PublicIpResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Ip:          "dryrun",
		GcpOverride: r.Args.GcpOverride,
	}
	if flags.DryRun {
		return out, nil
	}
	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[public_ip.GoogleComputeAddress](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.Ip = stateResource.Address
		out.GcpOutputs = &resourcespb.PublicIpGcpOutputs{
			ComputeAddressId: stateResource.SelfLink,
		}
	} else {
		statuses["gcp_compute_address"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r GcpPublicIp) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		public_ip.GoogleComputeAddress{
			GcpResource: common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
			NetworkTier: "STANDARD",
		},
	}, nil
}

func (r GcpPublicIp) GetMainResourceName() (string, error) {
	return output.GetResourceName(public_ip.GoogleComputeAddress{}), nil
}

func (r GcpPublicIp) GetAddress() string {
	return fmt.Sprintf("%s.%s.address", output.GetResourceName(public_ip.GoogleComputeAddress{}), r.ResourceId)
}
