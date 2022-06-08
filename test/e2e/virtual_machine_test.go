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

const (
	vnCidrBlock      = "10.0.0.0/16"
	publicSubnetCidr = "10.0.0.0/24"
)

// todo: azure enforces rule priority, aws doesnt
// todo: the provided RSA SSH key has 1024 bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure
func testVirtualMachine(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "vm")

	location := commonpb.Location_EU_WEST_1
	if cloud == commonpb.CloudProvider_AZURE {
		location = commonpb.Location_EU_WEST_2
	}

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:      "vm-test-vn",
		CidrBlock: vnCidrBlock,
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

	createPublicSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "vm-test-public-subnet",
		CidrBlock:        publicSubnetCidr,
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		AvailabilityZone: 1,
	}}
	publicSubnet, err := server.SubnetService.Create(ctx, createPublicSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create publicSubnet: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.SubnetService.Delete(ctx, &resourcespb.DeleteSubnetRequest{ResourceId: publicSubnet.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createRtRequest := &resourcespb.CreateRouteTableRequest{Resource: &resourcespb.RouteTableArgs{
		Name:             "vm-test-rt",
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
		SubnetId:     publicSubnet.CommonParameters.ResourceId,
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

	createPipRequest := &resourcespb.CreatePublicIpRequest{Resource: &resourcespb.PublicIpArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "pip-test-vm",
	}}
	pip, err := server.PublicIpService.Create(ctx, createPipRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create public ip: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.PublicIpService.Delete(ctx, &resourcespb.DeletePublicIpRequest{ResourceId: pip.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	// FIXME "network_interface": conflicts with vpc_security_group_ids
	// NSG might need to be associated via NIC
	createNicRequest := &resourcespb.CreateNetworkInterfaceRequest{Resource: &resourcespb.NetworkInterfaceArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:       "nsg-test-vm",
		SubnetId:   publicSubnet.CommonParameters.ResourceId,
		PublicIpId: pip.CommonParameters.ResourceId,
	}}
	nic, err := server.NetworkInterfaceService.Create(ctx, createNicRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create network interface: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.NetworkInterfaceService.Delete(ctx, &resourcespb.DeleteNetworkInterfaceRequest{ResourceId: nic.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	createNsgRequest := &resourcespb.CreateNetworkSecurityGroupRequest{Resource: &resourcespb.NetworkSecurityGroupArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:             "nsg-test-vm",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		Rules: []*resourcespb.NetworkSecurityRule{{
			Protocol: "tcp",
			Priority: 100,
			PortRange: &resourcespb.PortRange{
				From: 22,
				To:   22,
			},
			CidrBlock: "0.0.0.0/0",
			Direction: resourcespb.Direction_BOTH_DIRECTIONS,
		}, {
			Protocol: "tcp",
			Priority: 110,
			PortRange: &resourcespb.PortRange{
				From: 443,
				To:   443,
			},
			CidrBlock: "0.0.0.0/0",
			Direction: resourcespb.Direction_EGRESS,
		}, {
			Protocol: "tcp",
			Priority: 120,
			PortRange: &resourcespb.PortRange{
				From: 80,
				To:   80,
			},
			CidrBlock: "0.0.0.0/0",
			Direction: resourcespb.Direction_EGRESS,
		}},
	}}
	nsg, err := server.NetworkSecurityGroupService.Create(ctx, createNsgRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create network security group: %+v", err)
	}
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.NetworkSecurityGroupService.Delete(ctx, &resourcespb.DeleteNetworkSecurityGroupRequest{ResourceId: nsg.CommonParameters.ResourceId})
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
		if DestroyAfter {
			_, err := server.NetworkInterfaceSecurityGroupAssociationService.Delete(ctx, &resourcespb.DeleteNetworkInterfaceSecurityGroupAssociationRequest{ResourceId: nicNsgAssociation.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	pubKey, privKey, err := makeSSHKeyPair()
	if err != nil {
		t.Fatalf("unable to create ssh key: %+v", err)
	}

	createVmRequest := &resourcespb.CreateVirtualMachineRequest{Resource: &resourcespb.VirtualMachineArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "vm-test-vm",
		// FIXME validate NetworkSecurityGroupIds not with NetworkInterfaceIds
		//NetworkSecurityGroupIds: []string{nsg.CommonParameters.ResourceId},
		VmSize: commonpb.VmSize_GENERAL_NANO,
		UserDataBase64: base64.StdEncoding.EncodeToString([]byte(`#!/bin/bash
sudo echo "hello world" > /tmp/test.txt`)),
		SubnetId:         publicSubnet.CommonParameters.ResourceId,
		PublicSshKey:     pubKey,
		GeneratePublicIp: false,
		// FIXME this does nothing
		//PublicIpId:          pip.CommonParameters.ResourceId,
		NetworkInterfaceIds: []string{nic.CommonParameters.ResourceId},
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
		if DestroyAfter {
			_, err := server.VirtualMachineService.Delete(ctx, &resourcespb.DeleteVirtualMachineRequest{ResourceId: vm.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	var username string
	if cloud == commonpb.CloudProvider_AWS {
		username = "ubuntu"
	} else if cloud == commonpb.CloudProvider_AZURE {
		username = "adminuser"
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

	time.Sleep(3 * time.Minute)

	conn, err := ssh.Dial("tcp", pip.GetIp()+":22", config)
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
	updateNsgRules(t, ctx, nsg, pip.GetIp(), config)
}

func updateNsgRules(t *testing.T, ctx context.Context, nsg *resourcespb.NetworkSecurityGroupResource, ip string, config *ssh.ClientConfig) {
	nsg.Rules = []*resourcespb.NetworkSecurityRule{{
		Protocol: "tcp",
		Priority: 110,
		PortRange: &resourcespb.PortRange{
			From: 443,
			To:   443,
		},
		CidrBlock: "0.0.0.0/0",
		Direction: resourcespb.Direction_EGRESS,
	}, {
		Protocol: "tcp",
		Priority: 120,
		PortRange: &resourcespb.PortRange{
			From: 80,
			To:   80,
		},
		CidrBlock: "0.0.0.0/0",
		Direction: resourcespb.Direction_EGRESS,
	}}
	_, err := server.NetworkSecurityGroupService.Update(ctx, &resourcespb.UpdateNetworkSecurityGroupRequest{
		ResourceId: nsg.CommonParameters.ResourceId,
		Resource: &resourcespb.NetworkSecurityGroupArgs{
			CommonParameters: &commonpb.ResourceCommonArgs{
				Location:        nsg.CommonParameters.Location,
				CloudProvider:   nsg.CommonParameters.CloudProvider,
				ResourceGroupId: nsg.CommonParameters.ResourceGroupId,
			},
			Name:             nsg.Name,
			VirtualNetworkId: nsg.VirtualNetworkId,
			Rules:            nsg.Rules,
		},
	})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatal(fmt.Errorf("error updating nsg: %+v", err))
	}

	time.Sleep(20 * time.Second)

	conn, err := ssh.Dial("tcp", ip+":22", config)
	assert.Error(t, err)
	if err == nil {
		t.Cleanup(func() {
			conn.Close()
		})
	}
}

func TestAwsVirtualMachine(t *testing.T) {
	t.Parallel()
	testVirtualMachine(t, commonpb.CloudProvider_AWS)
}
func TestAzureVirtualMachine(t *testing.T) {
	t.Parallel()
	testVirtualMachine(t, commonpb.CloudProvider_AZURE)
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
