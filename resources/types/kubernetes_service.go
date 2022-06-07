package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type KubernetesCluster struct {
	resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]

	VirtualNetwork  *VirtualNetwork
	DefaultNodePool *KubernetesNodePool
}

func (r *KubernetesCluster) Create(resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) error {
	return CreateKubernetesCluster(r, resourceId, args, others)
}

func (r *KubernetesCluster) Update(args *resourcespb.KubernetesClusterArgs, _ *resources.Resources) error {
	r.Args = args
	return nil
}

func (r *KubernetesCluster) Import(resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) error {
	return NewKubernetesCluster(r, resourceId, args, others)
}

func (r *KubernetesCluster) Export(_ *resources.Resources) (*resourcespb.KubernetesClusterArgs, bool, error) {
	return r.Args, true, nil
}

func NewKubernetesCluster(cluster *KubernetesCluster, resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) error {
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
	if err != nil {
		return err
	}
	cluster.ResourceWithId = resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]{
		ResourceId: resourceId,
		Args:       args,
	}

	cluster.VirtualNetwork = vn

	cluster.DefaultNodePool = &KubernetesNodePool{}
	err = newKubernetesNodePool(cluster.DefaultNodePool, fmt.Sprintf("%s_default_pool", resourceId), args.DefaultNodePool, others, cluster)
	return err
}

func CreateKubernetesCluster(cluster *KubernetesCluster, resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
		if err != nil {
			return err
		}
		rgId, err := NewRgFromParent("ks", vn.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}
	if args.ServiceCidr == "" {
		args.ServiceCidr = "10.100.0.0/16"
	}

	return NewKubernetesCluster(cluster, resourceId, args, others)
}
func (r *KubernetesCluster) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if r.Args.GetDefaultNodePool() == nil {
		errs = append(errs, r.NewValidationError(fmt.Errorf("cluster must have a default node pool"), "default_node_pool"))
	}
	if r.Args.GetDefaultNodePool().GetClusterId() != "" {
		errs = append(errs, r.NewValidationError(fmt.Errorf("cluster id for default node pool can't be set"), "default_node_pool"))
	}
	errs = append(errs, r.DefaultNodePool.Validate(ctx)...)
	return errs
}
