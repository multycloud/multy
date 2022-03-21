package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services/database"
	"github.com/multycloud/multy/api/services/network_interface"
	"github.com/multycloud/multy/api/services/network_security_group"
	"github.com/multycloud/multy/api/services/object_storage"
	"github.com/multycloud/multy/api/services/object_storage_object"
	"github.com/multycloud/multy/api/services/public_ip"
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
	network_security_group.NetworkSecurityGroupService
	database.DatabaseService
	object_storage.ObjectStorageService
	object_storage_object.ObjectStorageObjectService
	public_ip.PublicIpService
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
		subnet.NewSubnetService(d),
		network_interface.NewNetworkInterfaceService(d),
		route_table.NewRouteTableService(d),
		route_table_association.NewRouteTableAssociationService(d),
		network_security_group.NewNetworkSecurityGroupService(d),
		database.NewDatabaseService(d),
		object_storage.NewObjectStorageService(d),
		object_storage_object.NewObjectStorageObjectService(d),
		public_ip.NewPublicIpService(d),
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

func (s *Server) CreateNetworkSecurityGroup(ctx context.Context, in *resources.CreateNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkSecurityGroup(ctx context.Context, in *resources.ReadNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkSecurityGroup(ctx context.Context, in *resources.UpdateNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkSecurityGroup(ctx context.Context, in *resources.DeleteNetworkSecurityGroupRequest) (*common.Empty, error) {
	return s.NetworkSecurityGroupService.Service.Delete(ctx, in)
}

func (s *Server) CreateDatabase(ctx context.Context, in *resources.CreateDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Create(ctx, in)
}
func (s *Server) ReadDatabase(ctx context.Context, in *resources.ReadDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Read(ctx, in)
}
func (s *Server) UpdateDatabase(ctx context.Context, in *resources.UpdateDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Update(ctx, in)
}
func (s *Server) DeleteDatabase(ctx context.Context, in *resources.DeleteDatabaseRequest) (*common.Empty, error) {
	return s.DatabaseService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorage(ctx context.Context, in *resources.CreateObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorage(ctx context.Context, in *resources.ReadObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorage(ctx context.Context, in *resources.UpdateObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorage(ctx context.Context, in *resources.DeleteObjectStorageRequest) (*common.Empty, error) {
	return s.ObjectStorageService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorageObject(ctx context.Context, in *resources.CreateObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorageObject(ctx context.Context, in *resources.ReadObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorageObject(ctx context.Context, in *resources.UpdateObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorageObject(ctx context.Context, in *resources.DeleteObjectStorageObjectRequest) (*common.Empty, error) {
	return s.ObjectStorageService.Service.Delete(ctx, in)
}

func (s *Server) CreatePublicIp(ctx context.Context, in *resources.CreatePublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Create(ctx, in)
}
func (s *Server) ReadPublicIp(ctx context.Context, in *resources.ReadPublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Read(ctx, in)
}
func (s *Server) UpdatePublicIp(ctx context.Context, in *resources.UpdatePublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Update(ctx, in)
}
func (s *Server) DeletePublicIp(ctx context.Context, in *resources.DeletePublicIpRequest) (*common.Empty, error) {
	return s.PublicIpService.Service.Delete(ctx, in)
}
