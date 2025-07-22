package itest

import (
	"context"
	"errors"
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

// mockBackend implements the Backend interface but always returns errors
type mockBackend struct{}

func (m *mockBackend) Put(ctx context.Context, key int64, value string) error {
	return errors.New("mock database error on Put")
}

func (m *mockBackend) Get(ctx context.Context, key int64) (string, error) {
	return "", errors.New("mock database error on Get")
}

func TestErrorHandling(t *testing.T) {
	// Use mock backend that always returns errors
	b := &mockBackend{}

	mux := http.NewServeMux()
	mux.Handle(frontendv1connect.NewFrontendServiceHandler(frontend.New(b)))
	server := httptest.NewServer(mux)
	t.Cleanup(func() { server.Close() })

	client := frontendv1connect.NewFrontendServiceClient(http.DefaultClient, server.URL)
	key, value := int64(1), "test-value"

	// Test Put error propagation
	_, err := client.Put(context.Background(), connect.NewRequest(frontendpb.PutRequest_builder{
		Key:   key,
		Value: value,
	}.Build()))
	require.Error(t, err)
	connectErr, ok := err.(*connect.Error)
	require.True(t, ok, "expected connect.Error")
	require.Equal(t, connect.CodeInternal, connectErr.Code())
	require.Contains(t, connectErr.Message(), "mock database error on Put")

	// Test Get error propagation
	_, err = client.Get(context.Background(), connect.NewRequest(frontendpb.GetRequest_builder{
		Key: key,
	}.Build()))
	require.Error(t, err)
	connectErr, ok = err.(*connect.Error)
	require.True(t, ok, "expected connect.Error")
	require.Equal(t, connect.CodeInternal, connectErr.Code())
	require.Contains(t, connectErr.Message(), "mock database error on Get")
}

// TestSQLErrors tests with a real SQL backend but triggers errors
func TestSQLErrors(t *testing.T) {
	b, err := sqlbackend.New(context.Background())
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.Handle(frontendv1connect.NewFrontendServiceHandler(frontend.New(b)))
	server := httptest.NewServer(mux)
	t.Cleanup(func() { server.Close() })

	client := frontendv1connect.NewFrontendServiceClient(http.DefaultClient, server.URL)

	// Try to get a key that doesn't exist, which should return sql.ErrNoRows
	_, err = client.Get(context.Background(), connect.NewRequest(frontendpb.GetRequest_builder{
		Key: 999, // non-existent key
	}.Build()))
	require.Error(t, err)
	connectErr, ok := err.(*connect.Error)
	require.True(t, ok, "expected connect.Error")
	require.Equal(t, connect.CodeInternal, connectErr.Code())
	// The specific SQL error message may vary by implementation, but it should be about no rows
	require.Contains(t, connectErr.Message(), "no rows")
}
