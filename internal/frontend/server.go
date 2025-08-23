package frontend

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

// NewServer creates a new gRPC server with health checks, reflection, and OpenTelemetry instrumentation
func NewServer(ctx context.Context, config *Config, backend sqlbackend.Backend) (*grpc.Server, error) {
	// Setup OpenTelemetry with default configuration
	if _, err := SetupOTEL(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to setup OpenTelemetry: %w", err)
	}

	// Create gRPC server with OTEL stats handler
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithMeterProvider(otel.GetMeterProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

	// Register the main service
	frontendpb.RegisterFrontendServiceServer(server, New(backend))

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus(frontendpb.FrontendService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service
	reflection.Register(server)

	return server, nil
}
