package frontend

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

// NewServer creates a new gRPC server with health checks and reflection enabled
func NewServer(backend sqlbackend.Backend) *grpc.Server {
	server := grpc.NewServer()

	// Register the main service
	frontendpb.RegisterFrontendServiceServer(server, New(backend))

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("frontend.v1.FrontendService", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service
	reflection.Register(server)

	return server
}
