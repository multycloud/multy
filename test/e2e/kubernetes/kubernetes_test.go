//go:build e2e
// +build e2e

package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/service_context"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

const DestroyAfter = true

type GetNodesOutput struct {
	Items []Node `json:"items"`
}

type NodeStatusCondition struct {
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

type Node struct {
	Status struct {
		Conditions []NodeStatusCondition `json:"conditions,omitempty"`
	} `json:"status"`
	Metadata struct {
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

var server *api.Server

func init() {
	if server != nil {
		return
	}
	flags.Environment = flags.Local
	flags.DryRun = false
	awsClient, err := aws_client.NewClient()
	if err != nil {
		log.Fatalf("failed to initialize aws client: %v", err)
	}
	database, err := db.NewDatabase(awsClient)
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	serviceContext := &service_context.ServiceContext{
		Database:  database,
		AwsClient: awsClient,
	}

	s := api.CreateServer(serviceContext)
	server = &s
}

func testKubernetes(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud)

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_1,
			CloudProvider: cloud,
		},
		Name:      "k8-test-vn",
		CidrBlock: "10.0.0.0/16",
	}}
	vn, err := server.VnService.Create(ctx, createVnRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}

	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.VnService.Delete(ctx, &resourcespb.DeleteVirtualNetworkRequest{ResourceId: vn.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %s", err)
			}
		}
	})

	createSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "k8-test-subnet",
		CidrBlock:        "10.0.0.0/24",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
	}}
	subnet, err := server.SubnetService.Create(ctx, createSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create subnet: %+v", err)
	}

	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.SubnetService.Delete(ctx, &resourcespb.DeleteSubnetRequest{ResourceId: subnet.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createRtRequest := &resourcespb.CreateRouteTableRequest{Resource: &resourcespb.RouteTableArgs{
		Name:             "k8-test-rt",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		Routes: []*resourcespb.Route{
			{
				CidrBlock:   "0.0.0.0/0",
				Destination: resourcespb.RouteDestination_INTERNET,
			},
		},
	}}
	rt, err := server.RouteTableService.Create(ctx, createRtRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create route table: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.RouteTableService.Delete(ctx, &resourcespb.DeleteRouteTableRequest{ResourceId: rt.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createRtaRequest := &resourcespb.CreateRouteTableAssociationRequest{Resource: &resourcespb.RouteTableAssociationArgs{
		SubnetId:     subnet.CommonParameters.ResourceId,
		RouteTableId: rt.CommonParameters.ResourceId,
	}}
	rta, err := server.RouteTableAssociationService.Create(ctx, createRtaRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create route table association: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.RouteTableAssociationService.Delete(ctx, &resourcespb.DeleteRouteTableAssociationRequest{ResourceId: rta.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createK8sClusterRequest := &resourcespb.CreateKubernetesClusterRequest{Resource: &resourcespb.KubernetesClusterArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_1,
			CloudProvider: cloud,
		},
		Name:             "k8testmulty",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		DefaultNodePool: &resourcespb.KubernetesNodePoolArgs{
			Name:              "default",
			SubnetId:          subnet.CommonParameters.ResourceId,
			StartingNodeCount: 1,
			MinNodeCount:      1,
			MaxNodeCount:      3,
			VmSize:            commonpb.VmSize_MEDIUM,
			DiskSizeGb:        20,
			Labels: map[string]string{
				"multy.dev/env": "test",
			},
		},
	}}
	k8s, err := server.KubernetesClusterService.Create(ctx, createK8sClusterRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create kubernetes cluster: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.KubernetesClusterService.Delete(ctx, &resourcespb.DeleteKubernetesClusterRequest{ResourceId: k8s.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %s", err)
	}
	kubecfg := path.Join(home, ".kube", fmt.Sprintf("config-%s", cloud.String()))
	// update kubectl configuration so that we can use kubectl commands - probably can't run this in parallel
	if cloud == commonpb.CloudProvider_AWS {
		// aws eks --region eu-west-2 update-kubeconfig --name kubernetes_test
		out, err := exec.Command("aws", "eks", "--region", "eu-west-2", "update-kubeconfig", "--name", k8s.Name, "--kubeconfig", kubecfg).CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
		}
	} else if cloud == commonpb.CloudProvider_AZURE {
		out, err := exec.Command("/usr/bin/az", "login", "--service-principal", "-u", os.Getenv("ARM_CLIENT_ID"), "-p", os.Getenv("ARM_CLIENT_SECRET"), "--tenant", os.Getenv("ARM_TENANT_ID")).CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err, string(out)))
		}
		// az aks get-credentials --resource-group ks-rg --name example
		out, err = exec.Command("/usr/bin/az", "aks", "get-credentials", "--resource-group", k8s.CommonParameters.ResourceGroupId, "--name", k8s.Name, "--file", kubecfg).CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
		}
	} else {
		t.Fatalf("unknown cloud: %s", cloud)
	}
	// kubectl get nodes -o json
	out, err := exec.Command("/usr/local/bin/kubectl", "--kubeconfig", kubecfg, "get", "nodes", "-o", "json").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	o := GetNodesOutput{}

	err = json.Unmarshal(out, &o)
	if err != nil {
		t.Fatal(fmt.Errorf("output cant be parsed: %s", err))
	}

	assert.Len(t, o.Items, 1)
	assert.Contains(t, o.Items[0].Status.Conditions, NodeStatusCondition{
		Type:   "Ready",
		Status: "True",
	})

	labels := o.Items[0].Metadata.Labels
	assert.Contains(t, maps.Keys(labels), "multy.dev/env")
	assert.Equal(t, labels["multy.dev/env"], "test")
}

func getCtx(t *testing.T, cloud commonpb.CloudProvider) context.Context {
	accessKeyId, exists := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !exists {
		t.Fatalf("AWS_ACCESS_KEY_ID not set")
	}
	accessKeySecret, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !exists {
		t.Fatalf("AWS_SECRET_ACCESS_KEY not set")
	}
	azSubscriptionId, exists := os.LookupEnv("ARM_SUBSCRIPTION_ID")
	if !exists {
		t.Fatalf("ARM_SUBSCRIPTION_ID not set")
	}
	azClientId, exists := os.LookupEnv("ARM_CLIENT_ID")
	if !exists {
		t.Fatalf("ARM_CLIENT_ID not set")
	}
	azClientSecret, exists := os.LookupEnv("ARM_CLIENT_SECRET")
	if !exists {
		t.Fatalf("ARM_CLIENT_SECRET not set")
	}
	azTenantId, exists := os.LookupEnv("ARM_TENANT_ID")
	if !exists {
		t.Fatalf("ARM_TENANT_ID not set")
	}

	credentials := &credspb.CloudCredentials{
		AwsCreds: &credspb.AwsCredentials{
			AccessKey: accessKeyId,
			SecretKey: accessKeySecret,
		},
		AzureCreds: &credspb.AzureCredentials{
			SubscriptionId: azSubscriptionId,
			TenantId:       azTenantId,
			ClientId:       azClientId,
			ClientSecret:   azClientSecret,
		},
	}
	b, err := proto.Marshal(credentials)
	if err != nil {
		t.Fatalf("unable to marshal creds: %s", err)
	}

	md := map[string][]string{"api_key": {fmt.Sprintf("test-%s", cloud.String())}, "cloud-creds-bin": {string(b)}}
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return ctx
}

func logGrpcErrorDetails(t *testing.T, err error) {
	if s, ok := status.FromError(err); ok {
		for _, details := range s.Details() {
			if msg, ok := details.(interface{ GetErrorMessage() string }); ok {
				t.Logf("server returned error: %s", msg.GetErrorMessage())
			}
		}
	}
}

func TestAwsKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AWS)
}
func TestAzureKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AZURE)
}
