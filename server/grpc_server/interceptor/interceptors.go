package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

const apiTokenMetadataKey = "api-token"

func WithUnaryServerInterceptors() grpc.ServerOption {
	interceptors := []grpc.UnaryServerInterceptor{}
	interceptors = append(interceptors, APITokenUnaryServerInterceptor())
	return grpc_middleware.WithUnaryServerChain(interceptors...)
}

func APITokenUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		values, ok := md[apiTokenMetadataKey]
		if !ok {
			return handler(ctx, req)
		}
		return handler(context_data.WithAPIToken(ctx, values[0]), req)
	}
}
