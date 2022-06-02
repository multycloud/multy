package types

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
	"github.com/multycloud/multy/validate"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v3"
	"net"
)

type KubeConfig struct {
	ApiVersion     string                   `yaml:"apiVersion"`
	Clusters       []NamedKubeConfigCluster `yaml:"clusters"`
	Contexts       []NamedKubeConfigContext `yaml:"contexts"`
	CurrentContext string                   `yaml:"current-context"`
	Users          []KubeConfigUser         `yaml:"users"`
	Kind           string                   `yaml:"kind"`
}

type NamedKubeConfigCluster struct {
	Name    string            `yaml:"name"`
	Cluster KubeConfigCluster `yaml:"cluster"`
}

type KubeConfigCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type NamedKubeConfigContext struct {
	Name    string            `yaml:"name"`
	Context KubeConfigContext `yaml:"context"`
}

type KubeConfigContext struct {
	User    string `yaml:"user"`
	Cluster string `yaml:"cluster"`
}

type KubeConfigUser struct {
	Name string `yaml:"name"`
	User struct {
		Exec KubeConfigExec `yaml:"exec"`
	} `yaml:"user"`
}

type KubeConfigExec struct {
	ApiVersion      string   `yaml:"apiVersion"`
	Command         string   `yaml:"command"`
	Args            []string `yaml:"args"`
	InteractiveMode string   `yaml:"interactiveMode"`
}

var kubernetesClusterMetadata = resources.ResourceMetadata[*resourcespb.KubernetesClusterArgs, *KubernetesCluster, *resourcespb.KubernetesClusterResource]{
	CreateFunc:        CreateKubernetesCluster,
	UpdateFunc:        UpdateKubernetesCluster,
	ReadFromStateFunc: KubernetesClusterFromState,
	ExportFunc: func(r *KubernetesCluster, _ *resources.Resources) (*resourcespb.KubernetesClusterArgs, bool, error) {
		return r.Args, true, nil
	},
	ImportFunc:      NewKubernetesCluster,
	AbbreviatedName: "ks",
}

type KubernetesCluster struct {
	resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]

	VirtualNetwork  *VirtualNetwork
	DefaultNodePool *KubernetesNodePool
}

func (r *KubernetesCluster) GetMetadata() resources.ResourceMetadataInterface {
	return &kubernetesClusterMetadata
}

