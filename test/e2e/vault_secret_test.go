//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func testVaultSecret(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "vault")

	location := commonpb.Location_EU_WEST_1
	if cloud == commonpb.CloudProvider_AZURE {
		location = commonpb.Location_EU_WEST_2
	}

	pubKey, config := createSshConfig(t, cloud)
	vm := setupVmForVaultTest(t, ctx, location, cloud, pubKey)

	createVaultRequest := &resourcespb.CreateVaultRequest{Resource: &resourcespb.VaultArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "test-vault-multy",
	}}
	vault, err := server.VaultService.Create(ctx, createVaultRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vault: %+v", err)
	}
	cleanup(t, ctx, server.VaultService, vault)

	readVault, err := server.VaultService.Read(ctx, &resourcespb.ReadVaultRequest{ResourceId: vault.CommonParameters.ResourceId})
	if err != nil {
		t.Fatalf("unable to read vault, %s", err)
	}

	assert.Equal(t, createVaultRequest.GetResource().GetCommonParameters().GetLocation(), readVault.GetCommonParameters().GetLocation())
	assert.Equal(t, createVaultRequest.GetResource().GetCommonParameters().GetCloudProvider(), readVault.GetCommonParameters().GetCloudProvider())
	assert.Nil(t, readVault.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createVaultRequest.GetResource().GetName(), readVault.GetName())

	createVaultSecretRequest := &resourcespb.CreateVaultSecretRequest{Resource: &resourcespb.VaultSecretArgs{
		Name:    "test-secret",
		Value:   "test-value",
		VaultId: vault.CommonParameters.ResourceId,
	}}
	vaultSecret, err := server.VaultSecretService.Create(ctx, createVaultSecretRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vault secret: %+v", err)
	}
	cleanup(t, ctx, server.VaultSecretService, vaultSecret)

	readVaultSecret, err := server.VaultSecretService.Read(ctx, &resourcespb.ReadVaultSecretRequest{ResourceId: vaultSecret.CommonParameters.ResourceId})
	if err != nil {
		t.Fatalf("unable to read vault secret, %s", err)
	}

	assert.Nil(t, readVaultSecret.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createVaultSecretRequest.GetResource().GetName(), readVaultSecret.GetName())
	assert.Equal(t, createVaultSecretRequest.GetResource().GetValue(), readVaultSecret.GetValue())
	assert.Equal(t, createVaultSecretRequest.GetResource().GetVaultId(), readVaultSecret.GetVaultId())

	t.Run("vault-access-policy-reader", func(t *testing.T) {
		testReaderAccess(t, ctx, cloud, vm, config, vault, vaultSecret)
	})

	t.Run("vault-access-policy-writer", func(t *testing.T) {
		testWriterAccess(t, ctx, cloud, vm, config, vault, vaultSecret)
	})

	// TODO: add test for OWNER policy
}

