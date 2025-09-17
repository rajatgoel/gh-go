package itest

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/bufconn"

	"github.com/rajatgoel/gh-go/client"
	"github.com/rajatgoel/gh-go/internal/config"
	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

// mockBackend implements the Backend interface but always returns errors
type mockBackend struct{}

func (m *mockBackend) Put(ctx context.Context, key int64, value string) error {
	return errors.New("mock database error on Put")
}

func (m *mockBackend) Get(ctx context.Context, key int64) (string, error) {
	return "", errors.New("mock database error on Get")
}

func (m *mockBackend) Close(context.Context) error {
	return nil
}

// setupTestServer creates a gRPC test server with the given backend
func setupTestServer(t *testing.T, backend sqlbackend.Backend) (*client.Client, func()) {
	// Create in-memory gRPC server using bufconn
	lis := bufconn.Listen(1024 * 1024)

	// Create test config
	cfg := &config.Config{
		ServiceName: "gh-go-frontend-test",
		Environment: "test",
		Port:        5051, // Not used with bufconn, but required by config
	}

	// Create server
	s, err := frontend.NewServer(t.Context(), cfg, backend)
	require.NoError(t, err)

	// Start server in background
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited: %v", err)
		}
	}()

	// Create client using custom dialer for bufconn
	c, err := client.New(
		client.WithTarget("passthrough:///bufnet"),
		client.WithDialer(func(ctx context.Context, target string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		client.WithInsecure(),
	)
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		c.Close()
		s.Stop()
		require.NoError(t, backend.Close(t.Context()))
	}

	return c, cleanup
}

func TestBasic(t *testing.T) {
	backend, err := sqlbackend.New(t.Context())
	require.NoError(t, err)

	c, cleanup := setupTestServer(t, backend)
	defer cleanup()

	key, value := int64(1), "value"

	// First Get should fail since key doesn't exist yet
	_, err = c.Get(t.Context(), key)
	require.Error(t, err)

	// Put the key-value pair
	err = c.Put(t.Context(), key, value)
	require.NoError(t, err)

	// Get should now succeed
	getResp, err := c.Get(t.Context(), key)
	require.NoError(t, err)
	require.Equal(t, value, getResp)
}

func TestErrorHandling(t *testing.T) {
	// Use mock backend that always returns errors
	backend := &mockBackend{}
	c, cleanup := setupTestServer(t, backend)
	defer cleanup()

	key, value := int64(1), "test-value"

	// Test Put error propagation
	err := c.Put(t.Context(), key, value)
	require.Error(t, err)

	// Test Get error propagation
	_, err = c.Get(t.Context(), key)
	require.Error(t, err)
}
