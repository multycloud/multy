//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestDriftDetection(t *testing.T) {
	t.Parallel()

	ctx := getCtx(t, commonpb.CloudProvider_AWS, "drift")
	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_1,
			CloudProvider: commonpb.CloudProvider_AWS,
		},
		Name:      "drift-test-vn",
		CidrBlock: vnCidrBlock,
	}}
	vn, err := server.VnService.Create(ctx, createVnRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}
	cleanup(t, ctx, server.VnService, vn)
	assert.Nil(t, vn.CommonParameters.ResourceStatus)

	// aws ec2 modify-vpc-attribute --no-enable-dns-support --vpc-id vpc-0af2d686ed858f734 --region us-west-1
	out, err := exec.Command("aws", "ec2", "modify-vpc-attribute", "--no-enable-dns-support", "--vpc-id", vn.AwsOutputs.VpcId, "--region", "us-west-1").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	_, err = server.RefreshState(ctx, &proto.RefreshStateRequest{Cloud: commonpb.CloudProvider_AWS})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to refresh state: %+v", err)
	}

	readVn, err := server.VnService.Read(ctx, &resourcespb.ReadVirtualNetworkRequest{ResourceId: vn.CommonParameters.ResourceId})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to read vn: %+v", err)
	}

	assert.Equal(t, commonpb.ResourceStatus_NEEDS_UPDATE, readVn.CommonParameters.ResourceStatus.Statuses["aws_vpc"])

	// aws ec2 detach-internet-gateway --internet-gateway-id igw-02a9aedfe64b9eca9 --vpc-id vpc-0af2d686ed858f734 --region us-west-1
	out, err = exec.Command("aws", "ec2", "detach-internet-gateway", "--vpc-id", vn.AwsOutputs.VpcId, "--internet-gateway-id", vn.AwsOutputs.InternetGatewayId, "--region", "us-west-1").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	// aws ec2 delete-internet-gateway --internet-gateway-id igw-02a9aedfe64b9eca9 --region us-west-1
	out, err = exec.Command("aws", "ec2", "delete-internet-gateway", "--internet-gateway-id", vn.AwsOutputs.InternetGatewayId, "--region", "us-west-1").CombinedOutput()
	if err != nil {
		t.Fatal(fmt.Errorf("command failed.\n err: %s\noutput: %s", err.Error(), string(out)))
	}

	_, err = server.RefreshState(ctx, &proto.RefreshStateRequest{Cloud: commonpb.CloudProvider_AWS})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to refresh state: %+v", err)
	}

	readVn, err = server.VnService.Read(ctx, &resourcespb.ReadVirtualNetworkRequest{ResourceId: vn.CommonParameters.ResourceId})
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to read vn: %+v", err)
	}

	assert.Equal(t, commonpb.ResourceStatus_NEEDS_UPDATE, readVn.CommonParameters.ResourceStatus.Statuses["aws_vpc"])
	assert.Equal(t, commonpb.ResourceStatus_NEEDS_CREATE, readVn.CommonParameters.ResourceStatus.Statuses["aws_internet_gateway"])
}
