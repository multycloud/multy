package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
)

var kubernetesNodePoolMetadata = resources.ResourceMetadata[*resourcespb.KubernetesNodePoolArgs, *KubernetesNodePool, *resourcespb.KubernetesNodePoolResource]{
	CreateFunc:        CreateKubernetesNodePool,
	UpdateFunc:        UpdateKubernetesNodePool,
	ReadFromStateFunc: KubernetesNodePoolFromState,
	ExportFunc: func(r *KubernetesNodePool, _ *resources.Resources) (*resourcespb.KubernetesNodePoolArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewKubernetesNodePool,
	AbbreviatedName: "ks",
}

type KubernetesNodePool struct {
	resources.ChildResourceWithId[*KubernetesCluster, *resourcespb.KubernetesNodePoolArgs]

	KubernetesCluster *KubernetesCluster
	Subnets           []*Subnet // azure??
}

func (r *KubernetesNodePool) GetMetadata() resources.ResourceMetadataInterface {
	return &kubernetesNodePoolMetadata
}

func CreateKubernetesNodePool(resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) (*KubernetesNodePool, error) {
	return NewKubernetesNodePool(resourceId, args, others)
}

func UpdateKubernetesNodePool(resource *KubernetesNodePool, vn *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func NewKubernetesNodePool(resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources) (*KubernetesNodePool, error) {
	cluster, err := resources.Get[*KubernetesCluster](resourceId, others, args.ClusterId)
	if err != nil {
		return nil, errors.ValidationErrors([]validate.ValidationError{{
			ErrorMessage: err.Error(),
			ResourceId:   resourceId,
			FieldName:    "cluster_id",
		}})
	}
	return newKubernetesNodePool(resourceId, args, others, cluster)
}

func KubernetesNodePoolFromState(resource *KubernetesNodePool, _ *output.TfState) (*resourcespb.KubernetesNodePoolResource, error) {
	return &resourcespb.KubernetesNodePoolResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  resource.ResourceId,
			NeedsUpdate: false,
		},
		Name:              resource.Args.Name,
		SubnetIds:         resource.Args.SubnetIds,
		ClusterId:         resource.Args.ClusterId,
		StartingNodeCount: resource.Args.StartingNodeCount,
		MinNodeCount:      resource.Args.MinNodeCount,
		MaxNodeCount:      resource.Args.MaxNodeCount,
		VmSize:            resource.Args.VmSize,
		DiskSizeGb:        resource.Args.DiskSizeGb,
		Labels:            resource.Args.Labels,
	}, nil
}

func newKubernetesNodePool(resourceId string, args *resourcespb.KubernetesNodePoolArgs, others *resources.Resources, cluster *KubernetesCluster) (*KubernetesNodePool, error) {
	knp := &KubernetesNodePool{
		ChildResourceWithId: resources.ChildResourceWithId[*KubernetesCluster, *resourcespb.KubernetesNodePoolArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
	}

	if args.StartingNodeCount == 0 {
		knp.Args.StartingNodeCount = args.MinNodeCount
	}

	knp.Parent = cluster
	knp.KubernetesCluster = cluster

	subnets, err := util.MapSliceValuesErr(args.SubnetIds, func(subnetId string) (*Subnet, error) {
		return resources.Get[*Subnet](resourceId, others, subnetId)
	})
	if err != nil {
		return nil, err
	}
	knp.Subnets = subnets
	return knp, nil
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

	return errs
}

func (r *KubernetesNodePool) GetMainResourceName() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return output.GetResourceName(kubernetes_node_pool.AwsKubernetesNodeGroup{}), nil
	}
	if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return output.GetResourceName(kubernetes_node_pool.AzureKubernetesNodePool{}), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesNodePool) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	subnetIds, err := util.MapSliceValuesErr(r.Subnets, func(v *Subnet) (string, error) {
		return resources.GetMainOutputId(v)
	})
	if err != nil {
		return nil, err
	}
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		roleName := fmt.Sprintf("iam_for_k8nodepool_%s", r.Args.Name)
		role := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.ResourceId, roleName),
			Name:             fmt.Sprintf("iam_for_k8nodepool_%s", r.Args.Name),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("ec2.amazonaws.com"),
		}
		clusterId, err := resources.GetMainOutputId(r.KubernetesCluster)
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			&role,
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSWorkerNodePolicy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKS_CNI_Policy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEC2ContainerRegistryReadOnly")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
			},
			&kubernetes_node_pool.AwsKubernetesNodeGroup{
				AwsResource:   common.NewAwsResourceWithIdOnly(r.ResourceId),
				ClusterName:   clusterId,
				NodeGroupName: r.Args.Name,
				NodeRoleArn:   fmt.Sprintf("aws_iam_role.%s.arn", r.ResourceId),
				SubnetIds:     subnetIds,
				ScalingConfig: kubernetes_node_pool.ScalingConfig{
					DesiredSize: int(r.Args.StartingNodeCount),
					MaxSize:     int(r.Args.MaxNodeCount),
					MinSize:     int(r.Args.MinNodeCount),
				},
				Labels:        r.Args.Labels,
				InstanceTypes: []string{common.VMSIZE[r.Args.VmSize][r.GetCloud()]},
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		pool, err := r.translateAzNodePool()
		if err != nil {
			return nil, err
		}
		return []output.TfBlock{
			pool,
		}, nil

	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesNodePool) translateAzNodePool() (*kubernetes_node_pool.AzureKubernetesNodePool, error) {
	clusterId, err := resources.GetMainOutputId(r.KubernetesCluster)
	if err != nil {
		return nil, err
	}
	subnetId, err := resources.GetMainOutputId(r.Subnets[0])
	if err != nil {
		return nil, err
	}
	return &kubernetes_node_pool.AzureKubernetesNodePool{
		AzResource: &common.AzResource{
			TerraformResource: output.TerraformResource{ResourceId: r.ResourceId},
			Name:              r.Args.Name,
		},
		ClusterId:         clusterId,
		NodeCount:         int(r.Args.StartingNodeCount),
		MaxSize:           int(r.Args.MaxNodeCount),
		MinSize:           int(r.Args.MinNodeCount),
		Labels:            r.Args.Labels,
		EnableAutoScaling: true,
		VmSize:            common.VMSIZE[r.Args.VmSize][r.GetCloud()],
		PodSubnetId:       subnetId,
	}, nil
}

func (r *KubernetesNodePool) GetCloudSpecificLocation() string {
	return r.KubernetesCluster.GetCloudSpecificLocation()
}
