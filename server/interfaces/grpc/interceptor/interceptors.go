package interceptor

import (
	"context"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/oinume/lekcije/server/ga_measurement"
)

const apiTokenMetadataKey = "api-token"

func WithUnaryServerInterceptors() grpc.ServerOption {
	interceptors := []grpc.UnaryServerInterceptor{
		APITokenUnaryServerInterceptor(),
		GAMeasurementEventInterceptor(),
	}
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

		eventValues := &ga_measurement.EventValues{}
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
		c := ga_measurement.WithEventValues(ctx, eventValues)
		return handler(c, req)
	}
}
