package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/util"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
)

type KubernetesCluster struct {
	resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]

	Subnets         []*Subnet
	DefaultNodePool *KubernetesNodePool
}

func NewKubernetesCluster(resourceId string, args *resourcespb.KubernetesClusterArgs, others resources.Resources) (*KubernetesCluster, error) {
	subnets, err := util.MapSliceValuesErr(args.SubnetIds, func(subnetId string) (*Subnet, error) {
		return resources.Get[*Subnet](resourceId, others, subnetId)
	})
	if err != nil {
		return nil, err
	}
	cluster := &KubernetesCluster{
		ResourceWithId: resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		Subnets: subnets,
	}

	cluster.DefaultNodePool, err = newKubernetesNodePool(fmt.Sprintf("%s_default_pool", resourceId), args.DefaultNodePool, others, cluster)
	return cluster, err
}

func (r *KubernetesCluster) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if len(r.Subnets) < 2 {
		errs = append(errs, r.NewValidationError(fmt.Errorf("at least 2 subnet ids must be provided"), "subnet_ids"))
	}
	if r.Args.GetDefaultNodePool() == nil {
		errs = append(errs, r.NewValidationError(fmt.Errorf("cluster must have a default node pool"), "default_node_pool"))
	}
	if r.Args.GetDefaultNodePool().GetClusterId() != "" {
		errs = append(errs, r.NewValidationError(fmt.Errorf("cluster id for default node pool can't be set"), "default_node_pool"))
	}
	errs = append(errs, r.DefaultNodePool.Validate(ctx)...)
	return errs
}

func (r *KubernetesCluster) GetMainResourceName() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return output.GetResourceName(kubernetes_service.AwsEksCluster{}), nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return output.GetResourceName(kubernetes_service.AzureEksCluster{}), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesCluster) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {

	subnetIds, err := util.MapSliceValuesErr(r.Subnets, func(v *Subnet) (string, error) {
		return resources.GetMainOutputId(v)
	})
	if err != nil {
		return nil, err
	}
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var outputs []output.TfBlock
		defaultNodePoolResources, err := r.DefaultNodePool.Translate(ctx)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, defaultNodePoolResources...)
		roleName := fmt.Sprintf("iam_for_k8cluster_%s", r.Args.Name)
		role := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.ResourceId, roleName),
			Name:             fmt.Sprintf("iam_for_k8cluster_%s", r.Args.Name),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("eks.amazonaws.com"),
		}
		outputs = append(outputs, &role,
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSClusterPolicy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSVPCResourceController")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController",
			},
			&kubernetes_service.AwsEksCluster{
				AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
				RoleArn:     fmt.Sprintf("aws_iam_role.%s.arn", r.ResourceId),
				VpcConfig:   kubernetes_service.VpcConfig{SubnetIds: subnetIds},
				Name:        r.Args.Name,
			})
		return outputs, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		var defaultPool *kubernetes_node_pool.AzureKubernetesNodePool
		defaultPool, err = r.DefaultNodePool.translateAzNodePool()
		if err != nil {
			return nil, err
		}
		defaultPool.Name = defaultPool.AzResource.Name
		defaultPool.AzResource = nil
		defaultPool.ClusterId = ""

		return []output.TfBlock{
			&kubernetes_service.AzureEksCluster{
				AzResource:      common.NewAzResource(r.ResourceId, r.Args.Name, rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId), r.GetCloudSpecificLocation()),
				DefaultNodePool: defaultPool,
				DnsPrefix:       common.UniqueId(r.Args.Name, "aks", common.LowercaseAlphanumericFormatFunc),
				Identity:        kubernetes_service.AzureIdentity{Type: "SystemAssigned"},
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesCluster) GetOutputValues(cloud commonpb.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"endpoint": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.endpoint}", output.GetResourceName(kubernetes_service.AwsEksCluster{}),
					r.ResourceId,
				),
			),
			"ca_certificate": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.certificate_authority[0].data}", output.GetResourceName(kubernetes_service.AwsEksCluster{}),
					r.ResourceId,
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"endpoint": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.kube_config.0.host}", output.GetResourceName(kubernetes_service.AzureEksCluster{}),
					r.ResourceId,
				),
			),
			"ca_certificate": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.kube_config.0.cluster_ca_certificate}", output.GetResourceName(kubernetes_service.AzureEksCluster{}),
					r.ResourceId,
				),
			),
		}
	}

	return nil
}
