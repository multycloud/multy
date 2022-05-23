//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"os"
	"os/exec"
	"path"
	"testing"
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
	} `json:"status"`
	Metadata struct {
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func testKubernetes(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud)

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_2,
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
			Location:      commonpb.Location_EU_WEST_2,
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

func TestAwsKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AWS)
}
func TestAzureKubernetes(t *testing.T) {
	t.Parallel()
	testKubernetes(t, commonpb.CloudProvider_AZURE)
}
