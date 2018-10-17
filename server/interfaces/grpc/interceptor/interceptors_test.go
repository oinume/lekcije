package interceptor

import (
	"context"
	"testing"

	"github.com/oinume/lekcije/server/context_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAPITokenUnaryServerInterceptor(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	const apiToken = "abc"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(apiTokenMetadataKey, apiToken))
	info := &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: "FooService.FooMethod",
	}
	var gotAPIToken string
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		v, err := context_data.GetAPIToken(ctx)
		if err != nil {
			return nil, err
		}
		gotAPIToken = v
		return "ok", nil
	}

	_, err := APITokenUnaryServerInterceptor()(ctx, nil, info, handler)
	r.NoError(err)
	a.Equal(apiToken, gotAPIToken)
}
