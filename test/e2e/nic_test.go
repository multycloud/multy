//go:build e2e
// +build e2e

package e2e

import (
	"encoding/base64"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"testing"
	"time"
)

func testNetworkInterface(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "ni")
	pubKey, config := createSshConfig(t, cloud)

	location := commonpb.Location_EU_WEST_1
	if cloud == commonpb.CloudProvider_AZURE {
		location = commonpb.Location_EU_WEST_2
	}

	subnet, nsg := createNetworkWithInternetAccess(t, ctx, location, cloud, "nsg")
	createPipRequest := &resourcespb.CreatePublicIpRequest{Resource: &resourcespb.PublicIpArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "ni-test-pip",
	}}
	pip, err := server.PublicIpService.Create(ctx, createPipRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create public ip: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter {
			_, err := server.PublicIpService.Delete(ctx, &resourcespb.DeletePublicIpRequest{ResourceId: pip.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})
	createNicRequest := &resourcespb.CreateNetworkInterfaceRequest{Resource: &resourcespb.NetworkInterfaceArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:       "ni-test-nsg",
		SubnetId:   subnet.CommonParameters.ResourceId,
		PublicIpId: pip.CommonParameters.ResourceId,
	}}
	nic, err := server.NetworkInterfaceService.Create(ctx, createNicRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create network interface: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter {
			_, err := server.NetworkInterfaceService.Delete(ctx, &resourcespb.DeleteNetworkInterfaceRequest{ResourceId: nic.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})
	createNicNsgRequest := &resourcespb.CreateNetworkInterfaceSecurityGroupAssociationRequest{Resource: &resourcespb.NetworkInterfaceSecurityGroupAssociationArgs{
		SecurityGroupId:    nsg.CommonParameters.ResourceId,
		NetworkInterfaceId: nic.CommonParameters.ResourceId,
	}}
	nicNsgAssociation, err := server.NetworkInterfaceSecurityGroupAssociationService.Create(ctx, createNicNsgRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create network interface nsg association: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter {
			_, err := server.NetworkInterfaceSecurityGroupAssociationService.Delete(ctx, &resourcespb.DeleteNetworkInterfaceSecurityGroupAssociationRequest{ResourceId: nicNsgAssociation.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createVmRequest := &resourcespb.CreateVirtualMachineRequest{Resource: &resourcespb.VirtualMachineArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:   "ni-test-vm",
		VmSize: commonpb.VmSize_GENERAL_MICRO,
		UserDataBase64: base64.StdEncoding.EncodeToString([]byte(`#!/bin/bash
sudo echo "hello world" > /tmp/test.txt`)),
		SubnetId:            subnet.CommonParameters.ResourceId,
		PublicSshKey:        pubKey,
		GeneratePublicIp:    false,
		NetworkInterfaceIds: []string{nic.CommonParameters.ResourceId},
		// FIXME this does nothing
		//PublicIpId:          pip.CommonParameters.ResourceId,
		ImageReference: &resourcespb.ImageReference{
			Os:      resourcespb.ImageReference_UBUNTU,
			Version: "18.04",
		},
	}}
	vm, err := server.VirtualMachineService.Create(ctx, createVmRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create virtual machine: %+v", err)
	}
	t.Cleanup(func() {
		if *destroyAfter {
			_, err := server.VirtualMachineService.Delete(ctx, &resourcespb.DeleteVirtualMachineRequest{ResourceId: vm.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	// wait a bit so that the vm is reachable
	time.Sleep(3 * time.Minute)
	testSSHConnection(t, pip.Ip, config)
}

func TestAwsNetworkInterface(t *testing.T) {
	t.Parallel()
	testNetworkInterface(t, commonpb.CloudProvider_AWS)
}
func TestAzureNetworkInterface(t *testing.T) {
	t.Parallel()
	testNetworkInterface(t, commonpb.CloudProvider_AZURE)
}