func testReaderAccess(t *testing.T, ctx context.Context, cloud commonpb.CloudProvider,
	vm *resourcespb.VirtualMachineResource, config *ssh.ClientConfig, vault *resourcespb.VaultResource,
	vaultSecret *resourcespb.VaultSecretResource) {
	createVaultAccessRequest := &resourcespb.CreateVaultAccessPolicyRequest{Resource: &resourcespb.VaultAccessPolicyArgs{
		VaultId:  vault.CommonParameters.ResourceId,
		Identity: vm.IdentityId,
		Access:   resourcespb.VaultAccess_READ,
	}}
	vap, err := server.VaultAccessPolicyService.Create(ctx, createVaultAccessRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vap: %+v", err)
	}
	cleanup(t, ctx, server.VaultAccessPolicyService, vap)

	// wait a bit so that the vm is reachable and policy has propagated
	time.Sleep(3 * time.Minute)

	conn, err := ssh.Dial("tcp", vm.PublicIp+":22", config)
	if err != nil {
		t.Fatal(fmt.Errorf("error in ssh connection: %+v", err))
	}
	t.Cleanup(func() {
		conn.Close()
	})

	err = fmt.Errorf("")
	var output []byte
	for i := 0; i < 6 && err != nil; i++ {
		var session *ssh.Session
		session, err = conn.NewSession()
		if err != nil {
			t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
		}
		output, err = session.CombinedOutput(getSecretReadCommand(vault.Name, vaultSecret.Name, cloud))
		if err != nil {
			t.Logf("command outputted: %s. waiting 1 min and trying again", output)
			time.Sleep(1 * time.Minute)
		}
		session.Close()
	}

	if err != nil {
		t.Fatal(fmt.Errorf("error running command: %+v", err))
	}

	assert.Equal(t, "test-value\n", string(output))

	session, err := conn.NewSession()
	if err != nil {
		t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
	}
	_, err = session.CombinedOutput(getSecretWriteCommand(vault.Name, vaultSecret.Name, cloud))
	assert.Error(t, err, "should not be able to write with read only perms")

	// remove the access policy and verify that eventually an error is returned when accessing the secret
	_, err = server.VaultAccessPolicyService.Delete(ctx, &resourcespb.DeleteVaultAccessPolicyRequest{ResourceId: vap.CommonParameters.ResourceId})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to delete vap: %+v", err)
	}

	err = nil
	for i := 0; i < 6 && err == nil; i++ {
		var session *ssh.Session
		session, err = conn.NewSession()
		if err != nil {
			t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
		}
		output, err = session.CombinedOutput(getSecretReadCommand(vault.Name, vaultSecret.Name, cloud))
		if err == nil {
			t.Logf("command outputted: %s. waiting 1 min and trying again", output)
			time.Sleep(1 * time.Minute)
		}
		session.Close()
	}
	assert.Error(t, err)
	t.Logf("comamnd returned %s", output)
}

func testWriterAccess(t *testing.T, ctx context.Context, cloud commonpb.CloudProvider,
	vm *resourcespb.VirtualMachineResource, config *ssh.ClientConfig, vault *resourcespb.VaultResource,
	vaultSecret *resourcespb.VaultSecretResource) {
	createVaultAccessRequest := &resourcespb.CreateVaultAccessPolicyRequest{Resource: &resourcespb.VaultAccessPolicyArgs{
		VaultId:  vault.CommonParameters.ResourceId,
		Identity: vm.IdentityId,
		Access:   resourcespb.VaultAccess_WRITE,
	}}
	vap, err := server.VaultAccessPolicyService.Create(ctx, createVaultAccessRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vap: %+v", err)
	}
	cleanup(t, ctx, server.VaultAccessPolicyService, vap)

	// wait a bit so that the vm is reachable and policy has propagated
	time.Sleep(3 * time.Minute)

	conn, err := ssh.Dial("tcp", vm.PublicIp+":22", config)
	if err != nil {
		t.Fatal(fmt.Errorf("error in ssh connection: %+v", err))
	}
	t.Cleanup(func() {
		conn.Close()
	})

	err = fmt.Errorf("")
	var output []byte
	for i := 0; i < 6 && err != nil; i++ {
		var session *ssh.Session
		session, err = conn.NewSession()
		if err != nil {
			t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
		}
		output, err = session.CombinedOutput(getSecretWriteCommand(vault.Name, vaultSecret.Name, cloud))
		if err != nil {
			t.Logf("command outputted: %s. waiting 1 min and trying again", output)
			time.Sleep(1 * time.Minute)
		}
		session.Close()
	}

	if err != nil {
		t.Fatal(fmt.Errorf("error running command: %+v", err))
	}

	session, err := conn.NewSession()
	if err != nil {
		t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
	}
	_, err = session.CombinedOutput(getSecretReadCommand(vault.Name, vaultSecret.Name, cloud))
	assert.Error(t, err, "should not be able to read with write only perms")

	// remove the access policy and verify that eventually an error is returned when accessing the secret
	_, err = server.VaultAccessPolicyService.Delete(ctx, &resourcespb.DeleteVaultAccessPolicyRequest{ResourceId: vap.CommonParameters.ResourceId})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to delete vap: %+v", err)
	}

	err = nil
	for i := 0; i < 6 && err == nil; i++ {
		var session *ssh.Session
		session, err = conn.NewSession()
		if err != nil {
			t.Fatal(fmt.Errorf("error creating ssh session: %+v", err))
		}
		output, err = session.CombinedOutput(getSecretWriteCommand(vault.Name, vaultSecret.Name, cloud))
		if err == nil {
			t.Logf("command outputted: %s. waiting 1 min and trying again", output)
			time.Sleep(1 * time.Minute)
		}
		session.Close()
	}
	assert.Error(t, err)
	t.Logf("comamnd returned %s", output)
}

