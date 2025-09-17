package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/earthboundkid/versioninfo/v2"
	"golang.org/x/sync/errgroup"

	"github.com/rajatgoel/gh-go/internal/config"
	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func main() {
	versioninfo.AddFlag(flag.CommandLine)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	backend, err := sqlbackend.New(ctx)
	if err != nil {
		slog.Error("failed to create backend", "error", err)
		os.Exit(1)
	}

	// Create gRPC server with OpenTelemetry instrumentation (enabled by default)
	server, otelCleanup, err := frontend.NewServer(ctx, cfg, backend)
	if err != nil {
		slog.Error("failed to create gRPC server", "error", err)
		os.Exit(1)
	}
	slog.Info("created gRPC server with OpenTelemetry instrumentation")

	// Listen on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// Create errgroup for coordinating goroutines
	g, ctx := errgroup.WithContext(ctx)

	// Start server in a goroutine
	g.Go(func() error {
		slog.Info("starting gRPC server", "port", cfg.Port)
		if err := server.Serve(lis); err != nil {
			slog.Error("server error", "error", err)
			return err
		}
		return nil
	})

	// Start signal handler in a separate goroutine
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-quit:
			slog.Info("received shutdown signal")
			server.GracefulStop() // Stop the server immediately
			cancel()              // Cancel the context to signal other goroutines
		case <-ctx.Done():
			// Context was cancelled by another goroutine (e.g., server error)
			slog.Info("context cancelled, shutting down signal handler")
		}

		return nil
	})

	// Wait for either goroutine to complete or context to be cancelled
	if err := g.Wait(); err != nil {
		slog.Error("goroutine error", "error", err)
	}

	slog.Info("shutting down gRPC server...")

	if otelCleanup != nil {
		// Flush telemetry with a fresh context so shutdown succeeds even if the main context was cancelled.
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cleanupCancel()
		otelCleanup(cleanupCtx)
	}

	slog.Info("gRPC server stopped")
}
