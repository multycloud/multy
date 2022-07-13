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

	// wait a bit so that the vm is reachable
	time.Sleep(3 * time.Minute)

	conn, err := ssh.Dial("tcp", vm.PublicIp+":22", config)
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
	output, err := session.CombinedOutput(getSecretCommand(vault.Name, vaultSecret.Name, cloud))
	if err != nil {
		t.Logf("command outputted: %s", output)
		t.Fatal(fmt.Errorf("error running command: %+v", err))
	}

	assert.Equal(t, "test-value\n", string(output), config)

	// TODO: add test to check if there's an error when no access policy is present
	// TODO: add test for OWNER and WRITER policies
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

func getSecretCommand(vaultName string, secretId string, cloud commonpb.CloudProvider) string {
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
