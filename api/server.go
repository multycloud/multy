package api

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"multy-go/api/proto"
	"multy-go/api/proto/common"
	"multy-go/api/proto/config"
	"multy-go/api/proto/resources"
	"multy-go/db"
	"net"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	db *db.Database
}

func (s Server) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := extractUserId(ctx)
	if err != nil {
		return nil, err
	}

	a, err := anypb.New(in)
	if err != nil {
		return nil, err
	}
	resource := config.Resource{
		ResourceId: uuid.NewString(),
		Resource:   a,
	}

	c, err := s.db.Load(userId)
	if err != nil {
		return nil, err
	}
	c.Resources = append(c.Resources, &resource)
	err = s.db.Store(c)
	if err != nil {
		return nil, err
	}
	return s.ReadVirtualNetwork(ctx, &resources.ReadVirtualNetworkRequest{
		ResourceId: resource.ResourceId,
	})
}

func extractUserId(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	userIds := md.Get("user_id")
	if len(userIds) == 0 {
		return "", fmt.Errorf("user id must be set")
	}
	if len(userIds) > 1 {
		return "", fmt.Errorf("only expected 1 user id, found %d", len(userIds))
	}
	return userIds[0], nil
}

func (s Server) ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := extractUserId(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.db.Load(userId)
	if err != nil {
		return nil, err
	}
	for _, r := range c.Resources {
		if r.ResourceId == in.ResourceId {
			converted, err := convertVirtualNetworks(r.Resource)
			if err != nil {
				return nil, err
			}
			return &resources.VirtualNetworkResource{
				CommonParameters: &common.CommonResourceParameters{ResourceId: r.ResourceId},
				Resources:        converted,
			}, nil
		}
	}

	return nil, fmt.Errorf("resource with id %s not found", in.ResourceId)
}

func convertVirtualNetworks(resource *any.Any) ([]*resources.CloudSpecificVirtualNetworkResource, error) {
	out := resources.CreateVirtualNetworkRequest{}
	err := resource.UnmarshalTo(&out)
	if err != nil {
		return nil, err
	}

	var result []*resources.CloudSpecificVirtualNetworkResource
	for _, r := range out.Resources {
		result = append(result, &resources.CloudSpecificVirtualNetworkResource{
			CommonParameters: convertCommonParams(r.CommonParameters),
			Name:             r.Name,
			CidrBlock:        r.CidrBlock,
		})
	}

	return result, nil
}

func convertCommonParams(parameters *common.CloudSpecificCreateResourceCommonParameters) *common.CloudSpecificCommonResourceParameters {
	return &common.CloudSpecificCommonResourceParameters{
		ResourceGroupId: parameters.ResourceGroupId,
		Location:        parameters.Location,
		CloudProvider:   parameters.CloudProvider,
		NeedsUpdate:     false,
	}
}

func (s Server) UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	userId, err := extractUserId(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.db.Load(userId)
	if err != nil {
		return nil, err
	}
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == in.ResourceId }); i != -1 {
		a, err := anypb.New(&resources.CreateVirtualNetworkRequest{Resources: in.Resources})
		if err != nil {
			return nil, err
		}
		c.Resources[i].Resource = a
		err = s.db.Store(c)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("resource with id %s not found", in.ResourceId)
	}

	return s.ReadVirtualNetwork(ctx, &resources.ReadVirtualNetworkRequest{
		ResourceId: in.ResourceId,
	})
}

func (s Server) DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	userId, err := extractUserId(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.db.Load(userId)
	if err != nil {
		return nil, err
	}
	if i := slices.IndexFunc(c.Resources, func(r *config.Resource) bool { return r.ResourceId == in.ResourceId }); i != -1 {
		c.Resources = append(c.Resources[:i], c.Resources[i+1:]...)
		err = s.db.Store(c)
		if err != nil {
			return nil, err
		}
	}
	return &common.Empty{}, nil
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
	server := Server{}
	d, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	server.db = d
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
