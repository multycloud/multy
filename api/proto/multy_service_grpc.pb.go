// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	common "github.com/multycloud/multy/api/proto/common"
	resources "github.com/multycloud/multy/api/proto/resources"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MultyResourceServiceClient is the client API for MultyResourceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultyResourceServiceClient interface {
	CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error)
	ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error)
	UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error)
	DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest, opts ...grpc.CallOption) (*common.Empty, error)
}

type multyResourceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMultyResourceServiceClient(cc grpc.ClientConnInterface) MultyResourceServiceClient {
	return &multyResourceServiceClient{cc}
}

func (c *multyResourceServiceClient) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error) {
	out := new(resources.VirtualNetworkResource)
	err := c.cc.Invoke(ctx, "/dev.multy.MultyResourceService/CreateVirtualNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multyResourceServiceClient) ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error) {
	out := new(resources.VirtualNetworkResource)
	err := c.cc.Invoke(ctx, "/dev.multy.MultyResourceService/ReadVirtualNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multyResourceServiceClient) UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest, opts ...grpc.CallOption) (*resources.VirtualNetworkResource, error) {
	out := new(resources.VirtualNetworkResource)
	err := c.cc.Invoke(ctx, "/dev.multy.MultyResourceService/UpdateVirtualNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multyResourceServiceClient) DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest, opts ...grpc.CallOption) (*common.Empty, error) {
	out := new(common.Empty)
	err := c.cc.Invoke(ctx, "/dev.multy.MultyResourceService/DeleteVirtualNetwork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MultyResourceServiceServer is the server API for MultyResourceService service.
// All implementations must embed UnimplementedMultyResourceServiceServer
// for forward compatibility
type MultyResourceServiceServer interface {
	CreateVirtualNetwork(context.Context, *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error)
	ReadVirtualNetwork(context.Context, *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error)
	UpdateVirtualNetwork(context.Context, *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error)
	DeleteVirtualNetwork(context.Context, *resources.DeleteVirtualNetworkRequest) (*common.Empty, error)
	mustEmbedUnimplementedMultyResourceServiceServer()
}

// UnimplementedMultyResourceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMultyResourceServiceServer struct {
}

func (UnimplementedMultyResourceServiceServer) CreateVirtualNetwork(context.Context, *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVirtualNetwork not implemented")
}
func (UnimplementedMultyResourceServiceServer) ReadVirtualNetwork(context.Context, *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadVirtualNetwork not implemented")
}
func (UnimplementedMultyResourceServiceServer) UpdateVirtualNetwork(context.Context, *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVirtualNetwork not implemented")
}
func (UnimplementedMultyResourceServiceServer) DeleteVirtualNetwork(context.Context, *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVirtualNetwork not implemented")
}
func (UnimplementedMultyResourceServiceServer) mustEmbedUnimplementedMultyResourceServiceServer() {}

// UnsafeMultyResourceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultyResourceServiceServer will
// result in compilation errors.
type UnsafeMultyResourceServiceServer interface {
	mustEmbedUnimplementedMultyResourceServiceServer()
}

func RegisterMultyResourceServiceServer(s grpc.ServiceRegistrar, srv MultyResourceServiceServer) {
	s.RegisterService(&MultyResourceService_ServiceDesc, srv)
}

func _MultyResourceService_CreateVirtualNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(resources.CreateVirtualNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultyResourceServiceServer).CreateVirtualNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.multy.MultyResourceService/CreateVirtualNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultyResourceServiceServer).CreateVirtualNetwork(ctx, req.(*resources.CreateVirtualNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultyResourceService_ReadVirtualNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(resources.ReadVirtualNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultyResourceServiceServer).ReadVirtualNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.multy.MultyResourceService/ReadVirtualNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultyResourceServiceServer).ReadVirtualNetwork(ctx, req.(*resources.ReadVirtualNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultyResourceService_UpdateVirtualNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(resources.UpdateVirtualNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultyResourceServiceServer).UpdateVirtualNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.multy.MultyResourceService/UpdateVirtualNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultyResourceServiceServer).UpdateVirtualNetwork(ctx, req.(*resources.UpdateVirtualNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultyResourceService_DeleteVirtualNetwork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(resources.DeleteVirtualNetworkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultyResourceServiceServer).DeleteVirtualNetwork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dev.multy.MultyResourceService/DeleteVirtualNetwork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultyResourceServiceServer).DeleteVirtualNetwork(ctx, req.(*resources.DeleteVirtualNetworkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MultyResourceService_ServiceDesc is the grpc.ServiceDesc for MultyResourceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultyResourceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dev.multy.MultyResourceService",
	HandlerType: (*MultyResourceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateVirtualNetwork",
			Handler:    _MultyResourceService_CreateVirtualNetwork_Handler,
		},
		{
			MethodName: "ReadVirtualNetwork",
			Handler:    _MultyResourceService_ReadVirtualNetwork_Handler,
		},
		{
			MethodName: "UpdateVirtualNetwork",
			Handler:    _MultyResourceService_UpdateVirtualNetwork_Handler,
		},
		{
			MethodName: "DeleteVirtualNetwork",
			Handler:    _MultyResourceService_DeleteVirtualNetwork_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/multy_service.proto",
}
