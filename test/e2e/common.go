//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"google.golang.org/protobuf/proto"
	"testing"
)

const (
	vnCidrBlock      = "10.0.0.0/16"
	publicSubnetCidr = "10.0.0.0/24"
)

func createNetworkWithInternetAccess(t *testing.T, ctx context.Context, location commonpb.Location,
	cloud commonpb.CloudProvider, prefix string) (*resourcespb.SubnetResource, *resourcespb.NetworkSecurityGroupResource) {
	createVnRequest := &resourcespb.CreateVirtualNetworkRequest{Resource: &resourcespb.VirtualNetworkArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:      prefix + "-test-vn",
		CidrBlock: vnCidrBlock,
	}}
	vn, err := server.VnService.Create(ctx, createVnRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}
	cleanup(t, ctx, server.VnService, vn)

	createPublicSubnetRequest := &resourcespb.CreateSubnetRequest{Resource: &resourcespb.SubnetArgs{
		Name:             prefix + "-test-public-subnet",
		CidrBlock:        publicSubnetCidr,
		VirtualNetworkId: vn.CommonParameters.ResourceId,
	}}
	publicSubnet, err := server.SubnetService.Create(ctx, createPublicSubnetRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create publicSubnet: %+v", err)
	}
	cleanup(t, ctx, server.SubnetService, publicSubnet)

	createRtRequest := &resourcespb.CreateRouteTableRequest{Resource: &resourcespb.RouteTableArgs{
		Name:             prefix + "-test-rt",
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

	createRtaRequest := &resourcespb.CreateRouteTableAssociationRequest{Resource: &resourcespb.RouteTableAssociationArgs{
		SubnetId:     publicSubnet.CommonParameters.ResourceId,
		RouteTableId: rt.CommonParameters.ResourceId,
	}}
	rta, err := server.RouteTableAssociationService.Create(ctx, createRtaRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create route table association: %+v", err)
	}
	cleanup(t, ctx, server.RouteTableAssociationService, rta)

	createNsgRequest := &resourcespb.CreateNetworkSecurityGroupRequest{Resource: &resourcespb.NetworkSecurityGroupArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      location,
			CloudProvider: cloud,
		},
		Name:             prefix + "-test-nsg",
		VirtualNetworkId: publicSubnet.VirtualNetworkId,
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
	cleanup(t, ctx, server.NetworkSecurityGroupService, nsg)

	return publicSubnet, nsg
}

func cleanup[Arg proto.Message, OutT proto.Message](t *testing.T, ctx context.Context, s services.Service[Arg, OutT], r OutT) {
	desc := r.ProtoReflect().Descriptor().Fields().ByName("common_parameters")
	resourceId := r.ProtoReflect().Get(desc).Message().Interface().(interface {
		GetResourceId() string
	}).GetResourceId()
	t.Cleanup(func() {
		if *destroyAfter {
			_, err := s.Delete(ctx, &resourcespb.DeleteVirtualNetworkRequest{ResourceId: resourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})
}
