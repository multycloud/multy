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
	ip := "dryrun"
	if !flags.DryRun {
		address, err := output.GetParsedById[public_ip.GoogleComputeAddress](state, r.ResourceId)
		if err != nil {
			return nil, err
		}
		ip = address.Address
	}

	return &resourcespb.PublicIpResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name: r.Args.Name,

		Ip: ip,
	}, nil
}

func (r GcpPublicIp) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		public_ip.GoogleComputeAddress{
			GcpResource: common.NewGcpResource(r.ResourceId, r.Args.Name),
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
