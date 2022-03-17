package types

import (
	"fmt"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
)

type KubernetesService struct {
	*resources.CommonResourceParams
	Name      string   `hcl:"name"`
	SubnetIds []string `hcl:"subnet_ids"`
}

func (r *KubernetesService) Validate(ctx resources.MultyContext, cloud common.CloudProvider) (errs []validate.ValidationError) {
	if len(r.SubnetIds) < 2 {
		errs = append(errs, r.NewError("subnet_ids", "at least 2 subnet ids must be provided"))
	}
	associatedNodes := 0
	for _, node := range resources.GetAllResources[*KubernetesServiceNodePool](ctx) {
		if node.ClusterId == resources.GetMainOutputId(r, cloud) && node.IsDefaultPool {
			associatedNodes += 1
		}
	}
	if associatedNodes != 1 {
		errs = append(errs, r.NewError("", fmt.Sprintf("cluster must have exactly 1 default node pool for cloud %s, found %d", cloud, associatedNodes)))
	}
	return errs
}

func (r *KubernetesService) GetMainResourceName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return common.GetResourceName(kubernetes_service.AwsEksCluster{})
	} else if cloud == common.AZURE {
		return common.GetResourceName(kubernetes_service.AzureEksCluster{})
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *KubernetesService) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {
	if cloud == common.AWS {
		roleName := fmt.Sprintf("iam_for_k8cluster_%s", r.Name)
		role := iam.AwsIamRole{
			AwsResource:      common.NewAwsResource(r.GetTfResourceId(cloud), roleName),
			Name:             fmt.Sprintf("iam_for_k8cluster_%s", r.Name),
			AssumeRolePolicy: iam.NewAssumeRolePolicy("eks.amazonaws.com"),
		}
		return []output.TfBlock{
			&role,
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "AmazonEKSClusterPolicy")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.GetTfResourceId(cloud)),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
			},
			iam.AwsIamRolePolicyAttachment{
				AwsResource: common.NewAwsResourceWithIdOnly(fmt.Sprintf("%s_%s", r.GetTfResourceId(cloud), "AmazonEKSVPCResourceController")),
				Role:        fmt.Sprintf("aws_iam_role.%s.name", r.GetTfResourceId(cloud)),
				PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController",
			},
			&kubernetes_service.AwsEksCluster{
				AwsResource: common.NewAwsResource(r.GetTfResourceId(cloud), r.Name),
				RoleArn:     fmt.Sprintf("aws_iam_role.%s.arn", r.GetTfResourceId(cloud)),
				VpcConfig:   kubernetes_service.VpcConfig{SubnetIds: r.SubnetIds},
				Name:        r.Name,
			},
		}
	} else if cloud == common.AZURE {
		var defaultPool *kubernetes_node_pool.AzureKubernetesNodePool
		for _, node := range resources.GetAllResources[*KubernetesServiceNodePool](ctx) {
			if node.ClusterId == resources.GetMainOutputId(r, cloud) && node.IsDefaultPool {
				defaultPool = node.translateAzNodePool(ctx)
				defaultPool.Name = defaultPool.AzResource.Name
				defaultPool.AzResource = nil
				defaultPool.ClusterId = ""
			}
		}
		return []output.TfBlock{
			&kubernetes_service.AzureEksCluster{
				AzResource:      common.NewAzResource(r.GetTfResourceId(cloud), r.Name, rg.GetResourceGroupName(r.ResourceGroupId, cloud), ctx.GetLocationFromCommonParams(r.CommonResourceParams, cloud)),
				DefaultNodePool: defaultPool,
				DnsPrefix:       r.Name,
				Identity:        kubernetes_service.AzureIdentity{Type: "SystemAssigned"},
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}

func (r *KubernetesService) GetOutputValues(cloud common.CloudProvider) map[string]cty.Value {
	switch cloud {
	case common.AWS:
		return map[string]cty.Value{
			"endpoint": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.endpoint}", common.GetResourceName(kubernetes_service.AwsEksCluster{}),
					r.GetTfResourceId(cloud),
				),
			),
			"ca_certificate": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.certificate_authority[0].data}", common.GetResourceName(kubernetes_service.AwsEksCluster{}),
					r.GetTfResourceId(cloud),
				),
			),
		}
	case common.AZURE:
		return map[string]cty.Value{
			"endpoint": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.kube_config.0.host}", common.GetResourceName(kubernetes_service.AzureEksCluster{}),
					r.GetTfResourceId(cloud),
				),
			),
			"ca_certificate": cty.StringVal(
				fmt.Sprintf(
					"${%s.%s.kube_config.0.cluster_ca_certificate}", common.GetResourceName(kubernetes_service.AzureEksCluster{}),
					r.GetTfResourceId(cloud),
				),
			),
		}
	}

	validate.LogInternalError("unknown cloud %s", cloud)
	return nil
}
