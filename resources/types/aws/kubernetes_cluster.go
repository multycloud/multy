package aws_resources

import (
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	"github.com/multycloud/multy/resources/output/route_table"
	"github.com/multycloud/multy/resources/output/route_table_association"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/validate"
	"gopkg.in/yaml.v3"
	"net"
)

type AwsKubernetesCluster struct {
	*types.KubernetesCluster
}

func InitKubernetesCluster(r *types.KubernetesCluster) resources.ResourceTranslator[*resourcespb.KubernetesClusterResource] {
	return AwsKubernetesCluster{r}
}

func createKubeConfig(clusterName string, certData string, endpoint string, awsRegion string) (string, error) {
	username := fmt.Sprintf("clusterUser_%s", clusterName)
	kubeConfig := &kubernetes_service.KubeConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Clusters: []kubernetes_service.NamedKubeConfigCluster{
			{
				Name: clusterName,
				Cluster: kubernetes_service.KubeConfigCluster{
					CertificateAuthorityData: certData,
					Server:                   endpoint,
				},
			},
		},
		Contexts: []kubernetes_service.NamedKubeConfigContext{
			{
				Name: clusterName,
				Context: kubernetes_service.KubeConfigContext{
					User:    username,
					Cluster: clusterName,
				},
			},
		},
		Users: []kubernetes_service.KubeConfigUser{
			{
				Name: username,
				User: struct {
					Exec kubernetes_service.KubeConfigExec `yaml:"exec"`
				}{
					Exec: kubernetes_service.KubeConfigExec{
						ApiVersion:      "client.authentication.k8s.io/v1beta1",
						Command:         "aws",
						Args:            []string{"--region", awsRegion, "eks", "get-token", "--cluster-name", clusterName},
						InteractiveMode: "IfAvailable",
						InstallHint:     "Install aws cli for use with kubectl by following\n\t\t\t\thttps://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html",
					},
				},
			},
		},
		CurrentContext: clusterName,
	}

	s, err := yaml.Marshal(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("unable to marshal kube config, %s", err)
	}

	return string(s), nil
}

func (r AwsKubernetesCluster) FromState(state *output.TfState) (*resourcespb.KubernetesClusterResource, error) {
	result := &resourcespb.KubernetesClusterResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:             r.Args.Name,
		ServiceCidr:      r.Args.ServiceCidr,
		VirtualNetworkId: r.Args.VirtualNetworkId,
		GcpOverride:      r.Args.GcpOverride,
	}
	result.Endpoint = "dryrun"
	if !flags.DryRun {
		cluster, err := output.GetParsedById[kubernetes_service.AwsEksCluster](state, r.ResourceId)
		if err != nil {
			return nil, err
		}
		result.Endpoint = cluster.Endpoint
		result.CaCertificate = cluster.CertificateAuthority[0].Data
		kubeCgfRaw, err := createKubeConfig(r.Args.Name, result.CaCertificate, result.Endpoint, r.GetCloudSpecificLocation())
		if err != nil {
			return nil, err
		}
		result.KubeConfigRaw = kubeCgfRaw
	}

	var err error
	result.DefaultNodePool, err = AwsKubernetesNodePool{r.DefaultNodePool}.FromState(state)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r AwsKubernetesCluster) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var outputs []output.TfBlock
	defaultNodePoolResources, err := AwsKubernetesNodePool{r.DefaultNodePool}.Translate(ctx)
	if err != nil {
		return nil, err
	}
	outputs = append(outputs, defaultNodePoolResources...)
	subnets, subnetResources, err := r.getAwsSubnets()
	if err != nil {
		return nil, err
	}
	var subnetIds []string
	for _, s := range subnets {
		subnetIds = append(subnetIds, fmt.Sprintf("%s.%s.id", output.GetResourceName(s), s.ResourceId))
	}

	outputs = append(outputs, subnetResources...)
	var deps []string
	for _, s := range subnetResources {
		// todo: get the id without casting
		deps = append(deps, fmt.Sprintf("%s.%s", output.GetResourceName(s), s.GetResourceId()))
	}

	roleName := fmt.Sprintf("multy-k8cluster-%s-role", r.Args.Name)
	role := iam.AwsIamRole{
		AwsResource:      common.NewAwsResource(r.ResourceId, roleName),
		Name:             roleName,
		AssumeRolePolicy: iam.NewAssumeRolePolicy("eks.amazonaws.com"),
	}

	policy1Id := fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSClusterPolicy")
	policy1 := iam.AwsIamRolePolicyAttachment{
		AwsResource: common.NewAwsResourceWithIdOnly(policy1Id),
		Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
		PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
	}
	policy2Id := fmt.Sprintf("%s_%s", r.ResourceId, "AmazonEKSVPCResourceController")
	policy2 := iam.AwsIamRolePolicyAttachment{
		AwsResource: common.NewAwsResourceWithIdOnly(policy2Id),
		Role:        fmt.Sprintf("aws_iam_role.%s.name", r.ResourceId),
		PolicyArn:   "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController",
	}
	deps = append(deps, fmt.Sprintf("%s.%s", output.GetResourceName(policy1), policy1Id),
		fmt.Sprintf("%s.%s", output.GetResourceName(policy1), policy2Id))

	outputs = append(outputs, &role,
		policy1,
		policy2,
		&kubernetes_service.AwsEksCluster{
			AwsResource: common.NewAwsResourceWithDeps(r.ResourceId, r.Args.Name, deps),
			RoleArn:     fmt.Sprintf("aws_iam_role.%s.arn", r.ResourceId),
			VpcConfig:   kubernetes_service.VpcConfig{SubnetIds: subnetIds, EndpointPrivateAccess: true},
			Name:        r.Args.Name,
			KubernetesNetworkConfig: kubernetes_service.KubernetesNetworkConfig{
				ServiceIpv4Cidr: r.Args.ServiceCidr,
			},
		})
	return outputs, nil
}

