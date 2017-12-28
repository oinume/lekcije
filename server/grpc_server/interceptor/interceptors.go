package interceptor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const apiTokenMetadataKey = "api-token"

func WithUnaryServerInterceptors() grpc.ServerOption {
	interceptors := []grpc.UnaryServerInterceptor{}
	interceptors = append(
		interceptors,
		DBUnaryServerInterceptor(),
		APITokenUnaryServerInterceptor(),
	)
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

const maxDBConnections = 5

func DBUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		db, err := model.OpenDB(
			bootstrap.ServerEnvVars.DBURL(),
			maxDBConnections,
			!config.IsProductionEnv(),
		)
		if err != nil {
			return handler(ctx, req)
		}
		defer db.Close()
		// TODO: redis
		return handler(context_data.SetDB(ctx, db), req)
	}
}
