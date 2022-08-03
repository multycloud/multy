package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"runtime/trace"
)

func TraceUnaryInterceptor() grpc.UnaryServerInterceptor {
	return TraceInterceptor
}

func TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, task := trace.NewTask(ctx, "request")
	trace.Log(ctx, "method", info.FullMethod)
	defer task.End()
	// Invoke the original method call
	return handler(ctx, req)
}
