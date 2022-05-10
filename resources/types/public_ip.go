package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/public_ip"
	"github.com/multycloud/multy/validate"
)

/*
Notes
AWS: NIC association done on public_ip
Azure: NIC association done on NIC creation
*/

var publicIpMetadata = resources.ResourceMetadata[*resourcespb.PublicIpArgs, *PublicIp, *resourcespb.PublicIpResource]{
	CreateFunc:        CreatePublicIp,
	UpdateFunc:        UpdatePublicIp,
	ReadFromStateFunc: PublicIpFromState,
	ExportFunc: func(r *PublicIp, _ *resources.Resources) (*resourcespb.PublicIpArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewPublicIp,
	AbbreviatedName: "pip",
}

type PublicIp struct {
	resources.ResourceWithId[*resourcespb.PublicIpArgs]
	NetworkInterface *NetworkInterface
}

func (r *PublicIp) GetMetadata() resources.ResourceMetadataInterface {
	return &publicIpMetadata
}

func CreatePublicIp(resourceId string, args *resourcespb.PublicIpArgs, others *resources.Resources) (*PublicIp, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		// todo: maybe put in the same RG as VM?
		rgId := NewRg("pip", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewPublicIp(resourceId, args, others)
}

func UpdatePublicIp(resource *PublicIp, vn *resourcespb.PublicIpArgs, others *resources.Resources) error {
	_, err := NewPublicIp(resource.ResourceId, vn, others)
	return err
}

func PublicIpFromState(resource *PublicIp, _ *output.TfState) (*resourcespb.PublicIpResource, error) {
	return &resourcespb.PublicIpResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:               resource.Args.Name,
		NetworkInterfaceId: resource.Args.NetworkInterfaceId,
	}, nil
}

func NewPublicIp(resourceId string, args *resourcespb.PublicIpArgs, others *resources.Resources) (*PublicIp, error) {
	ni, _, err := resources.GetOptional[*NetworkInterface](resourceId, others, args.NetworkInterfaceId)
	if err != nil {
		return nil, err
	}
	return &PublicIp{
		ResourceWithId: resources.ResourceWithId[*resourcespb.PublicIpArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		NetworkInterface: ni,
	}, nil
}

func (r *PublicIp) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		nid := ""
		if r.NetworkInterface != nil {
			var err error
			nid, err = resources.GetMainOutputId(r.NetworkInterface)
			if err != nil {
				return nil, err
			}
		}
		return []output.TfBlock{
			public_ip.AwsElasticIp{
				AwsResource:        common.NewAwsResource(r.ResourceId, r.Args.Name),
				NetworkInterfaceId: nid,
				//Vpc:        true,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return []output.TfBlock{
			public_ip.AzurePublicIp{
				AzResource: common.NewAzResource(
					r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.GetCommonParameters().ResourceGroupId),
					r.GetCloudSpecificLocation(),
				),
				AllocationMethod: "Static",
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *PublicIp) GetId(cloud commonpb.CloudProvider) string {
	types := map[commonpb.CloudProvider]string{common.AWS: public_ip.AwsResourceName, common.AZURE: public_ip.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.ResourceId)
}

func (r *PublicIp) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	//if r.NetworkInterfaceId != "" && r.InstanceId != "" {
	//	errs = append(errs, r.NewError(r.ResourceId, "instance_id", "cannot set both network_interface_id and instance_id"))
	//}
	return errs
}

func (r *PublicIp) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return public_ip.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return public_ip.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
