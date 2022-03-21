package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services/network_interface"
	"github.com/multycloud/multy/api/services/route_table"
	"github.com/multycloud/multy/api/services/route_table_association"
	"github.com/multycloud/multy/api/services/subnet"
	"github.com/multycloud/multy/api/services/virtual_network"
	"github.com/multycloud/multy/db"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	virtual_network.VnService
	subnet.SubnetService
	network_interface.NetworkInterfaceService
	route_table.RouteTableService
	route_table_association.RouteTableAssociationService
}

func RunServer(ctx context.Context, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	go func() {
		<-ctx.Done()
		s.GracefulStop()
		_ = lis.Close()
	}()
	d, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	server := Server{
		proto.UnimplementedMultyResourceServiceServer{},
		virtual_network.NewVnService(d),
		subnet.NewSubnetServiceService(d),
		network_interface.NewNetworkInterfaceServiceService(d),
		route_table.NewRouteTableServiceService(d),
		route_table_association.NewRouteTableAssociationServiceService(d),
	}
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Create(ctx, in)
}
func (s *Server) ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Read(ctx, in)
}
func (s *Server) UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Update(ctx, in)
}
func (s *Server) DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	return s.VnService.Service.Delete(ctx, in)
}

func (s *Server) CreateSubnet(ctx context.Context, in *resources.CreateSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Create(ctx, in)
}
func (s *Server) ReadSubnet(ctx context.Context, in *resources.ReadSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Read(ctx, in)
}
func (s *Server) UpdateSubnet(ctx context.Context, in *resources.UpdateSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Update(ctx, in)
}
func (s *Server) DeleteSubnet(ctx context.Context, in *resources.DeleteSubnetRequest) (*common.Empty, error) {
	return s.SubnetService.Service.Delete(ctx, in)
}

func (s *Server) CreateNetworkInterface(ctx context.Context, in *resources.CreateNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkInterface(ctx context.Context, in *resources.ReadNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkInterface(ctx context.Context, in *resources.UpdateNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkInterface(ctx context.Context, in *resources.DeleteNetworkInterfaceRequest) (*common.Empty, error) {
	return s.NetworkInterfaceService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTable(ctx context.Context, in *resources.CreateRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTable(ctx context.Context, in *resources.ReadRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTable(ctx context.Context, in *resources.UpdateRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTable(ctx context.Context, in *resources.DeleteRouteTableRequest) (*common.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTableAssociation(ctx context.Context, in *resources.CreateRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTableAssociation(ctx context.Context, in *resources.ReadRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTableAssociation(ctx context.Context, in *resources.UpdateRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTableAssociation(ctx context.Context, in *resources.DeleteRouteTableAssociationRequest) (*common.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}
