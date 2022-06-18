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
	"path/filepath"
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
	ctx := getCtx(t, cloud, "k8s")
	region := commonpb.Location_US_WEST_1
	if cloud == commonpb.CloudProvider_AZURE {
		region = commonpb.Location_EU_WEST_2
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
		DefaultNodePool: &resourcespb.KubernetesNodePoolArgs{
			Name:              "default",
			SubnetId:          subnet.CommonParameters.ResourceId,
			StartingNodeCount: 1,
			MinNodeCount:      1,
			MaxNodeCount:      3,
			VmSize:            commonpb.VmSize_GENERAL_MEDIUM,
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
	cleanup(t, ctx, server.KubernetesClusterService, k8s)

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
