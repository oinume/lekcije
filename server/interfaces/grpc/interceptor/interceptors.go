package interceptor

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/event_logger"
	"github.com/oinume/lekcije/server/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const apiTokenMetadataKey = "api-token"

func WithUnaryServerInterceptors() grpc.ServerOption {
	interceptors := []grpc.UnaryServerInterceptor{}
	interceptors = append(
		interceptors,
		APITokenUnaryServerInterceptor(),
		GAMeasurementEventInterceptor(),
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

func GAMeasurementEventInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		eventValues := &event_logger.GAMeasurementEventValues{}
		for key, values := range md {
			if !strings.HasPrefix(key, "http-") {
				continue
			}
			switch key {
			case "http-user-agent":
				eventValues.UserAgentOverride = values[0]
			case "http-referer":
				eventValues.DocumentReferrer = values[0]
			case "http-host":
				eventValues.DocumentHostName = values[0]
			case "http-url-path":
				eventValues.DocumentTitle = values[0]
				eventValues.DocumentPath = values[0]
			case "http-remote-addr":
				eventValues.IPOverride = values[0]
			case "http-tracking-id":
				eventValues.ClientID = values[0]
			}
		}
		c := event_logger.WithGAMeasurementEventValues(ctx, eventValues)
		return handler(c, req)
	}
}

const maxDBConnections = 5

func DBUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		db, err := model.OpenDB(
			config.DefaultVars.DBURL(),
			maxDBConnections,
			config.DefaultVars.DebugSQL,
		)
		if err != nil {
			return handler(ctx, req)
		}
		defer db.Close()
		// TODO: redis
		return handler(context_data.SetDB(ctx, db), req)
	}
}
