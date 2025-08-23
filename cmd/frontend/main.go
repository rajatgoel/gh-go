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

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func main() {
	port := flag.Int("port", 5051, "The server port")
	versioninfo.AddFlag(nil)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	backend, err := sqlbackend.New(ctx)
	if err != nil {
		slog.Error("failed to create backend", "error", err)
		os.Exit(1)
	}

	// Create gRPC server
	server := frontend.NewServer(backend)

	// Listen on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf("::%d", *port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		slog.Info("starting gRPC server", "port", *port)
		if err := server.Serve(lis); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down gRPC server...")

	// Give outstanding RPCs a chance to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdownCancel()

	// Gracefully stop the server
	server.GracefulStop()

	// Wait for server to stop
	<-shutdownCtx.Done()
	slog.Info("gRPC server stopped")
}
