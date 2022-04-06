package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/output/subnet"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

/*
Notes:
Azure: New subnets will be associated with a default route table to block traffic to internet
*/

type Subnet struct {
	resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.SubnetArgs]

	VirtualNetwork *VirtualNetwork
}

func Get[T resources.Resource](resources resources.Resources, id string) (T, error) {
	item, exists, err := GetOptional[T](resources, id)
	if err != nil {
		return item, err
	}
	if !exists {
		return item, fmt.Errorf("resource with id %s not found", id)
	}

	return item, nil
}

func GetOptional[T resources.Resource](resources resources.Resources, id string) (T, bool, error) {
	if r, ok := resources[id]; ok {
		if _, okType := r.(T); !okType {
			return *new(T), true, fmt.Errorf("resource with id %s is of the wrong type", id)
		}
		return r.(T), true, nil
	}

	return *new(T), false, nil
}

func NewSubnet(resourceId string, subnet *resourcespb.SubnetArgs, others resources.Resources) (*Subnet, error) {
	s := &Subnet{
		ChildResourceWithId: resources.ChildResourceWithId[*VirtualNetwork, *resourcespb.SubnetArgs]{
			ResourceId: resourceId,
			Args:       subnet,
		},
	}
	vn, err := Get[*VirtualNetwork](others, subnet.VirtualNetworkId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{s.NewValidationError(err.Error(), "virtual_network_id")})
	}
	s.Parent = vn
	s.VirtualNetwork = vn
	return s, nil
}

func (r *Subnet) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		location := r.VirtualNetwork.GetLocation()
		az, err := common.GetAvailabilityZone(location, int(r.Args.AvailabilityZone), r.GetCloud())
		if err != nil {
			return nil, err
		}
		awsSubnet := subnet.AwsSubnet{
			AwsResource:      common.NewAwsResource(r.ResourceId, r.Args.Name),
			CidrBlock:        r.Args.CidrBlock,
			VpcId:            r.VirtualNetwork.GetVirtualNetworkId(),
			AvailabilityZone: az,
		}
		// This flag needs to be set so that eks nodes can connect to the kubernetes cluster
		// https://aws.amazon.com/blogs/containers/upcoming-changes-to-ip-assignment-for-eks-managed-node-groups/
		// How to tell if this subnet is private?
		for _, resource := range resources.GetAllResources[*KubernetesNodePool](ctx) {
			for _, sn := range resource.Subnets {
				if sn.ResourceId == r.ResourceId {
					awsSubnet.MapPublicIpOnLaunch = true
				}
			}
		}
		return []output.TfBlock{awsSubnet}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		var azResources []output.TfBlock
		azSubnet := subnet.AzureSubnet{
			AzResource: &common.AzResource{
				TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				Name:              r.Args.Name,
				ResourceGroupName: rg.GetResourceGroupName(r.VirtualNetwork.Args.GetCommonParameters().GetResourceGroupId()),
			},
			AddressPrefixes:    []string{r.Args.CidrBlock},
			VirtualNetworkName: r.VirtualNetwork.GetVirtualNetworkName(),
		}
		azSubnet.ServiceEndpoints = getServiceEndpointSubnetReferences(ctx, r.ResourceId)
		azResources = append(azResources, azSubnet)

		// there must be a better way to do this
		if !checkSubnetRouteTableAssociated(ctx, r.ResourceId) {
			rt, err := r.VirtualNetwork.GetAssociatedRouteTableId()
			if err != nil {
				return nil, err
			}
			subnetId, err := resources.GetMainOutputId(r)
			if err != nil {
				return nil, err
			}
			rtAssociation := route_table_association.AzureRouteTableAssociation{
				AzResource: &common.AzResource{
					TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
				},
				SubnetId:     subnetId,
				RouteTableId: rt,
			}
			azResources = append(azResources, rtAssociation)
		}

		return azResources, nil
	}
	return nil, fmt.Errorf("cloud %r is not supported for this resource type ", r.GetCloud().String())
}

func (r *Subnet) GetId() string {
	t, _ := r.GetMainResourceName()
	return fmt.Sprintf("%s.%s.id", t, r.ResourceId)
}

func getServiceEndpointSubnetReferences(ctx resources.MultyContext, id string) []string {
	const (
		DATABASE = "Microsoft.Sql"
	)

	serviceEndpoints := map[string]bool{}
	for _, resource := range resources.GetAllResources[*Database](ctx) {
		for _, sn := range resource.Subnets {
			if sn.ResourceId == id {
				serviceEndpoints[DATABASE] = true
			}
		}
	}
	return util.Keys(serviceEndpoints)
}

func checkSubnetRouteTableAssociated(ctx resources.MultyContext, sId string) bool {
	for _, resource := range resources.GetAllResources[*RouteTableAssociation](ctx) {
		if sId == resource.Subnet.ResourceId {
			return true
		}
	}
	return false
}

func (r *Subnet) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	//if vn.Name contains not letters,numbers,_,- { return false }
	//if vn.Name length? { return false }
	//if vn.CidrBlock valid CIDR { return false }
	//if vn.AvailbilityZone valid { return false }
	if len(r.Args.CidrBlock) == 0 { // max len?
		errs = append(errs, r.NewValidationError(fmt.Sprintf("%r cidr_block length is invalid", r.ResourceId), "cidr_block"))
	}

	return errs
}

func (r *Subnet) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case common.AWS:
		return subnet.AwsResourceName, nil
	case common.AZURE:
		return subnet.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %r", r.GetCloud().String())
	}
}

func (r *Subnet) GetCloudSpecificLocation() string {
	return r.VirtualNetwork.GetCloudSpecificLocation()
}
