package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/iam"
	"multy-go/resources/output/kubernetes_node_pool"
	"multy-go/resources/output/kubernetes_service"
	rg "multy-go/resources/resource_group"
	"multy-go/validate"
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
