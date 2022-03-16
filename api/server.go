package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"multy-go/api/proto"
	"multy-go/api/proto/resources"
	"net"
)

type server struct {
	proto.UnimplementedMultyResourceServiceServer
}

func (s *server) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	// TODO: implement
	return nil, fmt.Errorf("unimplemented")
}

func RunServer(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	go func() {
		<-ctx.Done()
		s.GracefulStop()
		_ = lis.Close()
	}()
	proto.RegisterMultyResourceServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
