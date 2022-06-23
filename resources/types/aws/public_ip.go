package aws_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/resources/types"
)

type AwsPublicIp struct {
	*types.PublicIp
}

func InitPublicIp(vn *types.PublicIp) resources.ResourceTranslator[*resourcespb.PublicIpResource] {
	return AwsPublicIp{vn}
}

func (r AwsPublicIp) FromState(state *output.TfState) (*resourcespb.PublicIpResource, error) {
	ip := "dryrun"
	if !flags.DryRun {
		values, err := state.GetValues(public_ip.AwsElasticIp{}, r.ResourceId)
		if err != nil {
			return nil, err
		}
		ip = values["public_ip"].(string)
	}

	return &resourcespb.PublicIpResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Ip:          ip,
		GcpOverride: r.Args.GcpOverride,
	}, nil
}

func (r AwsPublicIp) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	return []output.TfBlock{
		public_ip.AwsElasticIp{
			AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
			//Vpc:        true,
		},
	}, nil
}

func (r AwsPublicIp) GetMainResourceName() (string, error) {
	return public_ip.AwsResourceName, nil
}
