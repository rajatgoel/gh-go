package itest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
	frontendv1connect "github.com/rajatgoel/gh-go/proto/frontend/v1/v1connect"
)

func TestStub(t *testing.T) {
	b, err := sqlbackend.New(context.Background())
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.Handle(frontendv1connect.NewFrontendServiceHandler(frontend.New(b)))
	server := httptest.NewServer(mux)
	t.Cleanup(func() { server.Close() })

	key, value := int64(1), "value"
	client := frontendv1connect.NewFrontendServiceClient(http.DefaultClient, server.URL)

	resp, err := client.Get(context.Background(), connect.NewRequest(&frontendpb.GetRequest{
		Key: key,
	}))
	require.NoError(t, err)
	require.Empty(t, resp.Msg.Value)

	_, err = client.Put(context.Background(), connect.NewRequest(&frontendpb.PutRequest{
		Key:   key,
		Value: value,
	}))
	require.NoError(t, err)

	resp, err = client.Get(context.Background(), connect.NewRequest(&frontendpb.GetRequest{
		Key: key,
	}))
	require.NoError(t, err)
	require.Equal(t, value, resp.Msg.Value)
}
