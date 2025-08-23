package itest

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

// mockBackend implements the Backend interface but always returns errors
type mockBackend struct{}

func (m *mockBackend) Put(ctx context.Context, key int64, value string) error {
	return errors.New("mock database error on Put")
}

func (m *mockBackend) Get(ctx context.Context, key int64) (string, error) {
	return "", errors.New("mock database error on Get")
}

// setupTestServer creates a gRPC test server with the given backend
func setupTestServer(t *testing.T, backend sqlbackend.Backend) frontendpb.FrontendServiceClient {
	// Create in-memory gRPC server
	lis := bufconn.Listen(1024 * 1024)

	s, err := frontend.NewServer(t.Context(), frontend.DefaultConfig(), backend)
	require.NoError(t, err)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()
	t.Cleanup(func() { s.Stop() })

	// Create client connection
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return frontendpb.NewFrontendServiceClient(conn)
}

func TestBasic(t *testing.T) {
	backend, err := sqlbackend.New(t.Context())
	require.NoError(t, err)

	client := setupTestServer(t, backend)
	key, value := int64(1), "value"

	// First Get should fail since key doesn't exist yet
	_, err = client.Get(t.Context(), frontendpb.GetRequest_builder{
		Key: key,
	}.Build())
	require.Error(t, err)
	grpcErr, ok := status.FromError(err)
	require.True(t, ok, "expected gRPC error")
	require.Equal(t, codes.Internal, grpcErr.Code())
	require.Contains(t, grpcErr.Message(), "no rows")

	// Put the key-value pair
	_, err = client.Put(t.Context(), frontendpb.PutRequest_builder{
		Key:   key,
		Value: value,
	}.Build())
	require.NoError(t, err)

	// Get should now succeed
	getResp, err := client.Get(t.Context(), frontendpb.GetRequest_builder{
		Key: key,
	}.Build())
	require.NoError(t, err)
	require.Equal(t, value, getResp.GetValue())
}

func TestErrorHandling(t *testing.T) {
	// Use mock backend that always returns errors
	backend := &mockBackend{}
	client := setupTestServer(t, backend)
	key, value := int64(1), "test-value"

	// Test Put error propagation
	_, err := client.Put(t.Context(), frontendpb.PutRequest_builder{
		Key:   key,
		Value: value,
	}.Build())
	require.Error(t, err)
	grpcErr, ok := status.FromError(err)
	require.True(t, ok, "expected gRPC error")
	require.Equal(t, codes.Internal, grpcErr.Code())
	require.Contains(t, grpcErr.Message(), "mock database error on Put")

	// Test Get error propagation
	_, err = client.Get(t.Context(), frontendpb.GetRequest_builder{
		Key: key,
	}.Build())
	require.Error(t, err)
	grpcErr, ok = status.FromError(err)
	require.True(t, ok, "expected gRPC error")
	require.Equal(t, codes.Internal, grpcErr.Code())
	require.Contains(t, grpcErr.Message(), "mock database error on Get")
}