func setupVmForVaultTest(t *testing.T, ctx context.Context, location commonpb.Location, cloud commonpb.CloudProvider, pubKey string) *resourcespb.VirtualMachineResource {
	subnet, nsg := createNetworkWithInternetAccess(t, ctx, location, cloud, "vault")
	createVmRequest := &resourcespb.CreateVirtualMachineRequest{Resource: &resourcespb.VirtualMachineArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name: "vault-test-vm",
		// FIXME validate NetworkSecurityGroupIds not with NetworkInterfaceIds
		NetworkSecurityGroupIds: []string{nsg.CommonParameters.ResourceId},
		VmSize:                  commonpb.VmSize_GENERAL_MICRO,
		UserDataBase64:          base64.StdEncoding.EncodeToString([]byte(getInitScript(cloud))),
		SubnetId:                subnet.CommonParameters.ResourceId,
		PublicSshKey:            pubKey,
		GeneratePublicIp:        true,
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
	return vm
}

func getInitScript(cloud commonpb.CloudProvider) string {
	switch cloud {
	case commonpb.CloudProvider_AWS:
		return `#!/bin/bash
sudo apt -y update
sudo apt -y install apt-transport-https ca-certificates curl
sudo apt -y install awscli
`
	case commonpb.CloudProvider_AZURE:
		return `#!/bin/bash
sudo apt -y update
sudo apt -y install apt-transport-https ca-certificates curl
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash
sudo su adminuser -c "az login --identity --allow-no-subscriptions" # az creds dont propagate across users apparently
`
	case commonpb.CloudProvider_GCP:
		return `#!/bin/bash
sudo apt -y update
sudo apt -y install apt-transport-https ca-certificates curl
sudo snap install google-cloud-cli --classic
`
	}
	return ""
}

func getSecretReadCommand(vaultName string, secretId string, cloud commonpb.CloudProvider) string {
	switch cloud {
	case commonpb.CloudProvider_GCP:
		return fmt.Sprintf("gcloud secrets versions access latest --secret=%s && echo", secretId)
	case commonpb.CloudProvider_AZURE:
		return fmt.Sprintf("az keyvault secret show --vault-name '%s' -n '%s' --query value -o tsv", vaultName, secretId)
	case commonpb.CloudProvider_AWS:
		return fmt.Sprintf("aws --region=eu-west-1 ssm get-parameter --name \"/%s/%s\" --with-decryption --output text --query Parameter.Value", vaultName, secretId)
	}
	return "exit 1"
}

func getSecretWriteCommand(vaultName string, secretId string, cloud commonpb.CloudProvider) string {
	switch cloud {
	case commonpb.CloudProvider_GCP:
		return fmt.Sprintf("printf \"new-secret\" | gcloud secrets versions add %s --data-file=-", secretId)
	case commonpb.CloudProvider_AZURE:
		return fmt.Sprintf("az keyvault secret set --vault-name '%s' -n '%s' --value new-secret", vaultName, secretId)
	case commonpb.CloudProvider_AWS:
		return fmt.Sprintf("aws --region=eu-west-1 ssm put-parameter --name \"/%s/%s\" --value new-value --overwrite", vaultName, secretId)
	}
	return "exit 1"
}

func TestAwsVault(t *testing.T) {
	t.Parallel()
	testVaultSecret(t, commonpb.CloudProvider_AWS)
}
func TestAzureVault(t *testing.T) {
	t.Parallel()
	testVaultSecret(t, commonpb.CloudProvider_AZURE)
}
func TestGcpVault(t *testing.T) {
	t.Parallel()
	testVaultSecret(t, commonpb.CloudProvider_GCP)
}
