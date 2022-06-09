package deploy

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	metadata2 "github.com/multycloud/multy/resources/types/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"testing"
)

// TODO: tests
// - provider from prev resource is included
// - get state calls terraform show
// - init calls terraform init

func init() {
	flags.Environment = flags.Local
}

func mockCreds(t *testing.T) context.Context {
	credentials := &credspb.CloudCredentials{
		AwsCreds: &credspb.AwsCredentials{
			AccessKey: "mock",
			SecretKey: "mock",
		},
		AzureCreds: &credspb.AzureCredentials{
			SubscriptionId: "mock",
			TenantId:       "mock",
			ClientId:       "mock",
			ClientSecret:   "mock",
		},
	}
	b, err := proto.Marshal(credentials)
	if err != nil {
		t.Fatalf("unable to marshal creds: %s", err)
	}

	md := map[string][]string{"api_key": {"test-deploy"}, "cloud-creds-bin": {string(b)}}
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return ctx
}

func TestDeploy_rollbacksIfSomethingFails(t *testing.T) {
	ctx := mockCreds(t)
	mockTfCmd := &MockTerraformCommand{}
	sut := deploy.DeploymentExecutor{
		TfCmd: mockTfCmd,
	}

	config, err := resources.LoadConfig(&configpb.Config{
		UserId:          "test",
		Resources:       nil,
		ResourceCounter: 1,
	}, metadata2.Metadatas)
	if err != nil {
		t.Fatalf("can't load config, %s", err)
	}

	res, err := config.CreateResource(&resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_EAST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Fatalf("can't create resource, %s", err)
	}

	mockTfCmd.
		On("Init", mock.Anything, mock.Anything).
		Return(nil)
	mockTfCmd.
		On("Apply", mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("fake error")).
		Once()
	mockTfCmd.
		On("Apply", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	mockTfCmd.
		On("GetState", mock.Anything, mock.Anything).
		Return(&output.TfState{}, nil)
	_, err = sut.Deploy(ctx, config, nil, nil)
	assert.Error(t, err)

	mockTfCmd.AssertNumberOfCalls(t, "Apply", 2)

	// assert that in the end the main.tf was applied without the new resource (aka rollback)
	file, err := os.ReadFile(filepath.Join(deploy.GetTempDirForUser(false, "test"), "main.tf"))
	if err != nil {
		t.Fatalf("can't read main.tf, %s", err)
	}
	assert.NotContains(t, string(file), res.GetResourceId())
}

func TestDeploy_callsTfApply(t *testing.T) {
	ctx := mockCreds(t)
	mockTfCmd := &MockTerraformCommand{}
	sut := deploy.DeploymentExecutor{
		TfCmd: mockTfCmd,
	}

	config, err := resources.LoadConfig(&configpb.Config{
		UserId:          "test",
		Resources:       nil,
		ResourceCounter: 1,
	}, metadata2.Metadatas)
	if err != nil {
		t.Fatalf("can't load config, %s", err)
	}

	res, err := config.CreateResource(&resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_EAST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Fatalf("can't create resource, %s", err)
	}

	mockTfCmd.
		On("Init", mock.Anything, mock.Anything).
		Return(nil)
	mockTfCmd.
		On("Apply", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTfCmd.
		On("GetState", mock.Anything, mock.Anything).
		Return(&output.TfState{}, nil)
	_, err = sut.Deploy(ctx, config, nil, nil)
	if err != nil {
		t.Fatalf("can't deploy, %s", err)
	}

	mockTfCmd.AssertNumberOfCalls(t, "Apply", 1)
	// assert that in the end the main.tf contains the resource id
	file, err := os.ReadFile(filepath.Join(deploy.GetTempDirForUser(false, "test"), "main.tf"))
	assert.NoError(t, err, "unable to read main.tf")
	assert.Contains(t, string(file), res.GetResourceId())
}

func TestDeploy_onlyAffectedResources(t *testing.T) {
	ctx := mockCreds(t)
	mockTfCmd := &MockTerraformCommand{}
	sut := deploy.DeploymentExecutor{
		TfCmd: mockTfCmd,
	}

	config, err := resources.LoadConfig(&configpb.Config{
		UserId:          "test",
		Resources:       nil,
		ResourceCounter: 1,
	}, metadata2.Metadatas)
	if err != nil {
		t.Fatalf("can't load config, %s", err)
	}

	r1, err := config.CreateResource(&resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_EAST_1,
			CloudProvider: commonpb.CloudProvider_AZURE,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Fatalf("can't create resource, %s", err)
	}

	mockTfCmd.
		On("Init", mock.Anything, mock.Anything).
		Return(nil)
	mockTfCmd.
		On("Apply", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	mockTfCmd.
		On("GetState", mock.Anything, mock.Anything).
		Return(&output.TfState{}, nil)
	_, err = sut.Deploy(ctx, config, nil, r1)
	if err != nil {
		t.Fatalf("can't deploy, %s", err)
	}

	res2, err := config.CreateResource(&resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_EAST_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "test-vn",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Fatalf("can't create resource, %s", err)
	}

	mockTfCmd.
		On("Apply", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	_, err = sut.Deploy(ctx, config, nil, res2)
	if err != nil {
		t.Fatalf("can't deploy, %s", err)
	}

	affectedResources := config.GetAffectedResources(res2.GetResourceId())
	assert.NotEmpty(t, affectedResources)
	mockTfCmd.AssertCalled(t, "Apply", mock.Anything, mock.Anything, affectedResources)
}

func TestRefresh(t *testing.T) {
	ctx := mockCreds(t)
	mockTfCmd := &MockTerraformCommand{}
	sut := deploy.DeploymentExecutor{
		TfCmd: mockTfCmd,
	}

	config, err := resources.LoadConfig(&configpb.Config{
		UserId:          "test",
		Resources:       nil,
		ResourceCounter: 1,
	}, metadata2.Metadatas)
	if err != nil {
		t.Fatalf("can't load config, %s", err)
	}

	mockTfCmd.
		On("Init", mock.Anything, mock.Anything).
		Return(nil)
	mockTfCmd.
		On("Refresh", mock.Anything, mock.Anything).
		Return(nil)
	err = sut.RefreshState(ctx, "test", config)
	if err != nil {
		t.Fatalf("can't refresh, %s", err)
	}

	mockTfCmd.AssertNumberOfCalls(t, "Refresh", 1)
}
