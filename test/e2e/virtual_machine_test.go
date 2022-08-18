//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"testing"
	"time"
)

// todo: azure enforces rule priority, aws doesnt
// todo: the provided RSA SSH key has 1024 bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure
func testVirtualMachine(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "vm")

	location := commonpb.Location_EU_WEST_1
	if cloud == commonpb.CloudProvider_AZURE {
		location = commonpb.Location_EU_WEST_2
	}

	pubKey, config := createSshConfig(t, cloud)
	vm, nsg := setUpVirtualMachine(t, ctx, location, cloud, pubKey)

	// wait a bit so that the vm is reachable
	time.Sleep(3 * time.Minute)

	t.Run("ssh_connection_success", func(t *testing.T) {
		testSSHConnection(t, vm.PublicIp, config)
	})

	t.Run("nsg_traffic_block", func(t *testing.T) {
		testNsgRules(t, ctx, nsg, vm.PublicIp, config)
	})

	t.Run("public_ip_id", func(t *testing.T) {
		testPublicIp(t, ctx, vm, config)
	})
}

func createSshConfig(t *testing.T, cloud commonpb.CloudProvider) (string, *ssh.ClientConfig) {
	pubKey, privKey, err := makeSSHKeyPair()
	if err != nil {
		t.Fatalf("unable to create ssh key: %+v", err)
	}

	//err = os.WriteFile("key", []byte(privKey), 0600)
	//if err != nil {
	//	t.Fatalf("unable to create ssh key: %+v", err)
	//}

	username := "adminuser"
	if cloud == commonpb.CloudProvider_AWS {
		username = "ubuntu"
	}

	signer, err := signerFromPem([]byte(privKey))
	if err != nil {
		t.Fatal(fmt.Errorf("error setting up cert"))
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return pubKey, config
}

func setUpVirtualMachine(t *testing.T, ctx context.Context,
	location commonpb.Location, cloud commonpb.CloudProvider,
	pubKey string) (*resourcespb.VirtualMachineResource, *resourcespb.NetworkSecurityGroupResource) {
	subnet, nsg := createNetworkWithInternetAccess(t, ctx, location, cloud, "vm")
	createVmRequest := &resourcespb.CreateVirtualMachineRequest{Resource: &resourcespb.VirtualMachineArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "vm-test-vm",
		// FIXME validate NetworkSecurityGroupIds not with NetworkInterfaceIds
		NetworkSecurityGroupIds: []string{nsg.CommonParameters.ResourceId},
		VmSize:                  commonpb.VmSize_GENERAL_MICRO,
		UserDataBase64: base64.StdEncoding.EncodeToString([]byte(`#!/bin/bash
sudo echo "hello world" > /tmp/test.txt`)),
		SubnetId:         subnet.CommonParameters.ResourceId,
		PublicSshKey:     pubKey,
		GeneratePublicIp: true,
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
	cleanup(t, ctx, server.VirtualMachineService, vm)

	readVm, err := server.VirtualMachineService.Read(ctx, &resourcespb.ReadVirtualMachineRequest{ResourceId: vm.CommonParameters.ResourceId})
	if err != nil {
		t.Fatalf("unable to read virtual machine, %s", err)
	}

	assert.Equal(t, readVm.GetCommonParameters().GetLocation(), createVmRequest.GetResource().GetCommonParameters().GetLocation())
	assert.Equal(t, readVm.GetCommonParameters().GetCloudProvider(), createVmRequest.GetResource().GetCommonParameters().GetCloudProvider())
	assert.Nil(t, readVm.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, readVm.GetName(), createVmRequest.GetResource().GetName())
	assert.Equal(t, readVm.GetVmSize(), createVmRequest.GetResource().GetVmSize())
	assert.Equal(t, readVm.GetPublicIpId(), createVmRequest.GetResource().GetPublicIpId())
	assert.Equal(t, readVm.GetSubnetId(), createVmRequest.GetResource().GetSubnetId())
	//assert.Equal(t, readVm.GetAvailabilityZone(), createVmRequest.GetResource().GetAvailabilityZone())
	assert.Equal(t, readVm.GetGeneratePublicIp(), createVmRequest.GetResource().GetGeneratePublicIp())
	assert.Equal(t, readVm.GetPublicSshKey(), createVmRequest.GetResource().GetPublicSshKey())
	assert.Equal(t, readVm.GetUserDataBase64(), createVmRequest.GetResource().GetUserDataBase64())
	assert.Equal(t, readVm.GetNetworkInterfaceIds(), createVmRequest.GetResource().GetNetworkInterfaceIds())
	assert.Equal(t, readVm.GetNetworkSecurityGroupIds(), createVmRequest.GetResource().GetNetworkSecurityGroupIds())
	assert.Equal(t, readVm.GetImageReference().GetOs(), createVmRequest.GetResource().GetImageReference().GetOs())
	assert.Equal(t, readVm.GetImageReference().GetVersion(), createVmRequest.GetResource().GetImageReference().GetVersion())

	return vm, nsg
}

func testSSHConnection(t *testing.T, pip string, config *ssh.ClientConfig) {
	conn, err := ssh.Dial("tcp", pip+":22", config)
	if err != nil {
		t.Fatal(fmt.Errorf("error in ssh connection: %+v", err))
	}
	t.Cleanup(func() {
		conn.Close()
	})

	session, err := conn.NewSession()
	if err != nil {
		t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
	}
	t.Cleanup(func() {
		session.Close()
	})

	// run command and capture stdout/stderr
	output, err := session.CombinedOutput("sudo cat /tmp/test.txt")
	if err != nil {
		t.Fatal(fmt.Errorf("error running command: %+v", err))
	}

	assert.Equal(t, "hello world\n", string(output), config)
}

func TestAwsVirtualMachine(t *testing.T) {
	t.Parallel()
	testVirtualMachine(t, commonpb.CloudProvider_AWS)
}
func TestAzureVirtualMachine(t *testing.T) {
	t.Parallel()
	testVirtualMachine(t, commonpb.CloudProvider_AZURE)
}
func TestGcpVirtualMachine(t *testing.T) {
	t.Parallel()
	testVirtualMachine(t, commonpb.CloudProvider_GCP)
}

func makeSSHKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// generate and write private key as PEM
	var privKeyBuf strings.Builder

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

	return pubKeyBuf.String(), privKeyBuf.String(), nil
}

func signerFromPem(pemBytes []byte) (ssh.Signer, error) {
	// read pem block
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, fmt.Errorf("Pem decode failed, no key found")
	}

	// generate signer instance from plain key
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Parsing plain private key failed %v", err)
	}

	return signer, nil
}
