package types

import (
	"fmt"
	"multy-go/resources"
	"multy-go/resources/common"
	"multy-go/resources/output"
	"multy-go/resources/output/iam"
	"multy-go/resources/output/kubernetes_service"
	"multy-go/validate"
)

type KubernetesService struct {
	*resources.CommonResourceParams
	Name      string   `hcl:"name"`
	SubnetIds []string `hcl:"subnet_ids"`
}

func (r *KubernetesService) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if len(r.SubnetIds) < 2 {
		errs = append(errs, r.NewError("subnet_ids", "at least 2 subnet ids must be provided"))
	}
	associatedNodes := 0
	for _, node := range resources.GetAllResources[*KubernetesServiceNodePool](ctx) {
		if node.Name == r.Name {
			associatedNodes += 1
		}
	}
	if associatedNodes == 0 {
		errs = append(errs, r.NewError("", "at least 1 node pool must be associated with this cluster, found 0"))
	}
	return errs
}

func (r *KubernetesService) GetMainResourceName(cloud common.CloudProvider) string {
	if cloud == common.AWS {
		return common.GetResourceName(kubernetes_service.AwsEksCluster{})
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return ""
}

func (r *KubernetesService) Translate(cloud common.CloudProvider, ctx resources.MultyContext) []output.TfBlock {

	roleName := fmt.Sprintf("iam_for_k8cluster_%s", r.Name)
	role := iam.AwsIamRole{
		AwsResource:      common.NewAwsResource(r.GetTfResourceId(cloud), roleName),
		Name:             fmt.Sprintf("iam_for_k8cluster_%s", r.Name),
		AssumeRolePolicy: iam.NewAssumeRolePolicy("eks.amazonaws.com"),
	}

	if cloud == common.AWS {
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
				AwsResource: common.NewAwsResource(r.ResourceId, r.Name),
				RoleArn:     fmt.Sprintf("aws_iam_role.%s.arn", r.GetTfResourceId(cloud)),
				VpcConfig:   kubernetes_service.VpcConfig{SubnetIds: r.SubnetIds},
				Name:        r.Name,
			},
		}
	}
	validate.LogInternalError("cloud %s is not supported for this resource type ", cloud)
	return nil
}
