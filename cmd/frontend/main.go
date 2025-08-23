package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/earthboundkid/versioninfo/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rajatgoel/gh-go/internal/frontend"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
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
	server := grpc.NewServer()

	// Register service
	frontendpb.RegisterFrontendServiceServer(server, frontend.New(backend))

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("frontend.v1.FrontendService", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service
	reflection.Register(server)

	// Listen on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf("::%d", *port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	slog.Info("starting gRPC server", "port", *port)
	if err := server.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
