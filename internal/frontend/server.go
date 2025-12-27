package frontend

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/dynoinc/gh-go/internal/sqlbackend"
	frontendpb "github.com/dynoinc/gh-go/proto/frontend/v1"
)

// ServerOption configures the server.
type ServerOption func(*serverConfig)

type serverConfig struct {
	noopTelemetry bool
}

// WithNoopTelemetry disables OTLP exporters and uses noop telemetry providers.
// This is useful for testing to avoid connection timeouts.
func WithNoopTelemetry() ServerOption {
	return func(c *serverConfig) {
		c.noopTelemetry = true
	}
}

// NewServer creates a new gRPC server with health checks, reflection, and OpenTelemetry instrumentation.
// The returned cleanup function must be called during shutdown to flush telemetry exporters.
func NewServer(ctx context.Context, backend sqlbackend.Backend, opts ...ServerOption) (*grpc.Server, func(), error) {
	cfg := &serverConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	cleanup := func() {}

	if !cfg.noopTelemetry {
		cleanupFn, err := setupOTLP(ctx)
		if err != nil {
			return nil, nil, err
		}
		cleanup = cleanupFn
	}

	// Create gRPC server with middleware chain
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(
				logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
					slog.Log(ctx, slog.Level(lvl), msg, fields...)
				}),
			),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// Register the main service
	frontendpb.RegisterFrontendServiceServer(server, New(backend))

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus(frontendpb.FrontendService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service
	reflection.Register(server)

	return server, cleanup, nil
}

func setupOTLP(ctx context.Context) (func(), error) {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcessPID(),
		resource.WithProcessExecutableName(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown trace provider", "error", err)
		}
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown metric provider", "error", err)
		}
	}, nil
}
