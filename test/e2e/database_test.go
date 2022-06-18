//go:build e2e
// +build e2e

package e2e

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func testDatabase(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "database")

	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_2,
			CloudProvider: cloud,
		},
		Name:      "db-test-vn",
		CidrBlock: "10.0.0.0/16",
	}}
	vn, err := server.VnService.Create(ctx, createVnRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}
	cleanup(t, ctx, server.VnService, vn)

	createPublicSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "db-test-public-subnet1",
		CidrBlock:        "10.0.0.0/24",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		AvailabilityZone: 1,
	}}
	publicSubnet, err := server.SubnetService.Create(ctx, createPublicSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create publicSubnet: %+v", err)
	}
	cleanup(t, ctx, server.SubnetService, publicSubnet)

	createSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             "db-test-public-subnet2",
		CidrBlock:        "10.0.1.0/24",
		VirtualNetworkId: vn.CommonParameters.ResourceId,
		AvailabilityZone: 2,
	}}
	subnet, err := server.SubnetService.Create(ctx, createSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create subnet: %+v", err)
	}
	cleanup(t, ctx, server.SubnetService, subnet)

	createRtRequest := &resourcespb.CreateRouteTableRequest{Resource: &resourcespb.RouteTableArgs{
		Name:             "db-test-rt",
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
	cleanup(t, ctx, server.RouteTableService, rt)

	subnetIds := []string{publicSubnet.CommonParameters.ResourceId, subnet.CommonParameters.ResourceId}
	for _, id := range subnetIds {
		createRtaRequest := &resourcespb.CreateRouteTableAssociationRequest{Resource: &resourcespb.RouteTableAssociationArgs{
			SubnetId:     id,
			RouteTableId: rt.CommonParameters.ResourceId,
		}}
		rta, err := server.RouteTableAssociationService.Create(ctx, createRtaRequest)
		if err != nil {
			logGrpcErrorDetails(t, err)
			t.Fatalf("unable to create route table association: %+v", err)
		}
		cleanup(t, ctx, server.RouteTableAssociationService, rta)
	}

	createDbRequest := &resourcespb.CreateDatabaseRequest{Resource: &resourcespb.DatabaseArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_2,
			CloudProvider: cloud,
		},
		Name:          "multydbtest" + common.RandomString(2),
		Engine:        resourcespb.DatabaseEngine_MYSQL,
		EngineVersion: "5.7",
		StorageGb:     10,
		Size:          commonpb.DatabaseSize_MICRO,
		Username:      "multyuser",
		// azure requires complex stuff
		Password:  common.RandomString(8) + "-2Ab",
		SubnetIds: subnetIds,
	}}
	db, err := server.DatabaseService.Create(ctx, createDbRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create route table association: %+v", err)
	}
	cleanup(t, ctx, server.DatabaseService, db)

	out, err := exec.Command("mysql", "-h", db.Host, "-P", "3306", "-u", db.ConnectionUsername, "--password="+db.Password, "-e", "select 12+34;").CombinedOutput()
	if err != nil {
		t.Fatalf("command failed.\n err: %s\noutput: %s", err.Error(), string(out))
	}
	assert.Contains(t, string(out), "46")
}

func TestAwsDatabase(t *testing.T) {
	t.Parallel()
	testDatabase(t, commonpb.CloudProvider_AWS)
}
func TestAzureDatabase(t *testing.T) {
	t.Parallel()
	testDatabase(t, commonpb.CloudProvider_AZURE)
}
