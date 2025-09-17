package frontend

import (
	"context"
	"testing"
	"time"

	"github.com/rajatgoel/gh-go/internal/config"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func TestNewServerReturnsCleanup(t *testing.T) {
	cfg := &config.Config{
		ServiceName: "gh-go-frontend-test",
		Environment: "test",
		Port:        0,
	}

	backend, err := sqlbackend.New(context.Background())
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}

	server, cleanup, err := NewServer(context.Background(), cfg, backend)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	if cleanup == nil {
		t.Fatal("expected non-nil cleanup function")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cleanup(shutdownCtx)
	server.Stop()
}
