package frontend

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"github.com/rajatgoel/gh-go/internal/config"
	"github.com/rajatgoel/gh-go/internal/sqlbackend"
	frontendpb "github.com/rajatgoel/gh-go/proto/frontend/v1"
)

// loggingInterceptor logs each RPC call in a compact format
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Log the RPC call in a compact format
	duration := time.Since(start)
	status := "OK"
	if err != nil {
		status = "ERROR"
	}

	slog.Info("RPC",
		"method", info.FullMethod,
		"duration", duration,
		"status", status,
	)

	return resp, err
}

// NewServer creates a new gRPC server with health checks, reflection, and OpenTelemetry instrumentation.
// The returned cleanup function must be called during shutdown to flush telemetry exporters.
func NewServer(ctx context.Context, cfg *config.Config, backend sqlbackend.Backend) (*grpc.Server, func(context.Context), error) {
	// Setup OpenTelemetry with default configuration
	otelCleanup, err := SetupOTEL(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup OpenTelemetry: %w", err)
	}

	// Create gRPC server with OTEL stats handler and logging interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
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

	return server, otelCleanup, nil
}
