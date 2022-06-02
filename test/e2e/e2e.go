//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/service_context"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"testing"
)

const DestroyAfter = true

func getCtx(t *testing.T, cloud commonpb.CloudProvider, testName string) context.Context {
	accessKeyId, exists := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !exists {
		t.Fatalf("AWS_ACCESS_KEY_ID not set")
	}
	accessKeySecret, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !exists {
		t.Fatalf("AWS_SECRET_ACCESS_KEY not set")
	}
	azSubscriptionId, exists := os.LookupEnv("ARM_SUBSCRIPTION_ID")
	if !exists {
		t.Fatalf("ARM_SUBSCRIPTION_ID not set")
	}
	azClientId, exists := os.LookupEnv("ARM_CLIENT_ID")
	if !exists {
		t.Fatalf("ARM_CLIENT_ID not set")
	}
	azClientSecret, exists := os.LookupEnv("ARM_CLIENT_SECRET")
	if !exists {
		t.Fatalf("ARM_CLIENT_SECRET not set")
	}
	azTenantId, exists := os.LookupEnv("ARM_TENANT_ID")
	if !exists {
		t.Fatalf("ARM_TENANT_ID not set")
	}

	credentials := &credspb.CloudCredentials{
		AwsCreds: &credspb.AwsCredentials{
			AccessKey: accessKeyId,
			SecretKey: accessKeySecret,
		},
		AzureCreds: &credspb.AzureCredentials{
			SubscriptionId: azSubscriptionId,
			TenantId:       azTenantId,
			ClientId:       azClientId,
			ClientSecret:   azClientSecret,
		},
	}
	b, err := proto.Marshal(credentials)
	if err != nil {
		t.Fatalf("unable to marshal creds: %s", err)
	}

	md := map[string][]string{"api_key": {fmt.Sprintf("test-%s-%s", testName, cloud.String())}, "cloud-creds-bin": {string(b)}}
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return ctx
}

func logGrpcErrorDetails(t *testing.T, err error) {
	if s, ok := status.FromError(err); ok {
		for _, details := range s.Details() {
			if msg, ok := details.(interface{ GetErrorMessage() string }); ok {
				t.Logf("server returned error: %s", msg.GetErrorMessage())
			}
		}
	}
}

var server *api.Server

func init() {
	if server != nil {
		return
	}
	flags.Environment = flags.Local
	flags.NoTelemetry = true
	flags.DryRun = false
	awsClient, err := aws_client.NewClient()
	if err != nil {
		log.Fatalf("failed to initialize aws client: %v", err)
	}
	database, err := db.NewDatabase(awsClient)
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	serviceContext := &service_context.ResourceServiceContext{
		Database:           database,
		AwsClient:          awsClient,
		DeploymentExecutor: deploy.NewDeploymentExecutor(),
	}

	s := api.CreateServer(serviceContext)
	server = &s
}
