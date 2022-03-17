package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services/virtual_network"
	"github.com/multycloud/multy/db"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	virtual_network.VnService
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
		virtual_network.VnService{Db: d},
	}
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Create(ctx, in)
}
func (s *Server) ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Read(ctx, in)
}
func (s *Server) UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Update(ctx, in)
}
func (s *Server) DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	return s.VnService.Delete(ctx, in)
}
