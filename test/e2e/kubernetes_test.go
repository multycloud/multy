//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"
)

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
		PodIp      string                `json:"podIp,omitempty"`
		Phase      string                `json:"phase,omitempty"`
	} `json:"status"`
	Metadata struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func testKubernetes(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "k8s")
	region := commonpb.Location_EU_WEST_2
	if cloud == commonpb.CloudProvider_AWS {
		region = commonpb.Location_US_WEST_1
	}

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      region,
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
	cleanup(t, ctx, server.VnService, vn)

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
	cleanup(t, ctx, server.SubnetService, subnet)

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
	cleanup(t, ctx, server.RouteTableService, rt)

	createRtaRequest := &resourcespb.CreateRouteTableAssociationRequest{Resource: &resourcespb.RouteTableAssociationArgs{
		SubnetId:     subnet.CommonParameters.ResourceId,
		RouteTableId: rt.CommonParameters.ResourceId,
	}}
	rta, err := server.RouteTableAssociationService.Create(ctx, createRtaRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create route table association: %+v", err)
	}
	cleanup(t, ctx, server.RouteTableAssociationService, rta)

	createK8sClusterRequest := &resourcespb.CreateKubernetesClusterRequest{Resource: &resourcespb.KubernetesClusterArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      region,
			CloudProvider: cloud,
		},
		Name:             "k8testmulty",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		ServiceCidr:      "10.100.0.0/16",
		DefaultNodePool: &resourcespb.KubernetesNodePoolArgs{
			Name:              "default",
			SubnetId:          subnet.CommonParameters.ResourceId,
			StartingNodeCount: 1,
			MinNodeCount:      1,
			MaxNodeCount:      3,
			VmSize:            commonpb.VmSize_GENERAL_MEDIUM,
			DiskSizeGb:        30,
			AvailabilityZone:  []int32{2},
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
	cleanup(t, ctx, server.KubernetesClusterService, k8s)

	readAndAssertEq(t, server, ctx, createK8sClusterRequest, k8s.CommonParameters.ResourceId)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("cannot get home dir: %s", err)
	}
	kubecfg := path.Join(home, ".kube", fmt.Sprintf("config-%s", cloud.String()))
	err = os.MkdirAll(filepath.Dir(kubecfg), os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		t.Fatal(fmt.Errorf("can't create kube dir: %s", err.Error()))
	}
	// update kubectl configuration so that we can use kubectl commands - probably can't run this in parallel
	err = os.WriteFile(kubecfg, []byte(k8s.KubeConfigRaw), 0777)
	if err != nil {
		t.Fatal(fmt.Errorf("can't create kube config: %s", err.Error()))
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

	assert.Greater(t, len(o.Items), 0)
	for _, item := range o.Items {
		assert.Contains(t, item.Status.Conditions, NodeStatusCondition{
			Type:   "Ready",
			Status: "True",
		})
		labels := item.Metadata.Labels
		assert.Contains(t, maps.Keys(labels), "multy.dev/env")
		assert.Equal(t, "test", labels["multy.dev/env"])
	}

	testKubernetesDeployment(t, kubecfg)

}

func readAndAssertEq(t *testing.T, s *api.Server, ctx context.Context, cluster *resourcespb.CreateKubernetesClusterRequest, clusterId string) {
	read, err := s.KubernetesClusterService.Read(ctx, &resourcespb.ReadKubernetesClusterRequest{ResourceId: clusterId})
	if err != nil {
		t.Fatalf("unable to read k8s cluster, %s", err)
	}

	assert.Equal(t, cluster.GetResource().GetCommonParameters().GetLocation(), read.GetCommonParameters().GetLocation())
	assert.Equal(t, cluster.GetResource().GetCommonParameters().GetCloudProvider(), read.GetCommonParameters().GetCloudProvider())
	assert.Nil(t, read.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, cluster.GetResource().GetName(), read.GetName())
	assert.Equal(t, cluster.GetResource().GetServiceCidr(), read.GetServiceCidr())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetName(), read.GetDefaultNodePool().GetName())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetVmSize(), read.GetDefaultNodePool().GetVmSize())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetMinNodeCount(), read.GetDefaultNodePool().GetMinNodeCount())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetMaxNodeCount(), read.GetDefaultNodePool().GetMaxNodeCount())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetStartingNodeCount(), read.GetDefaultNodePool().GetStartingNodeCount())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetDiskSizeGb(), read.GetDefaultNodePool().GetDiskSizeGb())
	assert.EqualValues(t, cluster.GetResource().GetDefaultNodePool().GetLabels(), read.GetDefaultNodePool().GetLabels())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetSubnetId(), read.GetDefaultNodePool().GetSubnetId())
	assert.Equal(t, cluster.GetResource().GetDefaultNodePool().GetAvailabilityZone(), read.GetDefaultNodePool().GetAvailabilityZone())

}

func testKubernetesDeployment(t *testing.T, kubecfg string) {
	// kubectl create deployment test-deployment --image=nginx --replicas=2
	out, err := exec.Command("/usr/local/bin/kubectl", "--kubeconfig", kubecfg, "create", "deployment", "test-deployment", "--image", "nginx", "--replicas", "2").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	var o GetNodesOutput
	allRunning := false
	for i := 0; i < 10 && !allRunning; i++ {
		t.Logf("waiting 5 seconds to check if pods are running")
		time.Sleep(5 * time.Second)
		// kubectl get pods -o json
		out, err = exec.Command("/usr/local/bin/kubectl", "--kubeconfig", kubecfg, "get", "pods", "-o", "json").CombinedOutput()
		if err != nil {
			t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
		}

		o = GetNodesOutput{}
		err = json.Unmarshal(out, &o)
		if err != nil {
			t.Fatal(fmt.Errorf("output cant be parsed: %s", err))
		}

		assert.Len(t, o.Items, 2)
		runningNodes := 0
		for _, item := range o.Items {
			if item.Status.Phase != "Running" {
				break
			}
			runningNodes += 1
		}
		allRunning = runningNodes == len(o.Items)
	}

	if !allRunning {
		t.Fatal("pods are not running after 30 seconds")
	}

	err = fmt.Errorf("")
	for i := 0; i < 5 && err != nil; i++ {
		// kubectl exec test-deployment-76cdbc6456-6mdv2 -- curl 10.0.0.171
		out, err = exec.Command("/usr/local/bin/kubectl", "--kubeconfig", kubecfg, "exec", o.Items[0].Metadata.Name, "--", "curl", o.Items[1].Status.PodIp).CombinedOutput()
		if err != nil {
			t.Logf(string(out))
			t.Logf("waiting 5 seconds to check if error goes away")
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	assert.Contains(t, string(out), "<h1>Welcome to nginx!</h1>")
}

func TestAwsKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AWS)
}
func TestAzureKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AZURE)
}

func TestGcpKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_GCP)
}
