package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/validate"
)

type KubernetesNodePool struct {
	resources.ChildResourceWithId[*KubernetesCluster, *resourcespb.KubernetesNodePoolArgs]

	KubernetesCluster *KubernetesCluster
	Subnet            *Subnet
}

func (r *KubernetesNodePool) Create(resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) error {
	return NewKubernetesNodePool(r, resourceId, args, others)
}

func (r *KubernetesNodePool) Update(args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) error {
	return NewKubernetesNodePool(r, r.ResourceId, args, others)
}

func (r *KubernetesNodePool) Import(resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) error {
	return NewKubernetesNodePool(r, resourceId, args, others)
}

func (r *KubernetesNodePool) Export(_ *resources.Resources) (*resourcespb.KubernetesNodePoolArgs, bool, error) {
	return r.Args, true, nil
}

func NewKubernetesNodePool(r *KubernetesNodePool, resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) error {
	cluster, err := resources.Get[*KubernetesCluster](resourceId, others, args.ClusterId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, r.ResourceId, "cluster_id"))
	}
	return newKubernetesNodePool(r, resourceId, args, others, cluster)
}

func newKubernetesNodePool(knp *KubernetesNodePool, resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources, cluster *KubernetesCluster) error {
	knp.ChildResourceWithId = resources.NewChildResource(resourceId, cluster, args)

	if args.StartingNodeCount == 0 {
		knp.Args.StartingNodeCount = args.MinNodeCount
	}

	knp.KubernetesCluster = cluster

	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return err
	}
	knp.Subnet = subnet
	return nil
}

func (r *KubernetesNodePool) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if r.Args.MinNodeCount < 1 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("node pool must have a min node count of at least 1"), "min_node_count"))
	}
	if r.Args.MaxNodeCount < 1 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("node pool must have a max node count of at least 1"), "max_node_count"))
	}
	if r.Args.MinNodeCount > r.Args.MaxNodeCount {
		errs = append(errs, r.NewValidationError(fmt.Errorf("min_node_count must be lower or equal to max_node_count"), "min_node_count"))
	}
	if r.Args.StartingNodeCount < r.Args.MinNodeCount || r.Args.StartingNodeCount > r.Args.MaxNodeCount {
		errs = append(errs, r.NewValidationError(fmt.Errorf("starting_node_count must be between min and max node count"), "starting_node_count"))
	}
	if r.Args.VmSize == commonpb.VmSize_UNKNOWN_VM_SIZE {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown vm size"), "vm_size"))
	}
	if r.Subnet.VirtualNetwork.ResourceId != r.KubernetesCluster.VirtualNetwork.ResourceId {
		errs = append(errs, r.NewValidationError(fmt.Errorf("subnet must be in the same virtual network as the cluster"), "subnet_id"))
	}
	if r.Subnet != r.KubernetesCluster.DefaultNodePool.Subnet {
		errs = append(errs, r.NewValidationError(fmt.Errorf("different subnets for node pools within the same cluster are not yet supported"), "subnet_id"))
	}
	if r.GetCloud() == commonpb.CloudProvider_GCP {
		numZones := 3
		if len(r.Args.AvailabilityZone) > 0 {
			numZones = len(r.Args.AvailabilityZone)
		}

		if int(r.Args.StartingNodeCount)%numZones != 0 {
			errs = append(errs, r.NewValidationError(
				fmt.Errorf("for gcp, starting node count must be a multiple of the number of availability zones (%d, %d, %d, etc)",
					numZones, numZones*2, numZones*3), "starting_node_count"))
		}

		if int(r.Args.MinNodeCount)%numZones != 0 {
			errs = append(errs, r.NewValidationError(
				fmt.Errorf("for gcp, minimum number of nodes must be a multiple of the number of availability zones (%d, %d, %d, etc)",
					numZones, numZones*2, numZones*3), "min_node_count"))
		}

		if int(r.Args.MaxNodeCount)%numZones != 0 {
			errs = append(errs, r.NewValidationError(
				fmt.Errorf("for gcp, maximum number of nodes must be a multiple of the number of availability zones (%d, %d, %d, etc)",
					numZones, numZones*2, numZones*3), "max_node_count"))
		}
	}

	return errs
}

func (r *KubernetesNodePool) ParseCloud(args *resourcespb.KubernetesNodePoolArgs) commonpb.CloudProvider {
	return common.ParseCloudFromResourceId(args.ClusterId)
}
