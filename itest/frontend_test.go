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

func TestBasic(t *testing.T) {
	b, err := sqlbackend.New(context.Background())
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.Handle(frontendv1connect.NewFrontendServiceHandler(frontend.New(b)))
	server := httptest.NewServer(mux)
	t.Cleanup(func() { server.Close() })

	key, value := int64(1), "value"
	client := frontendv1connect.NewFrontendServiceClient(http.DefaultClient, server.URL)

	// First Get should fail since key doesn't exist yet
	_, err = client.Get(context.Background(), connect.NewRequest(frontendpb.GetRequest_builder{
		Key: key,
	}.Build()))
	require.Error(t, err)
	connectErr, ok := err.(*connect.Error)
	require.True(t, ok, "expected connect.Error")
	require.Equal(t, connect.CodeInternal, connectErr.Code())

	_, err = client.Put(context.Background(), connect.NewRequest(frontendpb.PutRequest_builder{
		Key:   key,
		Value: value,
	}.Build()))
	require.NoError(t, err)

	getResp, err := client.Get(context.Background(), connect.NewRequest(frontendpb.GetRequest_builder{
		Key: key,
	}.Build()))
	require.NoError(t, err)
	require.Equal(t, value, getResp.Msg.GetValue())
}