func (r AwsKubernetesCluster) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_service.AwsEksCluster{}), nil
}

func (r AwsKubernetesCluster) getAwsSubnets() ([]subnet.AwsSubnet, []output.TfBlock, error) {
	block := r.VirtualNetwork.Args.CidrBlock
	_, vnNet, err := net.ParseCIDR(block)
	if err != nil {
		return nil, nil, err
	}
	tempSubnet, _ := cidr.NextSubnet(vnNet, 31)
	subnetBlock1, _ := cidr.PreviousSubnet(tempSubnet, 28)
	subnetBlock2, _ := cidr.PreviousSubnet(subnetBlock1, 28)
	validationError := validate.ValidationError{
		ErrorMessage: fmt.Sprintf("Not enough availabilty zones available in region %s. Kubernetes clusters in AWS require 2 availabilty zones.", r.VirtualNetwork.GetLocation()),
		ResourceId:   r.ResourceId,
		FieldName:    "virtual_network_id",
	}
	az1, err := common.GetAvailabilityZone(r.VirtualNetwork.GetLocation(), 1, r.GetCloud())
	if err != nil {
		return nil, nil, errors.ValidationErrors([]validate.ValidationError{validationError})
	}
	az2, err := common.GetAvailabilityZone(r.VirtualNetwork.GetLocation(), 2, r.GetCloud())
	if err != nil {
		return nil, nil, errors.ValidationErrors([]validate.ValidationError{validationError})
	}

	vpcId, err := resources.GetMainOutputId(AwsVirtualNetwork{r.VirtualNetwork})
	if err != nil {
		return nil, nil, err
	}
	rt := route_table.AwsRouteTable{
		AwsResource: common.NewAwsResource(r.ResourceId+"_public_rt", r.Args.Name+"_public_rt"),
		VpcId:       vpcId,
		Routes: []route_table.AwsRouteTableRoute{
			{
				CidrBlock: "0.0.0.0/0",
				GatewayId: AwsVirtualNetwork{r.VirtualNetwork}.GetAssociatedInternetGateway(),
			},
		},
	}

	subnet1 := subnet.AwsSubnet{
		AwsResource:      common.NewAwsResource(r.ResourceId+"_public_subnet", r.Args.Name+"_public_subnet"),
		CidrBlock:        subnetBlock1.String(),
		VpcId:            fmt.Sprintf("%s.%s.id", virtual_network.AwsResourceName, r.VirtualNetwork.ResourceId),
		AvailabilityZone: az1,
	}
	subnet2 := subnet.AwsSubnet{
		AwsResource:      common.NewAwsResource(r.ResourceId+"_private_subnet", r.Args.Name+"_private_subnet"),
		CidrBlock:        subnetBlock2.String(),
		VpcId:            fmt.Sprintf("%s.%s.id", virtual_network.AwsResourceName, r.VirtualNetwork.ResourceId),
		AvailabilityZone: az2,
	}

	rta := route_table_association.AwsRouteTableAssociation{
		AwsResource:  common.NewAwsResourceWithIdOnly(r.ResourceId + "_public_rta"),
		SubnetId:     fmt.Sprintf("%s.%s.id", output.GetResourceName(subnet1), subnet1.ResourceId),
		RouteTableId: fmt.Sprintf("%s.%s.id", output.GetResourceName(rt), rt.ResourceId),
	}

	return []subnet.AwsSubnet{subnet1, subnet2}, []output.TfBlock{subnet1, subnet2, rt, rta}, nil
}
