//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
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
		Name:             prefix + "-test-public-subnet",
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
	t.Cleanup(func() {
		if DestroyAfter {
			_, err := server.NetworkSecurityGroupService.Delete(ctx, &resourcespb.DeleteNetworkSecurityGroupRequest{ResourceId: nsg.CommonParameters.ResourceId})
			if err != nil {
				logGrpcErrorDetails(t, err)
				t.Logf("unable to delete resource %+v", err)
			}
		}
	})

	return publicSubnet, nsg
}
