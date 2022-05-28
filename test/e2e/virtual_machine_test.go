//go:build e2e
// +build e2e

package e2e

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

// todo: azure enforces rule priority, aws doesnt
// todo: the provided RSA SSH key has 1024 bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure
func testVirtualMachine(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "vm")

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
			CloudProvider: cloud,
		},
		Name:      "vm-test-vn",
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

	createPublicSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "vm-test-public-subnet",
		CidrBlock:        "10.0.0.0/24",
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
	fmt.Println(rta)
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

	createNsgRequest := &resourcespb.CreateNetworkSecurityGroupRequest{Resource: &resourcespb.NetworkSecurityGroupArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
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
			Direction: resourcespb.Direction_BOTH_DIRECTIONS,
		}, {
			Protocol: "tcp",
			Priority: 120,
			PortRange: &resourcespb.PortRange{
				From: 80,
				To:   80,
			},
			CidrBlock: "0.0.0.0/0",
			Direction: resourcespb.Direction_BOTH_DIRECTIONS,
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

	pubKey, privKey, err := makeSSHKeyPair()
	fmt.Println(privKey)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create ssh key: %+v", err)
	}

	var vmSize commonpb.VmSize_Enum
	if cloud == commonpb.CloudProvider_AWS {
		vmSize = commonpb.VmSize_MICRO
	} else if cloud == commonpb.CloudProvider_AZURE {
		vmSize = commonpb.VmSize_LARGE
	}

	createVmRequest := &resourcespb.CreateVirtualMachineRequest{Resource: &resourcespb.VirtualMachineArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_EU_WEST_1,
			CloudProvider: cloud,
		},
		Name:                    "vm-test-vm",
		NetworkSecurityGroupIds: []string{nsg.CommonParameters.ResourceId},
		VmSize:                  vmSize,
		UserDataBase64: base64.StdEncoding.EncodeToString([]byte(`#!/bin/bash
sudo echo "hello world" > /tmp/test.txt`)),
		SubnetId:         publicSubnet.CommonParameters.ResourceId,
		PublicSshKey:     pubKey,
		GeneratePublicIp: true,
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

	signer, err := signerFromPem([]byte(privKey), []byte(""))
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

	fmt.Println(vm.GetPublicIp())

	time.Sleep(3 * time.Minute)

	// connect ot ssh server
	conn, err := ssh.Dial("tcp", vm.GetPublicIp()+":22", config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		t.Fatal(fmt.Errorf("error creating ssh session"))
	}
	defer session.Close()

	// run command and capture stdout/stderr
	output, err := session.CombinedOutput("sudo cat /tmp/test.txt")
	assert.Equal(t, string(output), "hello world\n")
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

func signerFromPem(pemBytes []byte, password []byte) (ssh.Signer, error) {

	// read pem block
	err := errors.New("Pem decode failed, no key found")
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, err
	}

	// handle encrypted key
	if x509.IsEncryptedPEMBlock(pemBlock) {
		// decrypt PEM
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("Decrypting PEM block failed %v", err)
		}

		// get RSA, EC or DSA key
		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		// generate signer instance from key
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("Creating signer from encrypted key failed %v", err)
		}

		return signer, nil
	} else {
		// generate signer instance from plain key
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing plain private key failed %v", err)
		}

		return signer, nil
	}
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing PKCS private key failed %v", err)
		} else {
			return key, nil
		}
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing EC private key failed %v", err)
		} else {
			return key, nil
		}
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing DSA private key failed %v", err)
		} else {
			return key, nil
		}
	default:
		return nil, fmt.Errorf("Parsing private key failed, unsupported key type %q", block.Type)
	}
}
