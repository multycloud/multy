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

	Subnets []*Subnet
}

func NewKubernetesCluster(resourceId string, args *resourcespb.KubernetesClusterArgs, others resources.Resources) (*KubernetesCluster, error) {
	subnets, err := util.MapSliceValuesErr(args.SubnetIds, func(subnetId string) (*Subnet, error) {
		return Get[*Subnet](others, subnetId)
	})
	if err != nil {
		return nil, err
	}
	return &KubernetesCluster{
		ResourceWithId: resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		Subnets: subnets,
	}, nil
}

func (r *KubernetesCluster) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	if len(r.Subnets) < 2 {
		errs = append(errs, r.NewValidationError("at least 2 subnet ids must be provided", "subnet_ids"))
	}
	associatedNodes := 0
	for _, node := range resources.GetAllResourcesInCloud[*KubernetesNodePool](ctx, r.GetCloud()) {
		if node.KubernetesCluster.ResourceId == r.ResourceId && node.Args.IsDefaultPool {
			associatedNodes += 1
		}
	}
	if associatedNodes != 1 {
		errs = append(errs, r.NewValidationError(fmt.Sprintf("cluster must have exactly 1 default node pool, found %d", associatedNodes), ""))
	}
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
		roleName := fmt.Sprintf("iam_for_k8cluster_%s", r.Args.Name)
		role := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.ResourceId, roleName),
			Name:             fmt.Sprintf("iam_for_k8cluster_%s", r.Args.Name),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("eks.amazonaws.com"),
		}
		return []output.TfBlock{
			&role,
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
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		var defaultPool *kubernetes_node_pool.AzureKubernetesNodePool
		for _, node := range resources.GetAllResources[*KubernetesNodePool](ctx) {
			if node.KubernetesCluster.ResourceId == r.ResourceId && node.Args.IsDefaultPool {
				defaultPool, err = node.translateAzNodePool()
				if err != nil {
					return nil, err
				}
				defaultPool.Name = defaultPool.AzResource.Name
				defaultPool.AzResource = nil
				defaultPool.ClusterId = ""
			}
		}
		return []output.TfBlock{
			&kubernetes_service.AzureEksCluster{
				AzResource:      common.NewAzResource(r.ResourceId, r.Args.Name, rg.GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId), r.GetCloudSpecificLocation()),
				DefaultNodePool: defaultPool,
				DnsPrefix:       r.Args.Name,
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
