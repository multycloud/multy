//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, createVnRequest.GetResource().GetCommonParameters().GetLocation(), vn.GetCommonParameters().GetLocation())
	assert.Equal(t, createVnRequest.GetResource().GetCommonParameters().GetCloudProvider(), vn.GetCommonParameters().GetCloudProvider())
	assert.Nil(t, vn.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createVnRequest.GetResource().GetName(), vn.GetName())
	assert.Equal(t, createVnRequest.GetResource().GetCidrBlock(), vn.GetCidrBlock())

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

	assert.Nil(t, publicSubnet.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createPublicSubnetRequest.GetResource().GetName(), publicSubnet.GetName())
	assert.Equal(t, createPublicSubnetRequest.GetResource().GetCidrBlock(), publicSubnet.GetCidrBlock())
	assert.Equal(t, createPublicSubnetRequest.GetResource().GetVirtualNetworkId(), publicSubnet.GetVirtualNetworkId())

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

	assert.Nil(t, rt.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createRtRequest.GetResource().GetName(), rt.GetName())
	assert.Len(t, rt.GetRoutes(), 1)
	assert.Equal(t, createRtRequest.GetResource().GetRoutes()[0].GetCidrBlock(), rt.GetRoutes()[0].GetCidrBlock())
	assert.Equal(t, createRtRequest.GetResource().GetRoutes()[0].GetDestination(), rt.GetRoutes()[0].GetDestination())
	assert.Equal(t, createRtRequest.GetResource().GetVirtualNetworkId(), rt.GetVirtualNetworkId())

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

	assert.Nil(t, rta.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createRtaRequest.GetResource().GetSubnetId(), rta.GetSubnetId())
	assert.Equal(t, createRtaRequest.GetResource().GetRouteTableId(), rta.GetRouteTableId())

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

	assert.Equal(t, createNsgRequest.GetResource().GetCommonParameters().GetLocation(), nsg.GetCommonParameters().GetLocation())
	assert.Equal(t, createNsgRequest.GetResource().GetCommonParameters().GetCloudProvider(), nsg.GetCommonParameters().GetCloudProvider())
	assert.Nil(t, nsg.GetCommonParameters().GetResourceStatus())

	assert.Equal(t, createNsgRequest.GetResource().GetName(), nsg.GetName())
	assert.Equal(t, createNsgRequest.GetResource().GetVirtualNetworkId(), nsg.GetVirtualNetworkId())
	assert.Len(t, nsg.GetRules(), len(createNsgRequest.GetResource().GetRules()))
	if len(nsg.GetRules()) == len(createNsgRequest.GetResource().GetRules()) {
		for i, rule := range createNsgRequest.GetResource().GetRules() {
			assert.Equal(t, rule.GetCidrBlock(), nsg.GetRules()[i].GetCidrBlock())
			assert.Equal(t, rule.GetDirection(), nsg.GetRules()[i].GetDirection())
			assert.Equal(t, rule.GetPortRange().GetFrom(), nsg.GetRules()[i].GetPortRange().GetFrom())
			assert.Equal(t, rule.GetPortRange().GetTo(), nsg.GetRules()[i].GetPortRange().GetTo())
			assert.Equal(t, rule.GetPriority(), nsg.GetRules()[i].GetPriority())
		}
	}

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
