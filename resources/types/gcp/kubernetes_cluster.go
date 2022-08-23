package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/iam"
	"github.com/multycloud/multy/resources/output/kubernetes_node_pool"
	"github.com/multycloud/multy/resources/output/kubernetes_service"
	"github.com/multycloud/multy/resources/output/subnet"
	"github.com/multycloud/multy/resources/output/virtual_network"
	"github.com/multycloud/multy/resources/types"
	"gopkg.in/yaml.v3"
)

type GcpKubernetesCluster struct {
	*types.KubernetesCluster
}

func InitKubernetesCluster(r *types.KubernetesCluster) resources.ResourceTranslator[*resourcespb.KubernetesClusterResource] {
	return GcpKubernetesCluster{r}
}

func (r GcpKubernetesCluster) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.KubernetesClusterResource, error) {
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
		Endpoint:         "dryrun",
	}
	var err error
	result.DefaultNodePool, err = GcpKubernetesNodePool{r.DefaultNodePool}.FromState(state, plan)
	if err != nil {
		return nil, err
	}
	if flags.DryRun {
		return result, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if cluster, exists, err := output.MaybeGetParsedById[kubernetes_service.GoogleContainerCluster](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		result.Name = cluster.Name
		if len(cluster.IpAllocationPolicy) == 0 {
			result.ServiceCidr = ""
		} else {
			result.ServiceCidr = cluster.IpAllocationPolicy[0].ServicesIpv4CidrBlock
		}
		result.Endpoint = fmt.Sprintf("https://%s", cluster.Endpoint)
		result.CaCertificate = cluster.MasterAuth[0].ClusterCaCertificate

		rawConfig, err := createKubeConfig(r.Args.Name, result.CaCertificate, result.Endpoint)
		if err != nil {
			return nil, err
		}

		result.KubeConfigRaw = rawConfig

		result.GcpOutputs = &resourcespb.KubernetesClusterGcpOutputs{
			GkeClusterId: cluster.SelfLink,
		}
		output.AddToStatuses(statuses, "gcp_container_cluster", output.MaybeGetPlannedChageById[kubernetes_service.GoogleContainerCluster](plan, r.ResourceId))
	} else {
		statuses["gcp_container_cluster"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if stateResource, exists, err := output.MaybeGetParsedById[iam.GoogleServiceAccount](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		result.GcpOutputs.ServiceAccountEmail = stateResource.Email
		output.AddToStatuses(statuses, "gcp_service_account", output.MaybeGetPlannedChageById[iam.GoogleServiceAccount](plan, r.ResourceId))
	} else {
		statuses["gcp_service_account"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		result.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return result, nil
}

func createKubeConfig(clusterName string, certData string, endpoint string) (string, error) {
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
						ApiVersion:         "client.authentication.k8s.io/v1beta1",
						Command:            "gke-gcloud-auth-plugin",
						InteractiveMode:    "IfAvailable",
						ProvideClusterInfo: true,
						InstallHint:        "Install gke-gcloud-auth-plugin for use with kubectl by following\n        https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke",
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

func (r GcpKubernetesCluster) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	var outputs []output.TfBlock

	serviceAccountId := r.getServiceAccountId()
	serviceAccount := &iam.GoogleServiceAccount{
		GcpResource: common.NewGcpResource(r.ResourceId, "", r.Args.GetGcpOverride().GetProject()),
		AccountId:   serviceAccountId,
		DisplayName: fmt.Sprintf("Service Account for cluster %s - created by Multy", r.Args.Name),
	}
	outputs = append(outputs, serviceAccount)

	defaultNodePoolResources, err := GcpKubernetesNodePool{r.DefaultNodePool}.Translate(ctx)
	if err != nil {
		return nil, err
	}
	outputs = append(outputs, defaultNodePoolResources...)

	var tags []string
	tags = append(tags, GcpVirtualNetwork{r.DefaultNodePool.Subnet.VirtualNetwork}.getVnTag())
	tags = append(tags, GcpSubnet{r.DefaultNodePool.Subnet}.getNetworkTags()...)
	outputs = append(outputs, &kubernetes_service.GoogleContainerCluster{
		GcpResource:           common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
		RemoveDefaultNodePool: true,
		InitialNodeCount:      1,
		Location:              r.GetCloudSpecificLocation(),
		Subnetwork:            fmt.Sprintf("%s.%s.id", output.GetResourceName(subnet.GoogleComputeSubnetwork{}), r.DefaultNodePool.Subnet.ResourceId),
		Network:               fmt.Sprintf("%s.%s.id", output.GetResourceName(virtual_network.GoogleComputeNetwork{}), r.VirtualNetwork.ResourceId),
		IpAllocationPolicy: []kubernetes_service.GoogleContainerClusterIpAllocationPolicy{{
			ServicesIpv4CidrBlock: r.Args.ServiceCidr,
		}},
		NodeConfig: []kubernetes_node_pool.GoogleContainerNodeConfig{{
			MachineType: "e2-micro",
			Tags:        tags,
			// Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
			ServiceAccount: fmt.Sprintf("%s.%s.email", output.GetResourceName(iam.GoogleServiceAccount{}), r.ResourceId),
			OAuthScopes:    []string{"https://www.googleapis.com/auth/cloud-platform"},
		}},
	})
	return outputs, nil
}

func (r GcpKubernetesCluster) getServiceAccountId() string {
	return common.UniqueId(fmt.Sprintf("%s-%s", r.Args.Name, r.ResourceId), "-sa-", common.LowercaseAlphanumericAndDashFormatFunc)
}

func (r GcpKubernetesCluster) GetMainResourceName() (string, error) {
	return output.GetResourceName(kubernetes_service.GoogleContainerCluster{}), nil
}