func NewKubernetesCluster(resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) (*KubernetesCluster, error) {
	vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
	if err != nil {
		return nil, err
	}
	cluster := &KubernetesCluster{
		ResourceWithId: resources.ResourceWithId[*resourcespb.KubernetesClusterArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		VirtualNetwork: vn,
	}

	cluster.DefaultNodePool, err = newKubernetesNodePool(fmt.Sprintf("%s_default_pool", resourceId), args.DefaultNodePool, others, cluster)
	return cluster, err
}

func CreateKubernetesCluster(resourceId string, args *resourcespb.KubernetesClusterArgs, others *resources.Resources) (*KubernetesCluster, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		vn, err := resources.Get[*VirtualNetwork](resourceId, others, args.VirtualNetworkId)
		if err != nil {
			return nil, err
		}
		rgId, err := NewRgFromParent("ks", vn.Args.CommonParameters.ResourceGroupId, others,
			args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return nil, err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}
	if args.ServiceCidr == "" {
		args.ServiceCidr = "10.100.0.0/16"
	}

	return NewKubernetesCluster(resourceId, args, others)
}

func UpdateKubernetesCluster(resource *KubernetesCluster, vn *resourcespb.KubernetesClusterArgs, others *resources.Resources) error {
	resource.Args = vn
	return nil
}

func createKubeConfig(clusterName string, certData string, endpoint string, awsRegion string) (string, error) {
	username := fmt.Sprintf("clusterUser_%s", clusterName)
	kubeConfig := &KubeConfig{
		ApiVersion: "v1",
		Kind:       "Config",
		Clusters: []NamedKubeConfigCluster{
			{
				Name: clusterName,
				Cluster: KubeConfigCluster{
					CertificateAuthorityData: certData,
					Server:                   endpoint,
				},
			},
		},
		Contexts: []NamedKubeConfigContext{
			{
				Name: clusterName,
				Context: KubeConfigContext{
					User:    username,
					Cluster: clusterName,
				},
			},
		},
		Users: []KubeConfigUser{
			{
				Name: username,
				User: struct {
					Exec KubeConfigExec `yaml:"exec"`
				}{
					Exec: KubeConfigExec{
						ApiVersion:      "client.authentication.k8s.io/v1alpha1",
						Command:         "aws",
						Args:            []string{"--region", awsRegion, "eks", "get-token", "--cluster-name", clusterName},
						InteractiveMode: "IfAvailable",
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

func KubernetesClusterFromState(resource *KubernetesCluster, state *output.TfState) (*resourcespb.KubernetesClusterResource, error) {
	var err error
	result := &resourcespb.KubernetesClusterResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:             resource.Args.Name,
		ServiceCidr:      resource.Args.ServiceCidr,
		VirtualNetworkId: resource.Args.VirtualNetworkId,
	}
	result.Endpoint = "dryrun"
	if !flags.DryRun {
		switch resource.GetCloud() {
		case commonpb.CloudProvider_AWS:
			values, err := state.GetValues(kubernetes_service.AwsEksCluster{}, resource.ResourceId)
			if err != nil {
				return nil, err
			}
			result.Endpoint = values["endpoint"].(string)
			result.CaCertificate = values["certificate_authority"].([]interface{})[0].(map[string]interface{})["data"].(string)
			kubeCgfRaw, err := createKubeConfig(resource.Args.Name, result.CaCertificate, result.Endpoint, resource.GetCloudSpecificLocation())
			if err != nil {
				return nil, err
			}
			result.KubeConfigRaw = kubeCgfRaw
		case commonpb.CloudProvider_AZURE:
			values, err := state.GetValues(kubernetes_service.AzureEksCluster{}, resource.ResourceId)
			if err != nil {
				return nil, err
			}
			result.Endpoint = values["kube_config"].([]interface{})[0].(map[string]interface{})["host"].(string)
			result.CaCertificate = values["kube_config"].([]interface{})[0].(map[string]interface{})["cluster_ca_certificate"].(string)
			result.KubeConfigRaw = values["kube_config_raw"].(string)
		default:
			return nil, fmt.Errorf("unknown cloud %s", resource.GetCloud())
		}
	}

	result.DefaultNodePool, err = KubernetesNodePoolFromState(resource.DefaultNodePool, state)
	if err != nil {
		return nil, err
	}

	return result, err
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

func (r *KubernetesCluster) GetMainResourceName() (string, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return output.GetResourceName(kubernetes_service.AwsEksCluster{}), nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		return output.GetResourceName(kubernetes_service.AzureEksCluster{}), nil
	}
	return "", fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesCluster) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	if r.GetCloud() == commonpb.CloudProvider_AWS {
		var outputs []output.TfBlock
		defaultNodePoolResources, err := r.DefaultNodePool.Translate(ctx)
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

		role.GetFullResourceRef()

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
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		defaultPool, err := r.DefaultNodePool.translateAzNodePool()
		if err != nil {
			return nil, err
		}
		defaultPool.Name = defaultPool.AzResource.Name
		defaultPool.AzResource = nil
		defaultPool.ClusterId = ""

		return []output.TfBlock{
			&kubernetes_service.AzureEksCluster{
				AzResource:      common.NewAzResource(r.ResourceId, r.Args.Name, GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId), r.GetCloudSpecificLocation()),
				DefaultNodePool: defaultPool,
				DnsPrefix:       common.UniqueId(r.Args.Name, "aks", common.LowercaseAlphanumericFormatFunc),
				Identity:        kubernetes_service.AzureIdentity{Type: "SystemAssigned"},
				NetworkProfile: kubernetes_service.NetworkProfile{
					NetworkPlugin:    "azure",
					DnsServiceIp:     "10.100.0.10",
					DockerBridgeCidr: "172.17.0.1/16",
					ServiceCidr:      r.Args.ServiceCidr,
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *KubernetesCluster) getAwsSubnets() ([]subnet.AwsSubnet, []output.TfBlock, error) {
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

	vpcId, err := resources.GetMainOutputId(r.VirtualNetwork)
	if err != nil {
		return nil, nil, err
	}
	gtw, err := r.VirtualNetwork.GetAssociatedInternetGateway()
	if err != nil {
		return nil, nil, err
	}
	rt := route_table.AwsRouteTable{
		AwsResource: common.NewAwsResource(r.ResourceId+"_public_rt", r.Args.Name+"_public_rt"),
		VpcId:       vpcId,
		Routes: []route_table.AwsRouteTableRoute{
			{
				CidrBlock: "0.0.0.0/0",
				GatewayId: gtw,
			},
		},
	}

	subnet1 := subnet.AwsSubnet{
		AwsResource:      common.NewAwsResource(r.ResourceId+"_public_subnet", r.Args.Name+"_public_subnet"),
		CidrBlock:        subnetBlock1.String(),
		VpcId:            r.VirtualNetwork.GetVirtualNetworkId(),
		AvailabilityZone: az1,
	}
	subnet2 := subnet.AwsSubnet{
		AwsResource:      common.NewAwsResource(r.ResourceId+"_private_subnet", r.Args.Name+"_private_subnet"),
		CidrBlock:        subnetBlock2.String(),
		VpcId:            r.VirtualNetwork.GetVirtualNetworkId(),
		AvailabilityZone: az2,
	}

	rta := route_table_association.AwsRouteTableAssociation{
		AwsResource:  common.NewAwsResourceWithIdOnly(r.ResourceId + "_public_rta"),
		SubnetId:     fmt.Sprintf("%s.%s.id", output.GetResourceName(subnet1), subnet1.ResourceId),
		RouteTableId: fmt.Sprintf("%s.%s.id", output.GetResourceName(rt), rt.ResourceId),
	}

	return []subnet.AwsSubnet{subnet1, subnet2}, []output.TfBlock{subnet1, subnet2, rt, rta}, nil
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
