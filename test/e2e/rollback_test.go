//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestRollbackResourceGroup(t *testing.T) {
	ctx := getCtx(t, commonpb.CloudProvider_AZURE, "rollback")

	rgDeleted := false
	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_2,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:      "vn-rollback-test",
		CidrBlock: "10.0.0.0/16",
	}}
	vn, err := server.VnService.Create(ctx, createVnRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter && !rgDeleted {
			_, _ = server.VnService.Delete(ctx, &resourcespb.DeleteVirtualNetworkRequest{ResourceId: vn.CommonParameters.ResourceId})
		}
	})

	createSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "subnet-rollback-test",
		CidrBlock:        "10.0.0.0/16",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
	}}
	subnet, err := server.SubnetService.Create(ctx, createSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create subnet: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter && !rgDeleted {
			_, _ = server.SubnetService.Delete(ctx, &resourcespb.DeleteSubnetRequest{ResourceId: subnet.CommonParameters.ResourceId})
		}
	})

	out, err := exec.Command("az", "login", "--service-principal", "-u", os.Getenv("ARM_CLIENT_ID"), "-p", os.Getenv("ARM_CLIENT_SECRET"), "--tenant", os.Getenv("ARM_TENANT_ID")).CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err, string(out)))
	}
	// delete previously created resources
	out, err = exec.Command("az", "group", "delete", "-n", vn.CommonParameters.ResourceGroupId, "--yes").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}
	rgDeleted = true

	// expected resource group id: "nic-xxxx-rg"
	createNetworkInterfaceRequest := &resourcespb.CreateNetworkInterfaceRequest{Resource: &resourcespb.NetworkInterfaceArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_2,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:     "ni-rollback-test",
		SubnetId: subnet.CommonParameters.ResourceId,
	}}
	ni, err := server.NetworkInterfaceService.Create(ctx, createNetworkInterfaceRequest)
	if err == nil {
		t.Cleanup(func() {
			if *destroyAfter {
				_, err := server.NetworkInterfaceService.Delete(ctx, &resourcespb.DeleteNetworkInterfaceRequest{ResourceId: ni.CommonParameters.ResourceId})
				if err != nil {
					logGrpcErrorDetails(t, err)
					t.Logf("unable to delete resource %s", err)
				}
			}
		})
		t.Fatalf("network interface should have been unable to be creatde after VN and subnet was deleted")
	}

	matches := regexp.MustCompile("\\w+-(\\w+)-rg").FindStringSubmatch(vn.CommonParameters.ResourceGroupId)
	if len(matches) < 2 {
		t.Fatalf("resource group '%s' doesn't follow the expected format", vn.CommonParameters.ResourceGroupId)
	}

	out, err = exec.Command("az", "group", "exists", "--name", fmt.Sprintf("nic-%s-rg", matches[1])).CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	assert.Equal(t, "false\n", string(out))
}
