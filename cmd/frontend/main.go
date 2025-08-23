package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/earthboundkid/versioninfo/v2"

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
)

func main() {
	port := flag.Int("port", 5051, "The server port")
	versioninfo.AddFlag(nil)
	flag.Parse()

	ctx := context.Background()
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

	// Start server on specified port
	slog.Info("starting gRPC server", "port", *port)
	if err := server.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
