package aws_resources

import (
	"encoding/json"
	"fmt"
	"github.com/multy-dev/hclencoder"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/types"
)

type AwsVaultAccessPolicy struct {
	*types.VaultAccessPolicy
}

func InitVaultAccessPolicy(vn *types.VaultAccessPolicy) resources.ResourceTranslator[*resourcespb.VaultAccessPolicyResource] {
	return AwsVaultAccessPolicy{vn}
}

func (r AwsVaultAccessPolicy) FromState(state *output.TfState) (*resourcespb.VaultAccessPolicyResource, error) {
	return &resourcespb.VaultAccessPolicyResource{
		CommonParameters: &commonpb.CommonChildResourceParameters{
			ResourceId:  r.ResourceId,
			NeedsUpdate: false,
		},
		VaultId:  r.Args.VaultId,
		Identity: r.Args.Identity,
		Access:   r.Args.Access,
	}, nil
}

func (r AwsVaultAccessPolicy) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	var result []output.TfBlock
	result = append(result, AwsCallerIdentityData{TerraformDataSource: &output.TerraformDataSource{ResourceId: r.ResourceId}})

	policy, err := json.Marshal(iam.AwsIamPolicy{
		Statement: []iam.AwsIamPolicyStatement{{
			Action:   r.getAccessPolicyRules(),
			Effect:   "Allow",
			Resource: fmt.Sprintf("arn:aws:ssm:%s:${data.aws_caller_identity.%s.account_id}:parameter/%s/*", r.Vault.GetCloudSpecificLocation(), r.ResourceId, r.Vault.Args.Name),
		}, {
			Action:   []string{"ssm:DescribeParameters"},
			Effect:   "Allow",
			Resource: "*",
		}},
		Version: "2012-10-17",
	})

	if err != nil {
		return nil, fmt.Errorf("unable to encode aws policy: %s", err)
	}

	result = append(result, &iam.AwsIamRolePolicy{
		AwsResource: common.NewAwsResource(r.ResourceId, r.ResourceId),
		Name:        r.ResourceId,
		// we need to have an expression here because we use template strings within the policy json
		Policy: fmt.Sprintf("\"%s\"", hclencoder.EscapeString(string(policy))),
	})

	result = append(result, &iam.AwsIamRolePolicyAttachmentForVap{
		AwsResource: common.NewAwsResourceWithIdOnly(r.ResourceId),
		Role:        r.Args.Identity,
		PolicyArn:   fmt.Sprintf("%s.%s.arn", output.GetResourceName(iam.AwsIamRolePolicy{}), r.ResourceId),
	})

	return result, nil
}

func (r AwsVaultAccessPolicy) getAccessPolicyRules() []string {
	switch r.Args.Access {
	case resourcespb.VaultAccess_READ:
		return []string{"ssm:GetParameter*"}
	case resourcespb.VaultAccess_WRITE:
		return []string{"ssm:PutParameter", "ssm:DeleteParameter"}
	case resourcespb.VaultAccess_OWNER:
		return []string{"ssm:*"}
	default:
		return nil
	}
}

func (r AwsVaultAccessPolicy) GetMainResourceName() (string, error) {
	return "", fmt.Errorf("vault access policy doesn't output any resources in AWS")
}
