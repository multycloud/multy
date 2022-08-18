//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func testPublicIp(t *testing.T, ctx context.Context, vm *resourcespb.VirtualMachineResource, config *ssh.ClientConfig) {
	location := vm.CommonParameters.Location
	cloud := vm.CommonParameters.CloudProvider
	if cloud == commonpb.CloudProvider_AWS || cloud == commonpb.CloudProvider_AZURE {
		t.Skip("public ip not yet implemented (https://github.com/multycloud/multy/issues/323)")
	}

	createPipRequest := &resourcespb.CreatePublicIpRequest{Resource: &resourcespb.PublicIpArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "vm-test-pip",
	}}
	pip, err := server.PublicIpService.Create(ctx, createPipRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create public ip: %+v", err)
	}
	cleanup(t, ctx, server.PublicIpService, pip)

	assert.Equal(t, createPipRequest.GetResource().GetCommonParameters().GetLocation(), pip.GetCommonParameters().GetLocation())
	assert.Equal(t, createPipRequest.GetResource().GetCommonParameters().GetCloudProvider(), pip.GetCommonParameters().GetCloudProvider())
	assert.Nil(t, pip.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createPipRequest.GetResource().GetName(), pip.GetName())

	updateReq := getNoopUpdate(vm)
	updateReq.Resource.GeneratePublicIp = false
	updateReq.Resource.PublicIpId = pip.CommonParameters.ResourceId
	_, err = server.VirtualMachineService.Update(ctx, updateReq)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to update vm: %+v", err)
	}
	t.Cleanup(func() {
		_, err = server.VirtualMachineService.Update(ctx, getNoopUpdate(vm))
		if err != nil {
			logGrpcErrorDetails(t, err)
			t.Fatalf("unable to update vm: %+v", err)
		}
	})

	// wait a bit so that the vm is reachable
	time.Sleep(10 * time.Second)
	testSSHConnection(t, pip.Ip, config)

}

func getNoopUpdate(vm *resourcespb.VirtualMachineResource) *resourcespb.UpdateVirtualMachineRequest {
	return &resourcespb.UpdateVirtualMachineRequest{
		ResourceId: vm.CommonParameters.ResourceId,
		Resource: &resourcespb.VirtualMachineArgs{
			CommonParameters: &commonpb.ResourceCommonArgs{
				ResourceGroupId: vm.CommonParameters.ResourceGroupId,
				Location:        vm.CommonParameters.Location,
				CloudProvider:   vm.CommonParameters.CloudProvider,
			},
			Name:                    vm.Name,
			NetworkInterfaceIds:     vm.NetworkInterfaceIds,
			NetworkSecurityGroupIds: vm.NetworkSecurityGroupIds,
			VmSize:                  vm.VmSize,
			UserDataBase64:          vm.UserDataBase64,
			SubnetId:                vm.SubnetId,
			PublicSshKey:            vm.PublicSshKey,
			PublicIpId:              vm.PublicIpId,
			GeneratePublicIp:        vm.GeneratePublicIp,
			ImageReference:          vm.ImageReference,
			AwsOverride:             vm.AwsOverride,
			AzureOverride:           vm.AzureOverride,
		},
	}
}
